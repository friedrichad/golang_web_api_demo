package dtos

import("fmt")

type ComponentRequest struct {
	MetadataJSON string  `json:"metadata_json"`
	Unit         string  `json:"unit"`
	UnitPrice    float64 `json:"unit_price"`
	AddComponentToBinRequest `json:"component_bin"`
	ComponentCategoryID int32 `json:"component_category_id"`
}
type AddComponentToBinRequest struct{
	BinID int32 `json:"bin_id"`
	Quantity float64 `json:"quantity"`
}

func (r *ComponentRequest) Verify() (bool, error) {
	if r.MetadataJSON == "" {
		return false, fmt.Errorf("metadata_json is required")
	}
	if r.Unit == "" {
		return false, fmt.Errorf("unit is required")
	}
	if r.UnitPrice <= 0 {
		return false, fmt.Errorf("unit_price must be a positive number")
	}
	if r.AddComponentToBinRequest.BinID == 0 {
		return false, fmt.Errorf("bin_id is required")
	}
	if r.AddComponentToBinRequest.Quantity <= 0 {
		return false, fmt.Errorf("quantity must be a positive number")
	}
	if r.ComponentCategoryID <= 0 {
		return false, fmt.Errorf("component_category_id is required")
	}
	return true, nil
}
