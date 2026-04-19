package dtos

type TransferRequestComponentResponse struct {
	TransferRequestComponentID int32   `json:"transfer_request_component_id"`
	RequestID                  int32   `json:"request_id"`
	ComponentID                int32   `json:"component_id"`
	Quantity                   float64 `json:"quantity"`
	UnitPrice                  float64 `json:"unit_price"`
}
