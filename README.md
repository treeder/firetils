# Firestore utils for golang

Tools that make it easier to work with Google Cloud / Firebase's Firestore.

[![Go Reference](https://pkg.go.dev/badge/github.com/treeder/firetils.svg)](https://pkg.go.dev/github.com/treeder/firetils)

For instance:

```go
ob := &MyObject{}
err := firetils.Save(ctx, client, "myCollection", ob) // depending on interfaces on the object, it will update timestamps, call PreSave() function, etc
```

Fetching:

There's several functions here, they will all call populate into a struct and call AfterLoad(). 

```go
firetils.GetByID
firetils.GetOneByQuery
firetils.GetAllByQuery
```

## Data Handling

Add TimeStamped and Firestored to your objects, eg:

```go
type X struct {
    firetils.Firestored
    firetils.TimeStamped
    firetils.IDed
}
```

You can do pre-saving and after-loading by adding `PreSave(ctx context.Context)` and/or `AfterLoad(ctx context.Context)` function to your models.

## Authentication

`Authenticate` function will validate an auth token.

Or use `firetils.FireAuth` middleware to do it automatically.
