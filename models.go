package firetils

import (
	"time"

	"cloud.google.com/go/firestore"
)

type KitchenSink interface {
	StoredAndStamped
	OwnedI
}

type StoredAndStamped interface {
	FirestoredI
	IDedI
	TimestampedI
}

// todo: should FirestoredI just be merged with IdedI ?
type StoredAndIded interface {
	FirestoredI
	IDedI
}

type FirestoredI interface {
	GetRef() *firestore.DocumentRef
	SetRef(*firestore.DocumentRef)
}
type IDedI interface {
	GetID() string
	SetID(string)
}

type OwnedI interface {
	GetUserID() string
	SetUserID(string)
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

// Timestamped snake case version
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

// TimeStamped camel case version
type TimeStamped struct {
	UpdatedAt time.Time `firestore:"updatedAt" json:"updatedAt"`
	CreatedAt time.Time `firestore:"createdAt" json:"createdAt"`
}

func (ts *TimeStamped) GetCreatedAt() time.Time {
	return ts.CreatedAt
}
func (ts *TimeStamped) GetUpdatedAt() time.Time {
	return ts.UpdatedAt
}
func (ts *TimeStamped) SetCreatedAt(t time.Time) {
	ts.CreatedAt = t
}
func (ts *TimeStamped) SetUpdatedAt(t time.Time) {
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

// Owned stores the user_id for the owner of this object
type Owned struct {
	UserID string `firestore:"user_id" json:"user_id"`
}

func (f *Owned) GetUserID() string {
	return f.UserID
}
func (f *Owned) SetUserID(id string) {
	f.UserID = id
}

// OwnedBy same as Owned, but camel case version
type OwnedBy struct {
	UserID string `firestore:"userID" json:"userID"`
}

func (f *OwnedBy) GetUserID() string {
	return f.UserID
}
func (f *OwnedBy) SetUserID(id string) {
	f.UserID = id
}
