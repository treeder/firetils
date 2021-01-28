# Firestore utils for golang


## Data Handling

Add TimeStamped and Firestored to your objects, eg:

```go
type X struct {
    firetils.Firestored
    firetils.Timestamped
}
```

TODO

## Authentication

`Authenticate` function will validate an auth token.

Or use `FireAuth` middleware to do it automatically.

