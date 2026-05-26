package constants

// Inventory Ledger ReferenceType - identifies what triggered the ledger entry
const (
	LedgerReferenceTypeRequest    = 1 // From Request (Import/Export/Transfer)
	LedgerReferenceTypeAdjustment = 2 // From Inventory Adjustment
	LedgerReferenceTypeAudit      = 3 // From Inventory Audit
)

func IsValidLedgerReferenceType(referenceType int) bool {
	switch referenceType {
	case LedgerReferenceTypeRequest, LedgerReferenceTypeAdjustment, LedgerReferenceTypeAudit:
		return true
	default:
		return false
	}
}

// Get reference type name for display
func GetLedgerReferenceTypeName(referenceType int) string {
	switch referenceType {
	case LedgerReferenceTypeRequest:
		return "Request"
	case LedgerReferenceTypeAdjustment:
		return "Adjustment"
	case LedgerReferenceTypeAudit:
		return "Audit"
	default:
		return "Unknown"
	}
}
