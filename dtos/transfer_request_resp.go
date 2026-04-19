package dtos
import "time"

type TransferRequestResponse struct {
	RequestID       int32     `json:"request_id"`
	Description     string    `json:"description"`
	TypeInt         int32     `json:"type_int"`
	CreationTime    time.Time `json:"creation_time"`
	ExecutionTime   time.Time `json:"execution_time"`
	StatusInt       int32     `json:"status_int"`
	CreatorID       int32     `json:"creator_id"`
	ApproverID      int32     `json:"approver_id"`
	WarehouseFromID int32     `json:"warehouse_from_id"`
	WarehouseToID   int32     `json:"warehouse_to_id"`
	CustomerID      int32     `json:"customer_id"`
}