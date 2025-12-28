# Fiber v3 - Cách Đăng Ký Middleware Đúng

## ⚠️ Vấn Đề

Fiber v3 có bug với cách đăng ký middleware trực tiếp trong route. Cách sau **KHÔNG HOẠT ĐỘNG**:

```go
// ❌ SAI - Middleware không được gọi
router.Get("/auth/roles", middleware.AuthMiddleware(""), userHandler.HandleGetUserRoles)
router.Post("/auth/logout", middleware.AuthMiddleware(""), userHandler.HandleLogout)
```

## ✅ Cách Đúng

Phải dùng `.Use()` method thông qua `registerRouteWithMiddleware`:

```go
// ✅ ĐÚNG - Dùng registerRouteWithMiddleware với .Use()
authRolesMiddleware := middleware.AuthMiddleware("")
registerRouteWithMiddleware(router, "/auth", "GET", "/roles", []fiber.Handler{authRolesMiddleware}, userHandler.HandleGetUserRoles)
```

## 📝 Chi Tiết

### Hàm `registerRouteWithMiddleware`

```go
// registerRouteWithMiddleware đăng ký route với middleware sử dụng .Use() method (cách đúng theo Fiber v3)
func registerRouteWithMiddleware(router fiber.Router, prefix string, method string, path string, middlewares []fiber.Handler, handler fiber.Handler) {
	// Tạo group với prefix, middleware sẽ chỉ áp dụng cho routes trong group này
	routeGroup := router.Group(prefix)
	for _, mw := range middlewares {
		routeGroup.Use(mw)  // ← ĐÂY LÀ CÁCH ĐÚNG
	}

	// Đăng ký route với path tương đối (không có prefix vì đã có trong group)
	switch method {
	case "GET":
		routeGroup.Get(path, handler)
	case "POST":
		routeGroup.Post(path, handler)
	case "PUT":
		routeGroup.Put(path, handler)
	case "DELETE":
		routeGroup.Delete(path, handler)
	}
}
```

### Cách Sử Dụng

**Ví dụ 1: Route đơn giản với 1 middleware**
```go
authMiddleware := middleware.AuthMiddleware("")
registerRouteWithMiddleware(router, "/auth", "GET", "/roles", []fiber.Handler{authMiddleware}, userHandler.HandleGetUserRoles)
```

**Ví dụ 2: Route với nhiều middleware (như CRUD routes)**
```go
authReadMiddleware := middleware.AuthMiddleware("Permission.Read")
orgContextMiddleware := middleware.OrganizationContextMiddleware()
registerRouteWithMiddleware(router, "/permission", "GET", "/find", []fiber.Handler{authReadMiddleware, orgContextMiddleware}, permHandler.Find)
```

## 🔍 Lịch Sử

- **Ngày**: 2025-12-28
- **Vấn đề**: Endpoint `/api/v1/auth/roles` trả về 401 mặc dù token hợp lệ
- **Nguyên nhân**: Dùng cách trực tiếp `router.Get(path, middleware, handler)` - middleware không được gọi
- **Giải pháp**: Đổi sang dùng `registerRouteWithMiddleware` với `.Use()` method
- **Kết quả**: Đã test 7 cách khác nhau và chọn cách này là cách duy nhất hoạt động đúng

## 📌 Quy Tắc

**LUÔN LUÔN** dùng `registerRouteWithMiddleware` khi cần đăng ký route với middleware trong Fiber v3.

**KHÔNG BAO GIỜ** dùng cách trực tiếp:
```go
// ❌ KHÔNG DÙNG CÁCH NÀY
router.Get(path, middleware, handler)
router.Post(path, middleware, handler)
```

## 🔗 Tham Khảo

- File: `api/core/api/router/routes.go`
- Hàm: `registerRouteWithMiddleware()` (dòng 159-178)
- Tất cả CRUD routes đều dùng cách này
- Endpoint `/auth/roles` đã được sửa để dùng cách này
