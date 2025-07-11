package model

// RecordID defines a record id. Together with RecordType
// identifies unique records accross all types
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records accross all types
type RecordType string

// Existing Record Type
const (
	Movie = RecordType("movie")
)

// UserID defines a user id
type UserID string

// RatingValue defines a value of a rating record
type RatingValue int

type Rating struct {
	RecordID   string      `json:"recordId"`
	RecordType string      `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

// RatingEvent defines an event containing rating information
type RatingEvent struct {
	Rating
	ProviderID string          `json:"providerId"`
	EventType  RatingEventType `json:"eventType"`
}

// RatingEventType defines the type of a rating event
type RatingEventType string

// Rating event types
const (
	RatingEventTypePut    = RatingEventType("put")
	RatingEventTypeDelete = RatingEventType("delete")
)
