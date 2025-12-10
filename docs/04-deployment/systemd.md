# Systemd Service

Hướng dẫn cấu hình systemd service cho hệ thống.

## 📋 Tổng Quan

Tài liệu này hướng dẫn cách tạo systemd service để chạy ứng dụng như một service trên Linux.

## 📝 Tạo Service File

Tạo file `/etc/systemd/system/folkform-auth.service`:

```ini
[Unit]
Description=FolkForm Auth Backend
After=network.target mongodb.service

[Service]
Type=simple
User=folkform
WorkingDirectory=/opt/folkform-auth/api
ExecStart=/opt/folkform-auth/api/server
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=folkform-auth

# Environment variables
Environment="ADDRESS=8080"
Environment="MONGODB_CONNECTION_URI=mongodb://localhost:27017"
Environment="JWT_SECRET=your-secret-key"

[Install]
WantedBy=multi-user.target
```

## 🔧 Cấu Hình

### 1. Tạo User

```bash
sudo useradd -r -s /bin/false folkform
```

### 2. Copy Files

```bash
sudo mkdir -p /opt/folkform-auth
sudo cp -r api /opt/folkform-auth/
sudo chown -R folkform:folkform /opt/folkform-auth
```

### 3. Reload Systemd

```bash
sudo systemctl daemon-reload
```

### 4. Enable Service

```bash
sudo systemctl enable folkform-auth.service
```

### 5. Start Service

```bash
sudo systemctl start folkform-auth.service
```

## 🔍 Quản Lý Service

### Kiểm Tra Status

```bash
sudo systemctl status folkform-auth.service
```

### Xem Logs

```bash
# Xem logs
sudo journalctl -u folkform-auth.service

# Xem logs real-time
sudo journalctl -u folkform-auth.service -f

# Xem logs của ngày hôm nay
sudo journalctl -u folkform-auth.service --since today
```

### Restart Service

```bash
sudo systemctl restart folkform-auth.service
```

### Stop Service

```bash
sudo systemctl stop folkform-auth.service
```

### Disable Service

```bash
sudo systemctl disable folkform-auth.service
```

## 📝 Lưu Ý

- Đảm bảo user có quyền truy cập các file cần thiết
- Cấu hình environment variables trong service file hoặc file riêng
- Sử dụng `Restart=always` để tự động restart khi crash
- Kiểm tra logs thường xuyên để phát hiện lỗi

## 📚 Tài Liệu Liên Quan

- [Triển Khai Production](production.md)
- [Cấu Hình Server](cau-hinh-server.md)

