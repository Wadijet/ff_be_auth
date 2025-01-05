package middleware

import (
	"context"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	models "atk-go-server/app/models/mongodb"
	"atk-go-server/app/services"
	"atk-go-server/app/utility"
	"atk-go-server/config"
	"atk-go-server/global"
)

// JwtToken , basic jwt model
// Cấu trúc JwtToken, mô hình jwt cơ bản
type JwtToken struct {
	C                  *config.Configuration
	UserCRUD           services.RepositoryService
	RoleCRUD           services.RepositoryService
	PermissionCRUD     services.RepositoryService
	RolePermissionCRUD services.RepositoryService
	UserRoleCRUD       services.RepositoryService
}

// NewJwtToken , khởi tạo một JwtToken mới
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {

	newHandler := new(JwtToken)
	newHandler.C = c
	newHandler.UserCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Users)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newHandler.PermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newHandler.RolePermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)
	newHandler.UserRoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.UserRoles)

	return newHandler
}

// CheckUserAuth , kiểm tra xác thực người dùng
// Dành cho user
// CheckUserAuth là middleware kiểm tra quyền truy cập của người dùng dựa trên JWT token và các quyền yêu cầu.
//
// Tham số:
// - requirePermissions: Danh sách các quyền yêu cầu để truy cập vào tài nguyên.
// - next: fasthttp.RequestHandler tiếp theo sẽ được gọi nếu người dùng có quyền hợp lệ.
//
// Chức năng:
// - Kiểm tra xem header "Authorization" có chứa JWT token hợp lệ hay không.
// - Giải mã và xác thực JWT token.
// - Tìm kiếm người dùng dựa trên ID trong token.
// - Kiểm tra xem người dùng có bị khóa hay không.
// - Kiểm tra xem token có hợp lệ với người dùng hay không.
// - Nếu có quyền yêu cầu, kiểm tra xem người dùng có quyền đó hay không.
// - Nếu tất cả các kiểm tra đều thành công, gọi fasthttp.RequestHandler tiếp theo.
//
// Trả về:
// - fasthttp.RequestHandler: Handler sẽ được gọi nếu người dùng có quyền hợp lệ, nếu không sẽ trả về lỗi JSON.
func (jt *JwtToken) CheckUserAuth(requirePermission string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		notAuthError := "An unauthorized access!"
		notPermissionError := "You do not have permission to perform the action!"

		jwtTokenString := string(ctx.Request.Header.Peek("Authorization"))
		if jwtTokenString != "" {
			splitToken := strings.Split(jwtTokenString, "Bearer ")
			if len(splitToken) > 1 {
				jwtTokenString = splitToken[1]

				t := models.JwtToken{}
				jwtToken, err := jwt.ParseWithClaims(jwtTokenString, &t, func(token *jwt.Token) (interface{}, error) {
					return []byte(jt.C.JwtSecret), nil
				})

				if err != nil || !jwtToken.Valid {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				findUser, err := jt.UserCRUD.FindOneById(context.TODO(), utility.String2ObjectID(t.UserID), nil)
				if findUser == nil {
					utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					return
				}

				var user models.User
				bsonBytes, err := bson.Marshal(findUser)
				if err != nil {
					utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					return
				}
				err = bson.Unmarshal(bsonBytes, &user)
				if err != nil {
					utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					return
				}

				if user.IsBlock {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				isRightToken := false
				for _, _token := range user.Tokens {
					if _token.JwtToken == jwtTokenString {
						ctx.SetUserValue("userId", t.UserID)
						ctx.SetUserValue("userToken", _token.JwtToken)
						isRightToken = true
						break
					}
				}
				if !isRightToken {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				if requirePermission == "" {
					next(ctx)
					return
				}

				findPermission, err := jt.PermissionCRUD.FindOne(context.TODO(), bson.M{"name": requirePermission}, nil)
				if findPermission == nil {
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				var permission models.Permission
				bsonBytes, _ = bson.Marshal(findPermission)
				err = bson.Unmarshal(bsonBytes, &permission)
				if err != nil {
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				requirePermissionID := permission.ID
				findRoles, err := jt.UserRoleCRUD.FindAll(context.TODO(), bson.M{"userId": utility.String2ObjectID(t.UserID)}, nil)
				if findRoles == nil {
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				isRightRole := false
				var minScope byte = 1
				for _, findRoleData := range findRoles {
					var modelUserRole models.UserRole
					bsonBytes, _ = bson.Marshal(findRoleData)
					err = bson.Unmarshal(bsonBytes, &modelUserRole)
					if err != nil {
						continue
					}

					findRolePermissions, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId": modelUserRole.RoleID}, nil)
					if findRolePermissions == nil {
						continue
					}

					for _, findRolePermission := range findRolePermissions {
						var modelRolePermission models.RolePermission
						bsonBytes, _ = bson.Marshal(findRolePermission)
						err = bson.Unmarshal(bsonBytes, &modelRolePermission)
						if err != nil {
							continue
						}

						if modelRolePermission.PermissionID == requirePermissionID {
							isRightRole = true
							ctx.SetUserValue("RoleId", modelUserRole.RoleID)
							if modelRolePermission.Scope < minScope {
								minScope = modelRolePermission.Scope
							}
							break
						}
					}
				}
				if !isRightRole {
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				ctx.SetUserValue("minScope", minScope)
				next(ctx)
			} else {
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}
