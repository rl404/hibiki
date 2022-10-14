package entity

// Status is user manga status.
type Status string

// Available user manga status.
const (
	StatusReading   Status = "READING"
	StatusCompleted Status = "COMPLETED"
	StatusOnHold    Status = "ON_HOLD"
	StatusDropped   Status = "DROPPED"
	StatusPlanned   Status = "PLANNED"
)

// Priority is user manga priority.
type Priority string

// Available user manga priority.
const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

// RereadValue is user manga reread value.
type RereadValue string

// Available user manga reread value.
const (
	RereadValueVeryLow  RereadValue = "VERY_LOW"
	RereadValueLow      RereadValue = "LOW"
	RereadValueMedium   RereadValue = "MEDIUM"
	RereadValueHigh     RereadValue = "HIGH"
	RereadValueVeryHigh RereadValue = "VERY_HIGH"
)
