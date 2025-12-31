// Migration Script: Tạo lại indexes cho field ownerOrganizationId
// Chạy script này SAU KHI đã chạy migration_organizationid_to_ownerorganizationid.js
// Usage: mongo <database_name> migration_recreate_indexes.js

print("🚀 Bắt đầu tạo lại indexes cho ownerOrganizationId");
print("==========================================");

// Danh sách collections và indexes cần tạo
const collectionsWithIndexes = [
    {
        collection: "notification_senders",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" },
            { key: { ownerOrganizationId: 1, channelType: 1, name: 1 }, name: "ownerOrganizationId_1_channelType_1_name_1" }
        ]
    },
    {
        collection: "notification_templates",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" },
            { key: { ownerOrganizationId: 1, eventType: 1, channelType: 1 }, name: "ownerOrganizationId_1_eventType_1_channelType_1" }
        ]
    },
    {
        collection: "notification_channels",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "notification_queue",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "notification_history",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "roles",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" },
            { key: { ownerOrganizationId: 1, name: 1 }, name: "role_org_name_unique" }
        ]
    },
    {
        collection: "auth_logs",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_posts",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_conversations",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_messages",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_message_items",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_pages",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "fb_customers",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_orders",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_products",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_shops",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_customers",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_warehouses",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_variations",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "pc_pos_categories",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "customers",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    },
    {
        collection: "access_tokens",
        indexes: [
            { key: { ownerOrganizationId: 1 }, name: "ownerOrganizationId_1" }
        ]
    }
];

let totalCreated = 0;
let totalErrors = 0;

collectionsWithIndexes.forEach(({ collection: collectionName, indexes }) => {
    try {
        const collection = db.getCollection(collectionName);
        
        // Kiểm tra collection có tồn tại không
        if (!collection.exists()) {
            print(`⚠️  Collection "${collectionName}" không tồn tại, bỏ qua...`);
            return;
        }
        
        print(`\n📦 Tạo indexes cho collection: "${collectionName}"...`);
        
        indexes.forEach(index => {
            try {
                // Xóa index cũ nếu có (với tên tương tự)
                const oldIndexName = index.name.replace("ownerOrganizationId", "organizationId");
                try {
                    collection.dropIndex(oldIndexName);
                    print(`   🗑️  Đã xóa index cũ: ${oldIndexName}`);
                } catch (e) {
                    // Index cũ không tồn tại, bỏ qua
                }
                
                // Tạo index mới
                collection.createIndex(index.key, { name: index.name, background: true });
                print(`   ✅ Đã tạo index: ${index.name}`);
                totalCreated++;
            } catch (error) {
                if (error.code === 85) {
                    // Index đã tồn tại
                    print(`   ℹ️  Index "${index.name}" đã tồn tại, bỏ qua...`);
                } else {
                    print(`   ❌ Lỗi khi tạo index "${index.name}": ${error.message}`);
                    totalErrors++;
                }
            }
        });
        
    } catch (error) {
        print(`   ❌ LỖI khi xử lý collection "${collectionName}": ${error.message}`);
        totalErrors++;
    }
});

print("\n==========================================");
print("✅ Tạo indexes hoàn tất!");
print(`📊 Tổng số indexes đã tạo: ${totalCreated}`);
if (totalErrors > 0) {
    print(`❌ Số lỗi: ${totalErrors}`);
}
print("\n⚠️  LƯU Ý:");
print("1. Verify indexes đã được tạo đúng");
print("2. Kiểm tra performance sau khi migration");
