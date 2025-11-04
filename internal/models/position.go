package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Position represents a position in an act with cost information
type Position struct {
	ID                              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CurrentPeriodCost               *float64           `json:"currentPeriodCost,omitempty" bson:"currentPeriodCost,omitempty"`
	CurrentPeriodCostInspection     *float64           `json:"currentPeriodCostInspection,omitempty" bson:"currentPeriodCostInspection,omitempty"`
	CurrentPeriodCostConsiderations *float64           `json:"currentPeriodCostConsiderations,omitempty" bson:"currentPeriodCostConsiderations,omitempty"`
	AccumulatedCost                 *float64           `json:"accumulatedCost,omitempty" bson:"accumulatedCost,omitempty"`
}

// HasCurrentPeriodCost checks if position has any current period cost
func (p *Position) HasCurrentPeriodCost() bool {
	return p.CurrentPeriodCost != nil ||
		p.CurrentPeriodCostInspection != nil ||
		p.CurrentPeriodCostConsiderations != nil
}

// HasAccumulatedCost checks if position has accumulated cost
func (p *Position) HasAccumulatedCost() bool {
	return p.AccumulatedCost != nil
}
