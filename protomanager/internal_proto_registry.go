// protomanager/proto_registry.go
package protomanager

// ProtoRegistry defines the interface for registering and retrieving services.
type ProtoRegistry interface {
    RegisterService(serviceName string, metadata ServiceMetadata) error
    GetService(serviceName string) (ServiceMetadata, error)
}

// ServiceMetadata holds metadata about a service.
type ServiceMetadata struct {
    Domain   string
    Version  string
    // Additional fields as needed
}