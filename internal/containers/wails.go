package containers

// WailsFacade documents the service surface registered by main.go.
// Keeping this contract next to the adapter prevents transport concerns from
// leaking into the container use cases.
type WailsFacade = ContainerService
