package constants

// Request Status
const (
	RequestStatusUnknown   = 0
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


const (
	WarehouseRequestTypeImport   = "IMPORT"
	WarehouseRequestTypeExport   = "EXPORT"
	WarehouseRequestTypeTransfer = "TRANSFER"
)

func IsValidWarehouseRequestType(requestType string) bool {
	switch requestType {
	case WarehouseRequestTypeImport, WarehouseRequestTypeExport, WarehouseRequestTypeTransfer:
		return true
	default:
		return false
	}
}

func GetWarehouseRequestTypeNames() []string {
	return []string{
		WarehouseRequestTypeImport,
		WarehouseRequestTypeExport,
		WarehouseRequestTypeTransfer,
	}
}

const (
	OtherRequestTypeAccessPermission = "ACCESS_PERMISSION"
	OtherRequestTypeLeaveRequest     = "LEAVE_REQUEST"
	OtherRequestTypeChangeInfo       = "CHANGE_INFO"
	OtherRequestTypeOther            = "OTHER"
)

func IsValidOtherRequestType(requestType string) bool {
	switch requestType {
	case OtherRequestTypeAccessPermission, OtherRequestTypeLeaveRequest, OtherRequestTypeChangeInfo, OtherRequestTypeOther:
		return true
	default:
		return false
	}
}

func GetOtherRequestTypeNames() []string {
	return []string{
		OtherRequestTypeAccessPermission,
		OtherRequestTypeLeaveRequest,
		OtherRequestTypeChangeInfo,
		OtherRequestTypeOther,
	}
}

func IsValidRequestType(requestType string) bool {
	return IsValidWarehouseRequestType(requestType) || IsValidOtherRequestType(requestType)
}

func GetAllRequestTypeNames() []string {
	all := []string{}
	all = append(all, GetWarehouseRequestTypeNames()...)
	all = append(all, GetOtherRequestTypeNames()...)
	return all
}
