package dtos

type Finishedtransferrequestcomponent struct {
	FinishedTransferRequestComponentID int32   `json:"finished_transfer_request_component_id"`
	RequestID                          int32   `json:"request_id"`
	ComponentID                        int32   `json:"component_id"`
	BinID                              int32   `json:"bin_id"`
	Quantity                           float64 `json:"quantity"`
	TypeInt                            int32   `json:"type_int"`
}