package constants

// Inventory Adjustment Status
const(
	InventoryAdjustmentStatusUnknown = 0
	InventoryAdjustmentStatusPending  = 1
	InventoryAdjustmentStatusApproved = 2
	InventoryAdjustmentStatusRejected = 3
	InventoryAdjustmentStatusApplied = 4
	InventoryAdjustmentStatusCancelled = 5
)
func IsValidInventoryAdjustmentStatus(status int) bool {
	switch status {
	case InventoryAdjustmentStatusUnknown, InventoryAdjustmentStatusPending, InventoryAdjustmentStatusApproved, InventoryAdjustmentStatusRejected, InventoryAdjustmentStatusApplied, InventoryAdjustmentStatusCancelled:
		return true
	default:
		return false
	}
}

