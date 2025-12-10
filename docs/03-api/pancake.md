# Pancake Integration APIs

Tài liệu về các API endpoints tích hợp Pancake (Orders).

## 📋 Tổng Quan

Tất cả các API Pancake đều nằm dưới prefix `/api/v1/pancake/`.

## 🔐 Pancake Order APIs

Quản lý Pancake Orders.

**Prefix:** `/api/v1/pancake/order/`

**Endpoints (Full CRUD):**
- `POST /api/v1/pancake/order/insert-one` - Tạo order (Permission: `PcOrder.Insert`)
- `GET /api/v1/pancake/order/find` - Tìm orders (Permission: `PcOrder.Read`)
- `GET /api/v1/pancake/order/find-by-id/:id` - Tìm theo ID (Permission: `PcOrder.Read`)
- `GET /api/v1/pancake/order/find-by-ids` - Tìm nhiều orders theo IDs (Permission: `PcOrder.Read`)
- `GET /api/v1/pancake/order/find-with-pagination` - Tìm với phân trang (Permission: `PcOrder.Read`)
- `PUT /api/v1/pancake/order/update-by-id/:id` - Cập nhật order (Permission: `PcOrder.Update`)
- `DELETE /api/v1/pancake/order/delete-by-id/:id` - Xóa order (Permission: `PcOrder.Delete`)
- `GET /api/v1/pancake/order/count` - Đếm orders (Permission: `PcOrder.Read`)

### Ví Dụ: Tạo Order

**Request:**
```json
POST /api/v1/pancake/order/insert-one
{
  "orderId": "order-123",
  "customerId": "customer-456",
  "total": 100000,
  "status": "pending"
}
```

**Response:**
```json
{
  "data": {
    "_id": "507f1f77bcf86cd799439011",
    "orderId": "order-123",
    "customerId": "customer-456",
    "total": 100000,
    "status": "pending",
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

## 📝 Lưu Ý

- Tất cả endpoints đều yêu cầu authentication
- Mỗi endpoint yêu cầu permission `PcOrder.*` tương ứng
- Pancake integration đồng bộ dữ liệu đơn hàng từ hệ thống Pancake

## 📚 Tài Liệu Liên Quan

- [Facebook Integration APIs](facebook.md)
- [Agent Management APIs](agent.md)

