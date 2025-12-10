# Cấu Hình Server

Hướng dẫn cấu hình server cho production.

## 📋 Tổng Quan

Tài liệu này hướng dẫn cách cấu hình server (Nginx, Caddy) để reverse proxy cho ứng dụng.

## 🌐 Nginx Configuration

### Basic Configuration

Tạo file `/etc/nginx/sites-available/folkform-auth`:

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Proxy to application
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### Enable Site

```bash
sudo ln -s /etc/nginx/sites-available/folkform-auth /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🚀 Caddy Configuration

Tạo file `Caddyfile`:

```
api.yourdomain.com {
    reverse_proxy localhost:8080
}
```

### Run Caddy

```bash
caddy run
```

## 🔒 SSL Certificate

### Let's Encrypt với Certbot

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal
sudo certbot renew --dry-run
```

## 📊 Rate Limiting

### Nginx Rate Limiting

```nginx
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

server {
    location / {
        limit_req zone=api_limit burst=20;
        proxy_pass http://localhost:8080;
    }
}
```

## 📝 Lưu Ý

- Sử dụng HTTPS trong production
- Cấu hình rate limiting phù hợp
- Kiểm tra logs thường xuyên
- Backup configuration files

## 📚 Tài Liệu Liên Quan

- [Triển Khai Production](production.md)
- [Systemd Service](systemd.md)

