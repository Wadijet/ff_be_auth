# Phân Tích Cấu Trúc Organization Cho Tập Đoàn

## 1. Mô Hình Tập Đoàn

### 1.1. Cấu Trúc Phân Cấp

```
Tập Đoàn (Group/Corporation)
  └── Công Ty (Company)
      └── Phòng Ban (Department)
          └── Bộ Phận (Division)
              └── Team
```

### 1.2. Các Cấp Độ (Level)

| Level | Type | Tên Tiếng Việt | Mô Tả | Parent Type |
|-------|------|----------------|-------|-------------|
| 0 | `group` | Tập Đoàn | Cấp cao nhất, không có parent | `null` |
| 1 | `company` | Công Ty | Thuộc Tập Đoàn | `group` |
| 2 | `department` | Phòng Ban | Thuộc Công Ty | `company` |
| 3 | `division` | Bộ Phận | Thuộc Phòng Ban | `department` |
| 4+ | `team` | Team | Thuộc Bộ Phận hoặc Phòng Ban | `division` hoặc `department` |

## 2. Quy Tắc Validation

### 2.1. Level và Type Validation

```go
// Validate Level và Type phù hợp
func ValidateOrganizationLevelAndType(level int, orgType string, parentID *primitive.ObjectID) error {
    switch level {
    case 0: // Group
        if orgType != "group" {
            return errors.New("Level 0 phải là type 'group'")
        }
        if parentID != nil {
            return errors.New("Level 0 không được có parent")
        }
        
    case 1: // Company
        if orgType != "company" {
            return errors.New("Level 1 phải là type 'company'")
        }
        if parentID == nil {
            return errors.New("Level 1 phải có parent (group)")
        }
        // Validate parent là group
        parent, err := GetOrganization(parentID)
        if err != nil || parent.Type != "group" {
            return errors.New("Parent của Company phải là Group")
        }
        
    case 2: // Department
        if orgType != "department" {
            return errors.New("Level 2 phải là type 'department'")
        }
        if parentID == nil {
            return errors.New("Level 2 phải có parent (company)")
        }
        // Validate parent là company
        parent, err := GetOrganization(parentID)
        if err != nil || parent.Type != "company" {
            return errors.New("Parent của Department phải là Company")
        }
        
    case 3: // Division
        if orgType != "division" {
            return errors.New("Level 3 phải là type 'division'")
        }
        if parentID == nil {
            return errors.New("Level 3 phải có parent (department)")
        }
        // Validate parent là department
        parent, err := GetOrganization(parentID)
        if err != nil || parent.Type != "department" {
            return errors.New("Parent của Division phải là Department")
        }
        
    default: // Team (Level 4+)
        if orgType != "team" {
            return errors.New("Level 4+ phải là type 'team'")
        }
        if parentID == nil {
            return errors.New("Team phải có parent")
        }
        // Validate parent là division hoặc department
        parent, err := GetOrganization(parentID)
        if err != nil {
            return errors.New("Parent không tồn tại")
        }
        if parent.Type != "division" && parent.Type != "department" {
            return errors.New("Parent của Team phải là Division hoặc Department")
        }
    }
    
    return nil
}
```

### 2.2. Path Generation

```go
// GeneratePath tạo Path từ parent và code
func GeneratePath(parentID *primitive.ObjectID, code string) (string, error) {
    if parentID == nil {
        // Root (Group)
        return "/" + code, nil
    }
    
    // Lấy parent
    parent, err := GetOrganization(parentID)
    if err != nil {
        return "", err
    }
    
    // Path = parent.Path + "/" + code
    return parent.Path + "/" + code, nil
}
```

## 3. Ví Dụ Thực Tế

### 3.1. Tập Đoàn ABC

```
Tập Đoàn ABC (Level 0, Type: group)
  ├── Công Ty ABC Việt Nam (Level 1, Type: company)
  │   ├── Phòng Ban Kinh Doanh (Level 2, Type: department)
  │   │   ├── Bộ Phận Bán Hàng (Level 3, Type: division)
  │   │   │   ├── Team Bán Hàng Online (Level 4, Type: team)
  │   │   │   └── Team Bán Hàng Offline (Level 4, Type: team)
  │   │   └── Bộ Phận Chăm Sóc Khách Hàng (Level 3, Type: division)
  │   ├── Phòng Ban Kỹ Thuật (Level 2, Type: department)
  │   │   ├── Team Backend (Level 3, Type: team)
  │   │   ├── Team Frontend (Level 3, Type: team)
  │   │   └── Team DevOps (Level 3, Type: team)
  │   └── Phòng Ban Marketing (Level 2, Type: department)
  │       └── Team Content (Level 3, Type: team)
  ├── Công Ty ABC Singapore (Level 1, Type: company)
  │   └── Phòng Ban Kinh Doanh (Level 2, Type: department)
  └── Công Ty ABC Thái Lan (Level 1, Type: company)
      └── Phòng Ban Kinh Doanh (Level 2, Type: department)
```

### 3.2. Dữ Liệu Mẫu

```json
[
  {
    "id": "org_001",
    "name": "Tập Đoàn ABC",
    "code": "ABC_GROUP",
    "type": "group",
    "parentId": null,
    "path": "/abc_group",
    "level": 0,
    "isActive": true
  },
  {
    "id": "org_002",
    "name": "Công Ty ABC Việt Nam",
    "code": "ABC_VN",
    "type": "company",
    "parentId": "org_001",
    "path": "/abc_group/abc_vn",
    "level": 1,
    "isActive": true
  },
  {
    "id": "org_003",
    "name": "Phòng Ban Kinh Doanh",
    "code": "SALES_DEPT",
    "type": "department",
    "parentId": "org_002",
    "path": "/abc_group/abc_vn/sales_dept",
    "level": 2,
    "isActive": true
  },
  {
    "id": "org_004",
    "name": "Bộ Phận Bán Hàng",
    "code": "SALES_DIV",
    "type": "division",
    "parentId": "org_003",
    "path": "/abc_group/abc_vn/sales_dept/sales_div",
    "level": 3,
    "isActive": true
  },
  {
    "id": "org_005",
    "name": "Team Bán Hàng Online",
    "code": "SALES_ONLINE_TEAM",
    "type": "team",
    "parentId": "org_004",
    "path": "/abc_group/abc_vn/sales_dept/sales_div/sales_online_team",
    "level": 4,
    "isActive": true
  }
]
```

## 4. Query Patterns

### 4.1. Lấy Tất Cả Công Ty (Level 1)

```go
filter := bson.M{
    "type": "company",
    "level": 1,
    "isActive": true,
}
```

### 4.2. Lấy Tất Cả Phòng Ban Của Một Công Ty

```go
companyID := "org_002"
filter := bson.M{
    "type": "department",
    "parentId": utility.String2ObjectID(companyID),
    "isActive": true,
}
```

### 4.3. Lấy Tất Cả Con Của Một Organization (Scope = 1)

```go
parentPath := "/abc_group/abc_vn/sales_dept"
filter := bson.M{
    "path": bson.M{"$regex": "^" + parentPath},
    "isActive": true,
}
```

### 4.4. Lấy Tất Cả Organization Theo Cấp Độ

```go
// Lấy tất cả Level 2 (Phòng Ban)
filter := bson.M{
    "level": 2,
    "isActive": true,
}
```

## 5. Role Assignment

### 5.1. Role Phải Thuộc Organization

- **Bắt buộc**: Mỗi Role phải có `OrganizationID`
- **Không thể**: Role không có OrganizationID
- **Validation**: Khi tạo/cập nhật Role, phải validate OrganizationID tồn tại

### 5.2. Ví Dụ Role Assignment

```
Role: "Trưởng Phòng Kinh Doanh"
  - OrganizationID: "org_003" (Phòng Ban Kinh Doanh)
  - Scope = 1: Có quyền xem dữ liệu của Phòng Ban Kinh Doanh và tất cả con (Bộ Phận Bán Hàng, Team Online, Team Offline)

Role: "Nhân Viên Bán Hàng"
  - OrganizationID: "org_005" (Team Bán Hàng Online)
  - Scope = 2: Chỉ có quyền xem dữ liệu của Team Bán Hàng Online
```

## 6. Best Practices

### 6.1. Code Naming

- **Group**: `ABC_GROUP`, `XYZ_GROUP`
- **Company**: `ABC_VN`, `ABC_SG`, `ABC_TH`
- **Department**: `SALES_DEPT`, `TECH_DEPT`, `MARKETING_DEPT`
- **Division**: `SALES_DIV`, `SUPPORT_DIV`
- **Team**: `SALES_ONLINE_TEAM`, `BACKEND_TEAM`

### 6.2. Path Naming

- Path nên dùng code (không phải ID) để dễ đọc
- Format: `/code_parent/code_child/code_grandchild`
- Ví dụ: `/abc_group/abc_vn/sales_dept`

### 6.3. Level Management

- Không nên có quá 5 cấp độ (0-4)
- Nếu cần nhiều cấp hơn, cân nhắc flatten structure

### 6.4. IsActive Management

- Khi deactivate một Organization, nên deactivate tất cả con
- Hoặc cho phép query với `isActive: true` để tự động filter

