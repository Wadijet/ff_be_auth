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
	UserCRUD           services.Repository
	RoleCRUD           services.Repository
	PermissionCRUD     services.Repository
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
		notRoleError := "You have not selected a role to work with!"
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
											ctx.SetUserValue("roleId", _token.RoleID)      // set loggedIn user role id in contextß
											isRightToken = true
											break
										}
									}

									// Nếu token không hợp lệ
									if isRightToken == false {
										utility.JSON(ctx, utility.Payload(false, nil, notAuthError))
									} else {

										// Kiểm tra xem người dùng đã chọn RoleID để làm việc hay chưa
										if ctx.UserValue("roleId") == nil {
											utility.JSON(ctx, utility.Payload(false, nil, notRoleError))
										} else { // Nếu đã chọn RoleID thì tiếp tục kiểm tra
										

											// Nếu requirePermissions có số lượng bằng 0
											if len(requirePermissions) == 0 {
												next(ctx)
											} else { // Nếu requirePermissions có số lượng lớn hơn 0
												// Kiểm tra xem người dùng có quyền yêu cầu hay không
												roleID := ctx.UserValue("roleId").(string)
												// Tìm tất cả RolePermission dựa trên RoleID
												findRolePermission, err := jt.RolePermissionCRUD.FindAll(context.TODO(), bson.M{"roleId":roleID}, nil)
												if findRolePermission == nil {
													utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
												} else {
													// chuyển findRolePermission từ dạng []interface{} thành dạng []models.RolePermission
													var rolePermissions []models.RolePermission
													bsonBytes, err := bson.Marshal(findRolePermission)
													if err != nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													} 
													err = bson.Unmarshal(bsonBytes, &rolePermissions)
													if err != nil {
														utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
													}
													
													// Duyệt qua danh sách quyền của vai trò 
													// Nếu có quyền yêu cầu nào không nằm trong danh sách quyền của vai trò thì trả về lỗi JSON
													isRightPermission := false
													for _, requireRermission := range requirePermissions {
														for _, rolePermission := range rolePermissions {
															// Tìm permission trong danh sách các quyền theo PermissionID
															resultPermission, _ := jt.PermissionCRUD.FindOneById(context.TODO(), utility.ObjectID2String(rolePermission.PermissionID), nil)
															if resultPermission != nil {
																
																var permission models.Permission
																bsonBytes, err := bson.Marshal(resultPermission)
																if err != nil {
																	utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
																}

																err = bson.Unmarshal(bsonBytes, &permission)
																if err != nil {
																	utility.JSON(ctx, utility.Payload(false, err, notPermissionError))
																}

																if permission.Name == requireRermission {
																	isRightPermission = true
																	break
																}
															}
														}
													}
													if isRightPermission {
														next(ctx)
													} else {
														utility.JSON(ctx, utility.Payload(false, nil, notPermissionError))
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

