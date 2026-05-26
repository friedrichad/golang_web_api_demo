package constants

// Inventory Audit Status
const (
	InventoryAuditStatusUnknown   = 0
	InventoryAuditStatusPending   = 1
	InventoryAuditStatusApproved  = 2
	InventoryAuditStatusRejected  = 3
	InventoryAuditStatusApplied   = 4
	InventoryAuditStatusCancelled = 5
)

func IsValidInventoryAuditStatus(status int) bool {
	switch status {
	case InventoryAuditStatusUnknown, InventoryAuditStatusPending, InventoryAuditStatusApproved, InventoryAuditStatusRejected, InventoryAuditStatusApplied, InventoryAuditStatusCancelled:
		return true
	default:
		return false
	}
}
