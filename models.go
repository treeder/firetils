package firetils

import (
	"time"

	"cloud.google.com/go/firestore"
)

type StoredAndStamped interface {
	FirestoredI
	IDedI
	TimestampedI
}

type FirestoredI interface {
	GetRef() *firestore.DocumentRef
	SetRef(*firestore.DocumentRef)
}
type IDedI interface {
	GetID() string
	SetID(string)
}

type Firestored struct {
	Ref *firestore.DocumentRef `firestore:"-" json:"-"`
}

func (f *Firestored) GetRef() *firestore.DocumentRef {
	return f.Ref
}
func (f *Firestored) SetRef(ref *firestore.DocumentRef) {
	f.Ref = ref
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

type IDed struct {
	ID string `firestore:"-" json:"id"`
}

func (f *IDed) GetID() string {
	return f.ID
}
func (f *IDed) SetID(id string) {
	f.ID = id
}
