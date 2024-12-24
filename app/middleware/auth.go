package middleware

import (
	"context"
	"fmt"
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
	UserCRUD           services.Repository
	RoleCRUD           services.Repository
	PermissionCRUD     services.Repository
	RolePermissionCRUD services.Repository
	UserRoleCRUD       services.Repository
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

		// Lấy JWT token từ header "Authorization"
		jwtTokenString := string(ctx.Request.Header.Peek("Authorization"))
		if jwtTokenString != "" {
			// Tách JWT token từ header "Authorization Bearer" bằng cách tách chuỗi
			splitToken := strings.Split(jwtTokenString, "Bearer ")
			if len(splitToken) > 1 {

				// Lấy chuỗi JWT token từ phần tử thứ 2 của mảng sau khi tách
				jwtTokenString = splitToken[1]

				// Giải mã JWT token và kiểm tra tính hợp lệ của token
				t := models.JwtToken{}
				jwtToken, err := jwt.ParseWithClaims(jwtTokenString, &t, func(token *jwt.Token) (interface{}, error) {
					return []byte(jt.C.JwtSecret), nil
				})

				// Nếu có lỗi hoặc token không hợp lệ thì trả về lỗi JSON
				if err != nil || !jwtToken.Valid {
					utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
				} else { // Nếu token hợp lệ thì tiếp tục kiểm tra

					// Tìm người dùng dựa trên ID trong token
					findUser, err := jt.UserCRUD.FindOneById(context.TODO(), t.UserID, nil)
					// Nếu không tìm thấy người dùng thì trả về lỗi JSON
					if findUser == nil {
						utility.JSON(ctx, utility.Payload(false, err, notAuthError))
					} else { // Nếu tìm thấy người dùng thì tiếp tục kiểm tra
						// Chuyển kết quả tìm kiếm thành dạng []byte
						var user models.User
						bsonBytes, err := bson.Marshal(findUser)
						// Nếu có lỗi trong quá trình chuyển đổi thì trả về lỗi JSON
						if err != nil {
							utility.JSON(ctx, utility.Payload(false, err, notAuthError))
						} else { // Nếu không có lỗi thì tiếp tục kiểm tra

							// Chuyển kết quả tìm kiếm từ dạng []byte thành cấu trúc User
							err = bson.Unmarshal(bsonBytes, &user)
							if err != nil {
								utility.JSON(ctx, utility.Payload(false, err, notAuthError))
							} else {

								// Kiểm tra xem người dùng có bị khóa hay không
								if user.IsBlock {
									utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
								} else { // Nếu không bị khóa thì tiếp tục kiểm tra

									// Kiểm tra xem token có tồn tại trong danh sách Tokens của người dùng hay không. Nếu có thì lưu thông tin vào context
									isRightToken := false
									for _, _token := range user.Tokens {
										if _token.JwtToken == jwtTokenString {
											ctx.SetUserValue("userId", t.UserID)           // set loggedIn user id in context
											ctx.SetUserValue("userToken", _token.JwtToken) // set loggedIn user token in context
											isRightToken = true
											break
										}
									}
									// Nếu token không hợp lệ
									if isRightToken == false {
										utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
									} else {
										// Nếu requirePermissions có số lượng bằng 0
										if requirePermission == "" {
											next(ctx) // Gọi fasthttp.RequestHandler tiếp theo
										} else {
											// Tìm requirePermissionID dựa trên requirePermission
											findPermission, err := jt.PermissionCRUD.FindOne(context.TODO(), bson.M{"name": requirePermission}, nil)
											if findPermission == nil {
												utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
											} else {
												// chuyển findPermission từ dạng interface{} thành dạng models.Permission
												var permission models.Permission
												bsonBytes, _ := bson.Marshal(findPermission)
												err := bson.Unmarshal(bsonBytes, &permission)
												if err != nil {
													utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
												} else {
													requirePermissionID := permission.ID

													// Tìm tất cả các Role của người dùng dựa trên UserID
													findRoles, err := jt.UserRoleCRUD.FindAll(context.TODO(), bson.M{"userId": utility.String2ObjectID(t.UserID)}, nil)
													if findRoles == nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													} else {

														var isRightRole bool = false
														var minScope byte = 1

														// Duyệt qua danh sách các Role của người dùng
														for _, findRoleData := range findRoles {
															// decode permission từ bson.M về models.Permission
															var modelUserRole models.UserRole
															bsonBytes, _ := bson.Marshal(findRoleData)
															err := bson.Unmarshal(bsonBytes, &modelUserRole)
															if err != nil {
																fmt.Errorf("Failed to decode permission")
																continue
															}

															// Tìm tất cả các Permission của Role dựa trên RoleID
															findRolePermissions, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId": modelUserRole.RoleID}, nil)
															if findRolePermissions == nil {
																continue
															}

															// duyệt qua danh sách các Permission của Role
															for _, findRolePermission := range findRolePermissions {
																var modelRolePermission models.RolePermission
																bsonBytes, _ := bson.Marshal(findRolePermission)
																err := bson.Unmarshal(bsonBytes, &modelRolePermission)
																if err != nil {
																	fmt.Errorf("Failed to decode permission")
																	continue
																}

																if modelRolePermission.PermissionID == requirePermissionID {
																	isRightRole = true
																	if modelRolePermission.Scope < minScope {
																		minScope = modelRolePermission.Scope
																	}
																}
															}
															if isRightRole == false {
																utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
															} else {
																ctx.SetUserValue("minScope", minScope) // set minScope in context
																next(ctx)
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
				}
			} else {
				utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
			}
		} else {
			utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
		}
	}
}
