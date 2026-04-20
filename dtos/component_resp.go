package dtos

type ComponentResponse struct {
	ComponentID  int32   `json:"component_id"`
	MetadataJSON string  `json:"metadata_json"`
	Unit         string  `json:"unit"`
	UnitPrice    float64 `json:"unit_price"`
	AddComponentToBinRequest `json:"component_bin"`
	ComponentCategoryResponse `json:"component_category"`
}
