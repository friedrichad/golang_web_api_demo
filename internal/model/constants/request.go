package constants

// Request Status
const (
	RequestStatusUnknown = 0
	RequestStatusPending   = 1
	RequestStatusApproved  = 2
	RequestStatusRejected  = 3
	RequestStatusCompleted = 4
	RequestStatusCancelled = 5
)
func IsValidRequestStatus(status int) bool {	
	switch status {
	case RequestStatusUnknown, RequestStatusPending, RequestStatusApproved, RequestStatusRejected, RequestStatusCompleted, RequestStatusCancelled:
		return true
	default:
		return false
	}
}
// Request Types
const (
	RequestTypeImport = 1
	RequestTypeExport = 2
	RequestTypeTransfer = 3
)

func IsValidRequestType(requestType int) bool {
	switch requestType {
	case RequestTypeImport, RequestTypeExport, RequestTypeTransfer:
		return true
	default:
		return false
	}
}
