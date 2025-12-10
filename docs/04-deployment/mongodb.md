# MongoDB Setup

Hướng dẫn cài đặt và cấu hình MongoDB.

## 📋 Tổng Quan

Hệ thống sử dụng MongoDB để lưu trữ dữ liệu. Tài liệu này hướng dẫn cách setup MongoDB.

## 🚀 Cài Đặt MongoDB

### Windows

1. Tải MongoDB Community Server: https://www.mongodb.com/try/download/community
2. Chạy installer và làm theo hướng dẫn
3. MongoDB sẽ được cài đặt tại `C:\Program Files\MongoDB\Server\<version>\bin`

### Linux (Ubuntu)

```bash
# Import public key
wget -qO - https://www.mongodb.org/static/pgp/server-7.0.asc | sudo apt-key add -

# Add repository
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

# Update và cài đặt
sudo apt-get update
sudo apt-get install -y mongodb-org

# Khởi động MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod
```

### macOS

```bash
# Sử dụng Homebrew
brew tap mongodb/brew
brew install mongodb-community
brew services start mongodb-community
```

### Docker

```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

## ⚙️ Cấu Hình

### Connection String

**Local:**
```
mongodb://localhost:27017
```

**Với Authentication:**
```
mongodb://username:password@localhost:27017
```

**Replica Set:**
```
mongodb://host1:27017,host2:27017/?replicaSet=rs0
```

**Atlas:**
```
mongodb+srv://username:password@cluster.mongodb.net/
```

### Environment Variables

Thêm vào file `.env`:

```env
MONGODB_CONNECTION_URI=mongodb://localhost:27017
MONGODB_DBNAME_AUTH=folkform_auth
MONGODB_DBNAME_STAGING=folkform_staging
MONGODB_DBNAME_DATA=folkform_data
```

## ✅ Kiểm Tra

1. Khởi động MongoDB
2. Kiểm tra kết nối:
```bash
mongosh
# hoặc
mongo
```

3. Khởi động server và kiểm tra log

## 📚 Tài Liệu Liên Quan

- [Cài Đặt và Cấu Hình](../01-getting-started/cai-dat.md)
- [Cấu Hình Môi Trường](../01-getting-started/cau-hinh.md)

