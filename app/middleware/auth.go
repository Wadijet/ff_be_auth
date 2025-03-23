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
type JwtToken struct {
	C                  *config.Configuration
	UserCRUD           services.BaseService[models.User]
	RoleCRUD           services.BaseService[models.Role]
	PermissionCRUD     services.BaseService[models.Permission]
	RolePermissionCRUD services.BaseService[models.RolePermission]
	UserRoleCRUD       services.BaseService[models.UserRole]
}

// NewJwtToken , khởi tạo một JwtToken mới
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {
	newHandler := new(JwtToken)
	newHandler.C = c

	// Khởi tạo các collection
	userCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Users)).Collection(global.MongoDB_ColNames.Users)
	roleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Roles)).Collection(global.MongoDB_ColNames.Roles)
	permissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.Permissions)).Collection(global.MongoDB_ColNames.Permissions)
	rolePermissionCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.RolePermissions)).Collection(global.MongoDB_ColNames.RolePermissions)
	userRoleCol := db.Database(services.GetDBName(c, global.MongoDB_ColNames.UserRoles)).Collection(global.MongoDB_ColNames.UserRoles)

	// Khởi tạo các service với BaseService
	newHandler.UserCRUD = services.NewBaseService[models.User](userCol)
	newHandler.RoleCRUD = services.NewBaseService[models.Role](roleCol)
	newHandler.PermissionCRUD = services.NewBaseService[models.Permission](permissionCol)
	newHandler.RolePermissionCRUD = services.NewBaseService[models.RolePermission](rolePermissionCol)
	newHandler.UserRoleCRUD = services.NewBaseService[models.UserRole](userRoleCol)

	return newHandler
}

// CheckUserAuth , kiểm tra xác thực người dùng
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
					ctx.SetStatusCode(fasthttp.StatusUnauthorized)
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				findUser, err := jt.UserCRUD.FindOne(context.TODO(), t.UserID)
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusUnauthorized)
					utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					return
				}

				if findUser.IsBlock {
					ctx.SetStatusCode(fasthttp.StatusForbidden)
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				isRightToken := false
				for _, _token := range findUser.Tokens {
					if _token.JwtToken == jwtTokenString {
						ctx.SetUserValue("userId", t.UserID)
						ctx.SetUserValue("userToken", _token.JwtToken)
						isRightToken = true
						break
					}
				}
				if !isRightToken {
					ctx.SetStatusCode(fasthttp.StatusUnauthorized)
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
					return
				}

				if requirePermission == "" {
					next(ctx)
					return
				}

				findPermission, err := jt.PermissionCRUD.FindOneByFilter(context.TODO(), bson.M{"name": requirePermission}, nil)
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusForbidden)
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				requirePermissionID := findPermission.ID
				findRoles, err := jt.UserRoleCRUD.FindAll(context.TODO(), bson.M{"userId": utility.String2ObjectID(t.UserID)}, nil)
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusForbidden)
					utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
					return
				}

				isRightRole := false
				var minScope byte = 1
				for _, modelUserRole := range findRoles {
					findRolePermissions, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId": modelUserRole.RoleID}, nil)
					if err != nil {
						continue
					}

					for _, modelRolePermission := range findRolePermissions {
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
					ctx.SetStatusCode(fasthttp.StatusForbidden)
					utility.JSON(ctx, utility.Payload(false, nil, notPermissionError))
					return
				}

				ctx.SetUserValue("minScope", minScope)
				next(ctx)
			} else {
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}
