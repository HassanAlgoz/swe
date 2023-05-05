package test

// Serivce-wide test: port -> controller -> store .. and back: store -> controller -> port -> test
// Test by calling port/grpc.go functions directly (since they are the entry point)
// but first, initialize stuff just like a `main.go` would.
