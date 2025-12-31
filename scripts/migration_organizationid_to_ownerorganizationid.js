// Migration Script: Đổi tên field organizationId → ownerOrganizationId
// Chạy script này trong MongoDB shell hoặc MongoDB Compass
// Usage: mongo <database_name> migration_organizationid_to_ownerorganizationid.js

print("🚀 Bắt đầu migration: organizationId → ownerOrganizationId");
print("==========================================");

// Danh sách collections cần migrate
const collections = [
    // Notification collections
    "notification_senders",
    "notification_templates",
    "notification_channels",
    "notification_queue",
    "notification_history",
    
    // Auth collections
    "roles",
    "auth_logs",
    
    // Facebook collections
    "fb_posts",
    "fb_conversations",
    "fb_messages",
    "fb_message_items",
    "fb_pages",
    "fb_customers",
    
    // Pancake POS collections
    "pc_pos_orders",
    "pc_pos_products",
    "pc_pos_shops",
    "pc_pos_customers",
    "pc_pos_warehouses",
    "pc_pos_variations",
    "pc_pos_categories",
    
    // Other collections
    "customers",
    "access_tokens"
];

let totalUpdated = 0;
let totalErrors = 0;

collections.forEach(collectionName => {
    try {
        const collection = db.getCollection(collectionName);
        
        // Kiểm tra collection có tồn tại không
        if (!collection.exists()) {
            print(`⚠️  Collection "${collectionName}" không tồn tại, bỏ qua...`);
            return;
        }
        
        // Đếm số documents có field organizationId
        const count = collection.countDocuments({ organizationId: { $exists: true } });
        
        if (count === 0) {
            print(`✅ Collection "${collectionName}": Không có documents cần migrate (0 documents)`);
            return;
        }
        
        print(`\n📦 Migrating collection: "${collectionName}" (${count} documents)...`);
        
        // Migration: Đổi tên field organizationId → ownerOrganizationId
        const result = collection.updateMany(
            { organizationId: { $exists: true } },
            [
                {
                    $set: { ownerOrganizationId: "$organizationId" }
                },
                {
                    $unset: "organizationId"
                }
            ]
        );
        
        if (result.modifiedCount > 0) {
            print(`   ✅ Đã migrate ${result.modifiedCount} documents`);
            totalUpdated += result.modifiedCount;
        } else {
            print(`   ⚠️  Không có documents nào được migrate`);
        }
        
        // Verify: Kiểm tra không còn field organizationId
        const remainingCount = collection.countDocuments({ organizationId: { $exists: true } });
        if (remainingCount > 0) {
            print(`   ⚠️  CẢNH BÁO: Vẫn còn ${remainingCount} documents có field organizationId!`);
        }
        
        // Verify: Kiểm tra có field ownerOrganizationId
        const newCount = collection.countDocuments({ ownerOrganizationId: { $exists: true } });
        print(`   📊 Documents có ownerOrganizationId: ${newCount}`);
        
    } catch (error) {
        print(`   ❌ LỖI khi migrate collection "${collectionName}": ${error.message}`);
        totalErrors++;
    }
});

print("\n==========================================");
print("✅ Migration hoàn tất!");
print(`📊 Tổng số documents đã migrate: ${totalUpdated}`);
if (totalErrors > 0) {
    print(`❌ Số lỗi: ${totalErrors}`);
}
print("\n⚠️  LƯU Ý QUAN TRỌNG:");
print("1. Cần tạo lại indexes với tên field mới (ownerOrganizationId)");
print("2. Xóa indexes cũ với tên field cũ (organizationId) nếu có");
print("3. Verify dữ liệu sau khi migration");
print("4. Test các API endpoints để đảm bảo hoạt động đúng");
