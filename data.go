package firetils

import (
	"context"
	"reflect"

	"cloud.google.com/go/firestore"
	"github.com/treeder/gotils/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Save sets userID if it's in context, updates timestamps, calls PreSave(), then return object with Ref and ID set.
func Save(ctx context.Context, client *firestore.Client, collection string, v StoredAndStamped) (StoredAndStamped, error) {
	return Save2(ctx, client, collection, v, nil)
}

type SaveOptions struct {
	SkipOwned bool // won't set userID on the saved object
}

// Save2 same as Save, but with options
func Save2(ctx context.Context, client *firestore.Client, collection string, v StoredAndStamped, opts *SaveOptions) (StoredAndStamped, error) {
	if opts == nil {
		opts = &SaveOptions{}
	}
	UpdateTimeStamps(v)
	if !opts.SkipOwned {
		owned, ok := v.(OwnedI)
		if ok {
			if owned.GetUserID() == "" {
				// only set it if it's empty
				SetOwned(ctx, owned)
			}
		}
	}
	n := reflect.ValueOf(v)
	preSave := n.MethodByName("PreSave")
	if preSave.IsValid() {
		// fmt.Println("CALLING AfterLoad")
		_ = preSave.Call([]reflect.Value{reflect.ValueOf(ctx)})
		// TODO: check last value returned for an error!
	}
	var err error
	var ref *firestore.DocumentRef
	if v.GetRef() != nil {
		ref = v.GetRef()
		_, err = ref.Set(ctx, v)
	} else if v.GetID() != "" {
		// user set the ID
		ref = client.Collection(collection).Doc(v.GetID())
		_, err = ref.Set(ctx, v) // TODO should this be changed to create so we don't accidently overwrite something?
	} else {
		// make new ID
		ref, _, err = client.Collection(collection).Add(ctx, v)
	}
	if err != nil {
		return nil, gotils.C(ctx).Errorf("Failed to store object: %v", err)
	}
	v.SetRef(ref)
	v.SetID(ref.ID)
	return v, nil
}

// Merge does a Set with MergeAll
// Doesn't feel right...
// func Merge(ctx context.Context, client *firestore.Client, collection, id string, v map[string]interface{}) {
// 	// UpdateTimeStamps(v)
// 	v["updatedAt"] = time.Now()
// 	var err error
// 	var ref *firestore.DocumentRef
// 	ref = client.Collection(collection).Doc(id)
// 	_, err = ref.Set(ctx, v, firestore.MergeAll)
// 	if err != nil {
// 		return nil, gotils.C(ctx).Errorf("Failed to store object: %v", err)
// 	}
// }

// Collection returns CollectionRef by name
func Collection(client *firestore.Client, name string) *firestore.CollectionRef {
	return client.Collection(name)
}

// GetByID get a doc by ID
func GetByID(ctx context.Context, client *firestore.Client, collectionName, id string, t StoredAndStamped) error {
	ref := Collection(client, collectionName).Doc(id)
	return GetByRef(ctx, client, ref, t)
}

// GetByID2 same as GetByID, but returns new object
func GetByID2(ctx context.Context, client *firestore.Client, collectionName, id string, v StoredAndStamped) (StoredAndStamped, error) {
	t := reflect.TypeOf(v)
	// fmt.Printf("DATA: %+v\n", doc.Data())
	n := reflect.New(t.Elem())
	v2 := n.Interface()
	err := GetByID(ctx, client, collectionName, id, v2.(StoredAndStamped))
	if err != nil {
		return nil, err
	}
	return v2.(StoredAndStamped), nil
}

// NEED TO DO WHAT WE're DOING IN GET BY QUERY BELOW, and return the object, rather than filling in like DataTo
// GetByID get a doc by ID
// func GetByIDCached(ctx context.Context, collectionName, id string, t FirestoredI) error {
// 	ckey := fmt.Sprintf("%T-%v", t, id)
// 	fmt.Printf("XXX CKEY: %v\n", ckey)
// 	c2, _ := Cache.Get(ckey)
// 	if c2 != nil {
// 		user = c2.(*mmodels.User)
// 	}
// 	err := GetByID(ctx, collectionName, userID, user)
// 	if err != nil {
// 		return err
// 	}
// 	Cache.Set(ckey, user, 5)
// }

// GetByRef generic way to get a document
func GetByRef(ctx context.Context, client *firestore.Client, ref *firestore.DocumentRef, t StoredAndStamped) error {
	doc, err := ref.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return gotils.ErrNotFound
	}
	if err != nil {
		return gotils.C(ctx).Errorf("error in GetByRef: %v", err)
	}
	err = doc.DataTo(t)
	if err != nil {
		return gotils.C(ctx).Errorf("error in GetByRef: %v", err)
	}
	afterLoad(ctx, doc.Ref, t)
	return nil
}

// GetOneByQuery generic way to get a document
func GetOneByQuery(ctx context.Context, q firestore.Query, t StoredAndIded) error {
	iter := q.Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done {
		// gotils.WriteError(w, http.StatusInternalServerError, errors.New("Employee already added."))
		return gotils.ErrNotFound
	}
	if status.Code(err) == codes.NotFound {
		// gotils.WriteError(w, http.StatusBadRequest, gotils.ErrNotFound)
		return gotils.ErrNotFound
	}
	if err != nil {
		return gotils.C(ctx).Errorf("error in GetOneByQuery: %v", err)
	}
	err = doc.DataTo(t)
	if err != nil {
		return gotils.C(ctx).Errorf("error in GetOneByQuery: %v", err)
	}
	afterLoad(ctx, doc.Ref, t)
	return nil
}

// GetOneByQuery2 generic way to get a document
// This one returns a new object.
// Still trying to decide which way I like better...
func GetOneByQuery2(ctx context.Context, q firestore.Query, v StoredAndIded) (StoredAndIded, error) {
	t := reflect.TypeOf(v)
	// fmt.Printf("DATA: %+v\n", doc.Data())
	n := reflect.New(t.Elem())
	v2 := n.Interface()
	err := GetOneByQuery(ctx, q, v2.(StoredAndStamped))
	if err != nil {
		return nil, err
	}
	return v2.(StoredAndStamped), nil
}

// GetAllByQuery generic way to get a list of documents.
// NOTE: this doesn't seem to work well, best to use GetAllByQuery2
// limit restricts how many are returned. <=0 is all
// ret will be filled with the objects
func GetAllByQuery(ctx context.Context, q firestore.Query, limit int, ret []interface{}) error {
	// tType := t.Elem()
	// ret := []FirestoredI{}
	if limit > 0 {
		q = q.Limit(limit)
	}
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return gotils.C(ctx).Errorf("error iterating over query items: %v", err)
		}
		// HERE is how this all works: https://play.golang.org/p/tnsvwelTv4A
		t := reflect.TypeOf(ret)
		// fmt.Printf("DATA: %+v\n", doc.Data())
		n := reflect.New(t.Elem())
		v2 := n.Interface()
		err = doc.DataTo(v2)
		if err != nil {
			return gotils.C(ctx).Errorf("error on datato: %v", err)
		}
		fstored := v2.(StoredAndStamped)
		afterLoad(ctx, doc.Ref, fstored)
		ret = append(ret, fstored)
	}
	return nil
}

// GetAllByQuery2 generic way to get a list of documents, by just passing in the type.
// limit restricts how many we return. <=0 is all
// v is an instance of the type of object to be returned, it will not be modified or updated.
func GetAllByQuery2(ctx context.Context, q firestore.Query, v StoredAndStamped) ([]StoredAndStamped, error) {
	// tType := t.Elem()
	ret := []StoredAndStamped{}
	// if limit > 0 {
	// 	q = q.Limit(limit)
	// }
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, gotils.C(ctx).Errorf("error iterating over query items: %v", err)
		}
		// HERE is how this all works: https://play.golang.org/p/tnsvwelTv4A
		t := reflect.TypeOf(v)
		// fmt.Printf("DATA: %+v\n", doc.Data())
		n := reflect.New(t.Elem())
		v2 := n.Interface()
		err = doc.DataTo(v2)
		if err != nil {
			return nil, gotils.C(ctx).Errorf("error on DataTo: %v", err)
		}
		fstored := v2.(StoredAndStamped)
		afterLoad(ctx, doc.Ref, fstored)
		ret = append(ret, fstored)
	}
	return ret, nil
}
func afterLoad(ctx context.Context, ref *firestore.DocumentRef, v StoredAndIded) {
	v.SetRef(ref)
	v.SetID(ref.ID)
	// fmt.Printf("id: %v, status: %v\n", t.Ref.ID, t.Status)
	// could change this to another interface with PreSave and PostSave (or BeforeSave and AfterSave)
	n := reflect.ValueOf(v)
	afterLoad := n.MethodByName("AfterLoad")
	if afterLoad.IsValid() {
		// fmt.Println("CALLING AfterLoad")
		afterLoad.Call([]reflect.Value{reflect.ValueOf(ctx)})
		// TODO: check last value returned for an error!
	}
}

// GetByID get a doc by ID
func Delete(ctx context.Context, client *firestore.Client, collectionName, id string) error {
	ref := Collection(client, collectionName).Doc(id)
	_, err := ref.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
