package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Customer lưu thông tin khách hàng từ các nguồn (Pancake, POS, ...)
// Hỗ trợ multi-source data: panCakeData (Facebook), posData (POS)
type Customer struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // ID của customer

	// ===== COMMON FIELDS (Extract từ nhiều nguồn với conflict resolution) =====
	// Name: Ưu tiên POS (priority=1) hơn Pancake (priority=2) - vì sale thao tác trên POS
	Name string `json:"name" bson:"name" index:"text" extract:"PosData\\.name,converter=string,optional,priority=1,merge=priority|PanCakeData\\.name,converter=string,optional,priority=2,merge=priority"`

	// PhoneNumbers: Merge từ tất cả nguồn vào array
	// - POS: phone_numbers (array) - ưu tiên
	// - Pancake: phone_numbers (array)
	PhoneNumbers []string `json:"phoneNumbers" bson:"phoneNumbers" index:"text" extract:"PosData\\.phone_numbers,optional,priority=1,merge=merge_array|PanCakeData\\.phone_numbers,optional,priority=2,merge=merge_array"`

	// Email: Ưu tiên POS (priority=1) hơn Pancake (priority=2)
	// - POS: emails (array) → lấy email đầu tiên
	// - Pancake: email (string)
	Email string `json:"email" bson:"email" index:"text" extract:"PosData\\.emails,converter=array_first,optional,priority=1,merge=priority|PanCakeData\\.email,converter=string,optional,priority=2,merge=priority"`

	// ===== COMMON IDENTIFIER =====
	// CustomerId: ID chung để identify customer từ cả 2 nguồn (dùng cho filter khi upsert)
	// - POS: id (UUID string) - ưu tiên
	// - Pancake: id
	CustomerId string `json:"customerId" bson:"customerId" index:"text,unique,sparse" extract:"PosData\\.id,converter=string,optional,priority=1,merge=priority|PanCakeData\\.id,converter=string,optional,priority=2,merge=priority"` // ID chung để identify customer (unique, sparse)

	// ===== SOURCE-SPECIFIC IDENTIFIERS =====
	PanCakeCustomerId string `json:"panCakeCustomerId" bson:"panCakeCustomerId" index:"text" extract:"PanCakeData\\.id,converter=string,optional"` // Pancake Customer ID (extract từ PanCakeData["id"])
	Psid              string `json:"psid" bson:"psid" index:"text" extract:"PanCakeData\\.psid,converter=string,optional"`                         // PSID từ Pancake (Page Scoped ID, extract từ PanCakeData["psid"])
	PageId            string `json:"pageId" bson:"pageId" index:"text" extract:"PanCakeData\\.page_id,converter=string,optional"`                  // Page ID từ Pancake (extract từ PanCakeData["page_id"])

	PosCustomerId string `json:"posCustomerId" bson:"posCustomerId" index:"text,unique,sparse" extract:"PosData\\.id,converter=string,optional"` // UUID string - ID của hệ thống POS (unique, sparse)

	// ===== SOURCE-SPECIFIC DATA =====
	PanCakeData map[string]interface{} `json:"panCakeData,omitempty" bson:"panCakeData,omitempty"` // Dữ liệu gốc từ Pancake API
	PosData     map[string]interface{} `json:"posData,omitempty" bson:"posData,omitempty"`         // Dữ liệu gốc từ POS API

	// ===== EXTRACTED FIELDS (Từ các nguồn) =====
	// Common fields có thể có từ cả 2 nguồn - ưu tiên POS (priority=1)
	Birthday string `json:"birthday,omitempty" bson:"birthday,omitempty" extract:"PosData\\.date_of_birth,converter=string,optional,priority=1,merge=priority|PanCakeData\\.birthday,converter=string,optional,priority=2,merge=priority"` // Ngày sinh
	Gender   string `json:"gender,omitempty" bson:"gender,omitempty" extract:"PosData\\.gender,converter=string,optional,priority=1,merge=priority|PanCakeData\\.gender,converter=string,optional,priority=2,merge=priority"`              // Giới tính

	// Pancake-specific
	LivesIn          string `json:"livesIn,omitempty" bson:"livesIn,omitempty" extract:"PanCakeData\\.lives_in,converter=string,optional,merge=keep_existing"`             // Nơi ở
	PanCakeUpdatedAt int64  `json:"panCakeUpdatedAt" bson:"panCakeUpdatedAt" extract:"PanCakeData\\.updated_at,converter=time,format=2006-01-02T15:04:05.000000,optional"` // Thời gian cập nhật từ Pancake

	// POS-specific
	CustomerLevelId   string        `json:"customerLevelId,omitempty" bson:"customerLevelId,omitempty" extract:"PosData\\.level_id,converter=string,optional,merge=overwrite"`               // UUID string
	Point             int64         `json:"point,omitempty" bson:"point,omitempty" extract:"PosData\\.reward_point,converter=int64,optional,merge=overwrite"`                                // Điểm tích lũy
	TotalOrder        int64         `json:"totalOrder,omitempty" bson:"totalOrder,omitempty" extract:"PosData\\.order_count,converter=int64,optional,merge=overwrite"`                       // Tổng đơn hàng
	TotalSpent        float64       `json:"totalSpent,omitempty" bson:"totalSpent,omitempty" extract:"PosData\\.purchased_amount,converter=number,optional,merge=overwrite"`                 // Tổng tiền đã mua
	SucceedOrderCount int64         `json:"succeedOrderCount,omitempty" bson:"succeedOrderCount,omitempty" extract:"PosData\\.succeed_order_count,converter=int64,optional,merge=overwrite"` // Số đơn hàng thành công
	TagIds            []interface{} `json:"tagIds,omitempty" bson:"tagIds,omitempty" extract:"PosData\\.tags,optional,merge=overwrite"`                                                      // Tags (array)
	PosLastOrderAt    int64         `json:"posLastOrderAt,omitempty" bson:"posLastOrderAt,omitempty" extract:"PosData\\.last_order_at,converter=time,format=2006-01-02T15:04:05Z,optional"`  // Thời gian đơn hàng cuối
	PosAddresses      []interface{} `json:"posAddresses,omitempty" bson:"posAddresses,omitempty" extract:"PosData\\.shop_customer_address,optional,merge=overwrite"`                         // Địa chỉ (array)
	PosReferralCode   string        `json:"posReferralCode,omitempty" bson:"posReferralCode,omitempty" extract:"PosData\\.referral_code,converter=string,optional,merge=overwrite"`          // Mã giới thiệu
	PosIsBlock        bool          `json:"posIsBlock,omitempty" bson:"posIsBlock,omitempty" extract:"PosData\\.is_block,converter=bool,optional,merge=overwrite"`                           // Trạng thái block

	// ===== ORGANIZATION =====
	OwnerOrganizationID primitive.ObjectID `json:"ownerOrganizationId" bson:"ownerOrganizationId" index:"single:1"` // Tổ chức sở hữu dữ liệu (phân quyền)

	// ===== METADATA =====
	Sources   []string `json:"sources" bson:"sources"`     // ["pancake", "pos"] - Track nguồn dữ liệu
	CreatedAt int64    `json:"createdAt" bson:"createdAt"` // Thời gian tạo
	UpdatedAt int64    `json:"updatedAt" bson:"updatedAt"` // Thời gian cập nhật
}
