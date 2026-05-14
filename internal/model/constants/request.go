package constants

// Request Status
const (
	RequestStatusUnknown   = 0
	RequestStatusPending   = 1
	RequestStatusApproved  = 2
	RequestStatusRejected  = 3
	RequestStatusCompleted = 4
	RequestStatusCancelled = 5
	RequestStatusExpired = 6
)

func IsValidRequestStatus(status int) bool {
	switch status {
	case RequestStatusUnknown, RequestStatusPending, RequestStatusApproved, RequestStatusRejected, RequestStatusCompleted, RequestStatusCancelled:
		return true
	default:
		return false
	}
}

func IsValidApprovalStatus(status int) bool {
	switch status {
	case RequestStatusApproved, RequestStatusRejected:
		return true
	default:
		return false
	}
}


const (
	RequestTypeImport   = 1
	RequestTypeExport   = 2
	RequestTypeTransfer = 3
	RequestTypeAccessPermission = 4
)

func IsValidWarehouseRequestType(requestType int) bool {
	switch requestType {
	case RequestTypeImport, RequestTypeExport, RequestTypeTransfer,RequestTypeAccessPermission:
		return true
	default:
		return false
	}
}