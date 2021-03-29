# Firestore utils for golang


## Data Handling

Add TimeStamped and Firestored to your objects, eg:

```go
type X struct {
    firetils.Firestored
    firetils.TimeStamped
}
```

You can do pre-saving and after-loading by adding `PreSave(ctx context.Context)` and/or `AfterLoad(ctx context.Context)` function to your models.

## Authentication

`Authenticate` function will validate an auth token.

Or use `FireAuth` middleware to do it automatically.

