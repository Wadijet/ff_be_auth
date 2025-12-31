// Migration script: Đổi FromOrgID → OwnerOrganizationID và đổi tên collection
// Từ: organization_shares (field: fromOrgId)
// Sang: auth_organization_shares (field: ownerOrganizationId)

// Lưu ý: Script này giả định bạn sẽ xóa DB và làm lại từ đầu (theo yêu cầu user)
// Nếu cần migration trên DB hiện có, uncomment các dòng migration

print("==========================================");
print("Migration: OrganizationShare");
print("==========================================");
print("1. Copy fromOrgId → ownerOrganizationId");
print("2. Đổi tên collection: organization_shares → auth_organization_shares");
print("==========================================");

// ==========================================
// BƯỚC 1: Copy fromOrgId → ownerOrganizationId
// ==========================================
print("\n[Bước 1] Copy fromOrgId → ownerOrganizationId...");

// Lấy collection cũ
var oldCollection = db.organization_shares;

// Đếm số documents
var count = oldCollection.countDocuments();
print("Số documents cần migrate: " + count);

if (count > 0) {
    // Update tất cả documents: copy fromOrgId → ownerOrganizationId
    var result = oldCollection.updateMany(
        { ownerOrganizationId: { $exists: false } }, // Chỉ update nếu chưa có ownerOrganizationId
        [
            {
                $set: {
                    ownerOrganizationId: "$fromOrgId"
                }
            }
        ]
    );
    
    print("Đã copy " + result.modifiedCount + " documents");
    
    // Kiểm tra kết quả
    var checkResult = oldCollection.find({ ownerOrganizationId: { $exists: true } }).count();
    print("Số documents có ownerOrganizationId: " + checkResult);
} else {
    print("Không có documents nào cần migrate");
}

// ==========================================
// BƯỚC 2: Đổi tên collection
// ==========================================
print("\n[Bước 2] Đổi tên collection: organization_shares → auth_organization_shares...");

// Kiểm tra collection mới đã tồn tại chưa
var newCollectionExists = db.getCollectionNames().includes("auth_organization_shares");

if (newCollectionExists) {
    print("⚠️  Collection auth_organization_shares đã tồn tại!");
    print("   Bạn có muốn xóa và tạo lại không? (Script sẽ không tự động xóa)");
} else {
    // Đổi tên collection
    try {
        oldCollection.renameCollection("auth_organization_shares");
        print("✅ Đã đổi tên collection thành công");
    } catch (e) {
        print("❌ Lỗi khi đổi tên collection: " + e);
        print("   Có thể collection auth_organization_shares đã tồn tại");
    }
}

// ==========================================
// BƯỚC 3: Tạo indexes mới
// ==========================================
print("\n[Bước 3] Tạo indexes cho collection mới...");

var newCollection = db.auth_organization_shares;

// Tạo index cho ownerOrganizationId
try {
    newCollection.createIndex({ ownerOrganizationId: 1 });
    print("✅ Đã tạo index cho ownerOrganizationId");
} catch (e) {
    print("⚠️  Lỗi khi tạo index ownerOrganizationId: " + e);
}

// Tạo index cho toOrgId
try {
    newCollection.createIndex({ toOrgId: 1 });
    print("✅ Đã tạo index cho toOrgId");
} catch (e) {
    print("⚠️  Lỗi khi tạo index toOrgId: " + e);
}

// ==========================================
// BƯỚC 4: Xóa field cũ (tùy chọn - chỉ khi chắc chắn)
// ==========================================
print("\n[Bước 4] Xóa field fromOrgId cũ (tùy chọn)...");
print("⚠️  BƯỚC NÀY SẼ XÓA FIELD fromOrgId - CHỈ CHẠY KHI CHẮC CHẮN!");

// UNCOMMENT DÒNG DƯỚI ĐÂY NẾU MUỐN XÓA FIELD CŨ
// var deleteResult = newCollection.updateMany(
//     {},
//     { $unset: { fromOrgId: "" } }
// );
// print("Đã xóa field fromOrgId từ " + deleteResult.modifiedCount + " documents");

print("\n==========================================");
print("Migration hoàn tất!");
print("==========================================");
print("Collection mới: auth_organization_shares");
print("Field mới: ownerOrganizationId");
print("==========================================");
