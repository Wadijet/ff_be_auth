# ETL Pipeline Project Brief

## Mục tiêu
Xây dựng hệ thống ETL pipeline với khả năng cấu hình cao, giảm thiểu việc code và tái sử dụng tối đa.

## Kiến trúc tổng quan
1. **Registry Pattern**
   - Tận dụng registry pattern hiện có
   - Mở rộng để quản lý ETL components

2. **Components chính**
   - DataSource: Quản lý nguồn dữ liệu (chủ yếu là REST API)
   - Transformer: Xử lý và biến đổi dữ liệu
   - Destination: Điểm đích của dữ liệu (Internal API)
   - Pipeline: Workflow engine quản lý quy trình ETL

3. **Cấu hình**
   - Sử dụng YAML để định nghĩa components
   - Cấu trúc modular, dễ mở rộng
   - Tái sử dụng các định nghĩa

## Yêu cầu kỹ thuật
1. **Performance**
   - Hỗ trợ xử lý realtime
   - Tối ưu hiệu năng xử lý

2. **Maintainability**
   - Code tối thiểu, cấu hình tối đa
   - Dễ dàng thêm mới components

3. **Reliability**
   - Error handling
   - Monitoring và logging

## Phạm vi phase 1
1. Cấu trúc cơ bản của hệ thống
2. Registry cho ETL components
3. Định nghĩa và xử lý cấu hình
4. REST API connector cơ bản
5. Field mapping transformer
6. Internal API destination 