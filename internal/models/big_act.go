package models

// BigAct represents aggregated information for a big act
type BigAct struct {
	Changed                 bool                   `json:"changed" bson:"changed"`
	TotalCost               float64                `json:"totalCost,omitempty" bson:"totalCost,omitempty"`
	TotalCostInspection     float64                `json:"totalCostInspection,omitempty" bson:"totalCostInspection,omitempty"`
	TotalCostConsiderations float64                `json:"totalCostConsiderations,omitempty" bson:"totalCostConsiderations,omitempty"`
	PositionIDs             string                 `json:"positionIds,omitempty" bson:"positionIds,omitempty"`
	BigActLink              string                 `json:"bigActLink,omitempty" bson:"bigActLink,omitempty"`
	TextFields              map[string]interface{} `json:"textFields,omitempty" bson:"textFields,omitempty"`
}
