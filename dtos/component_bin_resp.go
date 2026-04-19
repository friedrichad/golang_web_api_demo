package dtos

type ComponentBinResponse struct {
	ComponentBinID int32   `json:"component_bin_id"`
	Quantity       float64 `json:"quantity"`
	ComponentID    int32   `json:"component_id"`
	BinID          int32   `json:"bin_id"`
}