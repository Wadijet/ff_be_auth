# Git Workflow

Quy trình làm việc với Git trong dự án.

## 📋 Tổng Quan

Tài liệu này mô tả quy trình làm việc với Git.

## 🌿 Branch Strategy

### Main Branches

- `main`: Production code
- `develop`: Development code

### Feature Branches

- `feature/<feature-name>`: Feature development
- `bugfix/<bug-name>`: Bug fixes
- `hotfix/<hotfix-name>`: Hotfixes

## 🔄 Workflow

### 1. Tạo Feature Branch

```bash
git checkout develop
git pull origin develop
git checkout -b feature/new-feature
```

### 2. Commit Changes

```bash
git add .
git commit -m "feat: add new feature"
```

### 3. Push và Tạo Pull Request

```bash
git push origin feature/new-feature
```

Tạo Pull Request từ `feature/new-feature` vào `develop`.

### 4. Merge

Sau khi review và approve, merge PR.

## 📝 Commit Messages

Format: `<type>: <message>`

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style
- `refactor`: Refactoring
- `test`: Tests
- `chore`: Maintenance

**Ví dụ:**
```
feat: add user authentication
fix: resolve login issue
docs: update API documentation
```

## 📚 Tài Liệu Liên Quan

- [Coding Standards](coding-standards.md)

