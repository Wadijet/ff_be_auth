# Frontend Context Switching - Quick Reference

## 🎯 Core Concept

**Context làm việc = ROLE ID** (không phải organization ID)

## 📋 Flow Ngắn Gọn

1. **Login** → Lưu JWT token
2. **GET /auth/roles** → Lấy danh sách roles
3. **Chọn role** → Lưu `roleId` vào localStorage
4. **Mọi request** → Gửi header `X-Active-Role-ID: <roleId>` (Optional: Backend tự động detect nếu không có)

## 🔑 Key Points

### Header Phải Gửi
```
Authorization: Bearer <jwt-token>
X-Active-Role-ID: <role-id>  ← QUAN TRỌNG: ROLE ID, không phải org ID
```

### API Endpoints

**1. Login:**
```
POST /api/v1/auth/login/firebase
Body: { "idToken": "...", "hwid": "..." }
Response: { "data": { "token": "...", ... } }
```

**2. Get Roles:**
```
GET /api/v1/auth/roles
Headers: { "Authorization": "Bearer <token>" }
Response: {
  "data": [
    {
      "roleId": "...",
      "roleName": "...",
      "organizationId": "...",
      "organizationName": "..."
    }
  ]
}
```

### Logic Chọn Role

```javascript
const roles = await api.get('/auth/roles');

if (roles.length === 0) {
  // Error: Không có role
} else if (roles.length === 1) {
  // Tự động chọn role duy nhất
  setActiveRole(roles[0]);
} else {
  // User chọn role
  const selectedRole = await showRoleSelector(roles);
  setActiveRole(selectedRole);
}

function setActiveRole(role) {
  localStorage.setItem('activeRoleId', role.roleId);
  axios.defaults.headers.common['X-Active-Role-ID'] = role.roleId;
}
```

### Axios Interceptor

```javascript
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwt_token');
  const roleId = localStorage.getItem('activeRoleId');
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  // Optional: Backend tự động detect role đầu tiên nếu không gửi
  // Nhưng khuyến nghị nên gửi để user có thể chọn role
  if (roleId) {
    config.headers['X-Active-Role-ID'] = roleId;
  }
  
  return config;
});
```

### Backend Tự Động Detect

**Nếu frontend KHÔNG gửi header `X-Active-Role-ID`:**
- ✅ Backend tự động lấy role đầu tiên của user
- ⚠️ User không thể chọn role nếu có nhiều roles
- ✅ Hữu ích cho trường hợp user chỉ có 1 role

## ⚠️ Common Mistakes

❌ **SAI:** Gửi organization ID trong header
```javascript
headers: { 'X-Active-Organization-ID': orgId } // SAI!
```

✅ **ĐÚNG:** Gửi role ID trong header
```javascript
headers: { 'X-Active-Role-ID': roleId } // ĐÚNG!
```

## 📝 Checklist

- [ ] Setup axios interceptor
- [ ] Gọi `/auth/roles` sau khi login
- [ ] Chọn role (tự động hoặc cho user chọn)
- [ ] Lưu `activeRoleId` vào localStorage
- [ ] Gửi `X-Active-Role-ID` trong mọi request
- [ ] Implement switch role function
- [ ] Reload data sau khi đổi role
