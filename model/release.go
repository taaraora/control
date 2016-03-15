package model

type Release struct {
	ID        int        `json:"id"`   // Should be auto-generated by incrementing
	Type      string     `json:"type"` // initial, bluegreen, canary, rolling, hard, etc.
	Component *Component `json:"component"`
}
