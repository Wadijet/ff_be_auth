# Agent Management APIs

Tài liệu về các API endpoints quản lý Agent (trợ lý tự động).

## 📋 Tổng Quan

Tất cả các API Agent đều nằm dưới prefix `/api/v1/agent/`.

## 🔐 Agent CRUD APIs

Quản lý Agents.

**Prefix:** `/api/v1/agent/`

**Endpoints (Full CRUD):**
- `POST /api/v1/agent/insert-one` - Tạo agent (Permission: `Agent.Insert`)
- `GET /api/v1/agent/find` - Tìm agents (Permission: `Agent.Read`)
- `GET /api/v1/agent/find-by-id/:id` - Tìm theo ID (Permission: `Agent.Read`)
- `PUT /api/v1/agent/update-by-id/:id` - Cập nhật agent (Permission: `Agent.Update`)
- `DELETE /api/v1/agent/delete-by-id/:id` - Xóa agent (Permission: `Agent.Delete`)

## 🔐 Agent Check-In/Check-Out APIs

### Check-In

Đánh dấu agent check-in.

**Endpoint:** `POST /api/v1/agent/check-in/:id`

**Authentication:** Cần (Permission: `Agent.CheckIn`)

**Path Parameters:**
- `id`: Agent ID

**Response 200:**
```json
{
  "data": {
    "message": "Agent checked in successfully",
    "checkInTime": "2024-01-01T08:00:00Z"
  },
  "error": null
}
```

### Check-Out

Đánh dấu agent check-out.

**Endpoint:** `POST /api/v1/agent/check-out/:id`

**Authentication:** Cần (Permission: `Agent.CheckOut`)

**Path Parameters:**
- `id`: Agent ID

**Response 200:**
```json
{
  "data": {
    "message": "Agent checked out successfully",
    "checkOutTime": "2024-01-01T17:00:00Z"
  },
  "error": null
}
```

## 📝 Lưu Ý

- Tất cả endpoints đều yêu cầu authentication
- Mỗi endpoint yêu cầu permission tương ứng
- Check-in/check-out được sử dụng để theo dõi thời gian làm việc của agent

## 📚 Tài Liệu Liên Quan

- [Facebook Integration APIs](facebook.md)
- [Pancake Integration APIs](pancake.md)

