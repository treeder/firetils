package firetils

import (
	"time"

	"cloud.google.com/go/firestore"
)

type Firestored struct {
	Ref *firestore.DocumentRef `firestore:"-" json:"-"`
}

type TimestampedI interface {
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}
type Timestamped struct {
	UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
}

func (ts *Timestamped) GetCreatedAt() time.Time {
	return ts.CreatedAt
}
func (ts *Timestamped) GetUpdatedAt() time.Time {
	return ts.UpdatedAt
}
func (ts *Timestamped) SetCreatedAt(t time.Time) {
	ts.CreatedAt = t
}
func (ts *Timestamped) SetUpdatedAt(t time.Time) {
	ts.UpdatedAt = t
}

// UpdateTimeStamps call this right before storing it in a database
func UpdateTimeStamps(obj TimestampedI) TimestampedI {
	if obj.GetCreatedAt().IsZero() {
		obj.SetCreatedAt(time.Now())
	}
	obj.SetUpdatedAt(time.Now())
	return obj
}
