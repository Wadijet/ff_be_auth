# Báo Cáo Test

Hướng dẫn xem và phân tích báo cáo test.

## 📋 Tổng Quan

Sau khi chạy test, báo cáo tự động được tạo trong `api-tests/reports/`.

## 📊 Format Báo Cáo

Báo cáo được lưu dưới dạng Markdown với tên file:
```
test_report_YYYY-MM-DD_HH-MM-SS.md
```

## 📝 Nội Dung Báo Cáo

Báo cáo bao gồm:
- Tổng số test cases
- Số test passed
- Số test failed
- Pass rate (%)
- Chi tiết từng test case

## 🔍 Xem Báo Cáo

### PowerShell

```powershell
# Mở file report mới nhất
Get-ChildItem api-tests\reports\*.md | Sort-Object LastWriteTime -Descending | Select-Object -First 1 | ForEach-Object { notepad $_.FullName }
```

### Command Line

```bash
# Xem report mới nhất
ls -t api-tests/reports/*.md | head -1 | xargs cat
```

## 📈 Phân Tích

- **Pass Rate > 95%**: Tốt
- **Pass Rate 80-95%**: Cần cải thiện
- **Pass Rate < 80%**: Có vấn đề nghiêm trọng

## 📚 Tài Liệu Liên Quan

- [Tổng Quan Testing](tong-quan.md)
- [Chạy Test Suite](chay-test.md)
- [Viết Test Case](viet-test.md)

