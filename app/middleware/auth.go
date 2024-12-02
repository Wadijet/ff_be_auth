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
	C              *config.Configuration
	UserCRUD       services.Repository
	RoleCRUD       services.Repository
	PermissionCRUD services.Repository
	RolePermissionCRUD services.Repository
}

// NewJwtToken , khởi tạo một JwtToken mới
func NewJwtToken(c *config.Configuration, db *mongo.Client) *JwtToken {

	newHandler := new(JwtToken)
	newHandler.C = c
	newHandler.UserCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Users)
	newHandler.RoleCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Roles)
	newHandler.PermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.Permissions)
	newHandler.RolePermissionCRUD = *services.NewRepository(c, db, global.MongoDB_ColNames.RolePermissions)

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
func (jt *JwtToken) CheckUserAuth(requirePermissions []string, next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		notAuthError := "An unauthorized access!"
		notPermissionError := "You do not have permission to perform the action!"

		tokenString := string(ctx.Request.Header.Peek("Authorization"))
		if tokenString != "" {
			splitToken := strings.Split(tokenString, "Bearer ")
			if len(splitToken) > 1 {
				tokenString = splitToken[1]

				// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
				t := models.JwtToken{}
				token, err := jwt.ParseWithClaims(tokenString, &t, func(token *jwt.Token) (interface{}, error) {
					return []byte(jt.C.JwtSecret), nil
				})

				if err != nil || !token.Valid {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
				} else {
					findUser, err := jt.UserCRUD.FindOneById(context.TODO(), t.ID, nil)
					if findUser == nil {
						utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					} else {
						var user models.User
						bsonBytes, err := bson.Marshal(findUser)
						if err != nil {
							utility.JSON(ctx, utility.Payload(false, err, notAuthError))
						} else {
							err = bson.Unmarshal(bsonBytes, &user)
							if err != nil {
								utility.JSON(ctx, utility.Payload(false, err, notAuthError))
							} else {
								if user.IsBlock {
									utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
								} else {
									isRightToken := false
									for _, _token := range user.Tokens {
										if _token.Token == tokenString {
											ctx.SetUserValue("userId", t.ID) // set loggedIn user id in context
											isRightToken = true
											break
										}
									}

									if isRightToken == false {
										utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
									} else {
										if len(requirePermissions) == 0 {
											next(ctx)
										} else {
											strRoleID := utility.ObjectID2String(user.Role)
											findRole, err := jt.RoleCRUD.FindOneById(context.TODO(), strRoleID, nil)
											if findRole == nil {
												utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
											} else {
												var result_findRole models.Role
												bsonBytes, err := bson.Marshal(findRole)
												if err != nil {
													utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
												} else {
													err = bson.Unmarshal(bsonBytes, &result_findRole)
													if err != nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													} else {
														totalCheck := true
														for _, requirePermisson := range requirePermissions {
															checkPermission := false
															for _, s := range result_findRole.Permissions {
																if requirePermisson == s.Name {
																	checkPermission = true
																	break
																}
															}

															if checkPermission == false {
																totalCheck = false
																break
															}
														}

														if totalCheck == true {
															next(ctx)
														} else {
															utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
														}

													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			} else {
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}
