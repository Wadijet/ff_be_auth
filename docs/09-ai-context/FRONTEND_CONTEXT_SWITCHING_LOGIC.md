# Logic Chọn Context Làm Việc - Hướng Dẫn Frontend

## 📋 Tổng Quan

Hệ thống sử dụng **Context Switching** để quản lý quyền truy cập dữ liệu theo organization. 

**QUAN TRỌNG:** 
- Context làm việc = **ROLE** (không phải organization)
- Frontend gửi **ROLE ID** trong header `X-Active-Role-ID`
- Backend tự động suy ra organization từ role

## 🔄 Flow Đầy Đủ

### Bước 1: User Đăng Nhập

**Endpoint:** `POST /api/v1/auth/login/firebase`

**Request:**
```json
{
  "idToken": "firebase-id-token",
  "hwid": "hardware-id-optional"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "User Name",
    "token": "jwt-token-here"
  }
}
```

**Lưu token:**
```javascript
localStorage.setItem('jwt_token', response.data.token);
localStorage.setItem('user', JSON.stringify(response.data));
```

---

### Bước 2: Lấy Danh Sách Roles (Context Làm Việc)

**Endpoint:** `GET /api/v1/auth/roles`

**Headers:**
```
Authorization: Bearer <jwt-token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "roleId": "role-id-1",
      "roleName": "Manager",
      "organizationId": "org-id-1",
      "organizationName": "Company A",
      "organizationCode": "COMPANY_A",
      "organizationType": "company",
      "organizationLevel": 1
    },
    {
      "roleId": "role-id-2",
      "roleName": "Employee",
      "organizationId": "org-id-2",
      "organizationName": "Company B",
      "organizationCode": "COMPANY_B",
      "organizationType": "company",
      "organizationLevel": 1
    }
  ]
}
```

**Lưu ý:**
- Response chỉ chứa các role **trực tiếp** của user
- KHÔNG bao gồm children/parents organizations
- Mỗi role = một context làm việc

---

### Bước 3: Chọn Context (Role)

**Logic chọn:**

```javascript
// Lấy danh sách roles
const roles = await api.get('/auth/roles');

if (roles.length === 0) {
  // User không có role nào
  // Hiển thị thông báo lỗi
  showError('Bạn chưa được gán role nào. Vui lòng liên hệ admin.');
} else if (roles.length === 1) {
  // Chỉ có 1 role → Tự động chọn
  const selectedRole = roles[0];
  setActiveContext(selectedRole);
} else {
  // Có nhiều roles → User phải chọn
  const selectedRole = await showRoleSelector(roles);
  setActiveContext(selectedRole);
}
```

**Function setActiveContext:**
```javascript
function setActiveContext(role) {
  // Lưu ROLE ID (không phải organization ID)
  localStorage.setItem('activeRoleId', role.roleId);
  localStorage.setItem('activeOrganizationId', role.organizationId); // Lưu để hiển thị, không gửi trong header
  localStorage.setItem('activeRoleName', role.roleName);
  localStorage.setItem('activeOrganizationName', role.organizationName);
  
  // Cập nhật header cho các request tiếp theo
  axios.defaults.headers.common['X-Active-Role-ID'] = role.roleId;
  
  // Reload data với context mới
  reloadApplicationData();
}
```

---

### Bước 4: Mỗi Request Gửi Kèm Context

**QUAN TRỌNG:** Gửi **ROLE ID** trong header, không phải organization ID

**Headers cho mọi request:**
```
Authorization: Bearer <jwt-token>
X-Active-Role-ID: <role-id>  ← Optional: Backend tự động detect nếu không có
```

**Lưu ý:** 
- ✅ **Khuyến nghị:** Frontend nên gửi header `X-Active-Role-ID` để user có thể chọn role
- ✅ **Tự động:** Nếu không gửi header, backend sẽ tự động dùng role đầu tiên của user
- ⚠️ **Hạn chế:** Nếu không gửi header và user có nhiều roles, backend sẽ dùng role đầu tiên (user không thể chọn)

**Ví dụ với Axios:**
```javascript
// Setup interceptor để tự động thêm header
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('jwt_token');
  const activeRoleId = localStorage.getItem('activeRoleId');
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  if (activeRoleId) {
    config.headers['X-Active-Role-ID'] = activeRoleId;
  }
  
  return config;
});
```

---

### Bước 5: Đổi Context (Switch Role)

**Khi user muốn đổi role:**

```javascript
async function switchContext(newRoleId) {
  // Validate role mới có trong danh sách roles của user không
  const roles = await api.get('/auth/roles');
  const newRole = roles.find(r => r.roleId === newRoleId);
  
  if (!newRole) {
    showError('Role không hợp lệ');
    return;
  }
  
  // Cập nhật context
  setActiveContext(newRole);
  
  // Reload toàn bộ data với context mới
  window.location.reload(); // Hoặc update state nếu dùng React/Vue
}
```

---

## 📝 Implementation Example (React)

```javascript
// Context/Store
const [activeRole, setActiveRole] = useState(null);
const [userRoles, setUserRoles] = useState([]);

// Sau khi login
useEffect(() => {
  const loadUserRoles = async () => {
    try {
      const response = await api.get('/auth/roles');
      const roles = response.data;
      
      setUserRoles(roles);
      
      // Tự động chọn role nếu chỉ có 1
      if (roles.length === 1) {
        setActiveRole(roles[0]);
        localStorage.setItem('activeRoleId', roles[0].roleId);
        axios.defaults.headers.common['X-Active-Role-ID'] = roles[0].roleId;
      } else if (roles.length > 1) {
        // Hiển thị dialog cho user chọn
        // Hoặc lấy từ localStorage nếu đã chọn trước đó
        const savedRoleId = localStorage.getItem('activeRoleId');
        if (savedRoleId) {
          const savedRole = roles.find(r => r.roleId === savedRoleId);
          if (savedRole) {
            setActiveRole(savedRole);
            axios.defaults.headers.common['X-Active-Role-ID'] = savedRoleId;
          }
        }
      }
    } catch (error) {
      console.error('Failed to load roles:', error);
    }
  };
  
  loadUserRoles();
}, []);

// Component hiển thị role selector
function RoleSelector({ roles, onSelect }) {
  return (
    <div className="role-selector">
      <h3>Chọn context làm việc:</h3>
      {roles.map(role => (
        <button 
          key={role.roleId}
          onClick={() => onSelect(role)}
        >
          {role.roleName} - {role.organizationName}
        </button>
      ))}
    </div>
  );
}

// Function đổi context
const handleSwitchRole = (newRole) => {
  setActiveRole(newRole);
  localStorage.setItem('activeRoleId', newRole.roleId);
  axios.defaults.headers.common['X-Active-Role-ID'] = newRole.roleId;
  
  // Reload data
  window.location.reload();
};
```

---

## 📝 Implementation Example (Vue)

```javascript
// Store/Pinia
export const useAuthStore = defineStore('auth', {
  state: () => ({
    activeRole: null,
    userRoles: [],
  }),
  
  actions: {
    async loadUserRoles() {
      try {
        const response = await api.get('/auth/roles');
        this.userRoles = response.data;
        
        if (this.userRoles.length === 1) {
          this.setActiveRole(this.userRoles[0]);
        } else if (this.userRoles.length > 1) {
          const savedRoleId = localStorage.getItem('activeRoleId');
          if (savedRoleId) {
            const savedRole = this.userRoles.find(r => r.roleId === savedRoleId);
            if (savedRole) {
              this.setActiveRole(savedRole);
            }
          }
        }
      } catch (error) {
        console.error('Failed to load roles:', error);
      }
    },
    
    setActiveRole(role) {
      this.activeRole = role;
      localStorage.setItem('activeRoleId', role.roleId);
      axios.defaults.headers.common['X-Active-Role-ID'] = role.roleId;
    },
    
    switchRole(newRoleId) {
      const newRole = this.userRoles.find(r => r.roleId === newRoleId);
      if (newRole) {
        this.setActiveRole(newRole);
        window.location.reload();
      }
    }
  }
});
```

---

## ⚠️ Lưu Ý Quan Trọng

### 1. Context Là ROLE, Không Phải Organization

✅ **ĐÚNG:**
```javascript
// Gửi ROLE ID trong header
headers: {
  'X-Active-Role-ID': 'role-id-123'
}
```

❌ **SAI:**
```javascript
// KHÔNG gửi organization ID
headers: {
  'X-Active-Organization-ID': 'org-id-123' // SAI!
}
```

### 2. Response Format

Endpoint `/auth/roles` trả về:
- ✅ Chỉ các role **trực tiếp** của user
- ✅ Organization **trực tiếp** của mỗi role
- ❌ KHÔNG bao gồm children organizations
- ❌ KHÔNG bao gồm parent organizations
- ❌ KHÔNG có tree structure

### 3. Khi Nào Cần Reload?

- ✅ Sau khi đổi role (switch context)
- ✅ Sau khi login lần đầu
- ❌ KHÔNG cần reload khi chỉ xem data

### 4. Error Handling

```javascript
// Nếu không có role
if (roles.length === 0) {
  showError('Bạn chưa được gán role nào. Vui lòng liên hệ admin.');
  // Redirect về trang login hoặc hiển thị thông báo
}

// Nếu role không hợp lệ (backend trả về 403)
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 403) {
      // Role không hợp lệ, reload danh sách roles
      loadUserRoles();
    }
    return Promise.reject(error);
  }
);
```

---

## 🔍 API Endpoints Summary

| Endpoint | Method | Headers | Mô Tả |
|----------|--------|---------|-------|
| `/auth/login/firebase` | POST | - | Đăng nhập với Firebase |
| `/auth/roles` | GET | `Authorization` | Lấy danh sách roles (context làm việc) |
| Tất cả endpoints khác | * | `Authorization`, `X-Active-Role-ID` | Gửi kèm context trong mọi request |

---

## 📊 Data Flow

```
1. Login
   ↓
2. GET /auth/roles
   ↓
3. User chọn role (hoặc tự động nếu chỉ có 1)
   ↓
4. Lưu roleId vào localStorage
   ↓
5. Set header X-Active-Role-ID cho mọi request
   ↓
6. Backend tự động suy ra organization từ role
   ↓
7. Tất cả data operations dùng organization đó
```

---

## 🎯 Checklist Implementation

- [ ] Setup axios interceptor để tự động thêm `X-Active-Role-ID` header
- [ ] Sau khi login, gọi `GET /auth/roles` để lấy danh sách roles
- [ ] Nếu có 1 role → Tự động chọn
- [ ] Nếu có nhiều roles → Hiển thị selector cho user chọn
- [ ] Lưu `activeRoleId` vào localStorage
- [ ] Hiển thị role hiện tại ở UI (header/sidebar)
- [ ] Implement chức năng đổi role (switch context)
- [ ] Reload data sau khi đổi role
- [ ] Handle error khi không có role hoặc role không hợp lệ

---

## 💡 Tips

1. **Lưu role vào localStorage** để giữ context khi refresh page
2. **Validate role** khi load lại từ localStorage (có thể role đã bị xóa)
3. **Hiển thị role hiện tại** ở UI để user biết đang làm việc với context nào
4. **Reload data** sau khi đổi role để đảm bảo data đúng với context mới
5. **Error handling** khi role không hợp lệ hoặc không có role
