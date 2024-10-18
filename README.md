# ProtoManager

ProtoManager is a robust and scalable tool designed to manage Protocol Buffer (protobuf) workflows within microservices architectures. It automates the registration of microservices, manages global proto files, and ensures consistent generation of protobuf-based code across services. Leveraging Go's concurrency features and Viper for configuration management, ProtoManager provides a flexible and efficient solution for large-scale microservices management.

## Features

- **Dynamic Microservice Registration**: Easily register new microservices and manage their protobuf definitions.
- **Global Proto File Management**: Maintain a central proto file that aggregates all service definitions for consistency.
- **Automated Code Generation**: Automatically generate language-specific protobuf code whenever the global proto file is updated.
- **Configuration Management**: Use Viper for flexible and dynamic configuration handling.
- **Signal Handling**: Gracefully handle system signals for shutdowns and configuration reloads.
- **Event Subscription**: Subscribe to system events for real-time monitoring and logging.

#### Notes

	•	No Direct Dependency: The protomanager package does not import or depend on any external Registry. It relies on the ProtoRegistry interface, ensuring loose coupling.
	•	Event-Driven Architecture: The ProtoManager emits events that can be listened to by other components, facilitating communication without direct dependencies.
	•	Modularity and Flexibility: The use of interfaces and dependency injection allows the protomanager to function independently or integrate with external systems as needed.
	•	Generics and Concurrency: The ClusterManager and tasks utilize Go’s generics and concurrency features for efficient and type-safe task execution.

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git
- Protobuf Compiler (`protoc`)
- GitHub Actions (for CI workflows)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/Cdaprod/protomanager.git
   cd protomanager

2.	Install Dependencies

go mod tidy


3.	Setup Configuration

Create a config.yaml file in the root directory with the following content:

global_proto_path: "./proto/global.proto"
microservice_proto_dir: "./proto/microservices"
output_dir: "./generated"

Ensure that the proto directory contains a global.proto file that acts as the central repository of all microservice proto definitions.

### Usage

#### Building the CLI

cd Cdaprod/protomanager
go build -o protomanager .

#### Registering a Microservice

To register a new microservice and update the global proto file:

go run cmd/main.go --repo /path/to/repo --name myservice --domain mydomain --action register --config ./config.yaml

	•	--repo: Path to the microservice repository.
	•	--name: Name of the microservice.
	•	--domain: Domain under which the microservice is registered.
	•	--action: Action to perform (register, shutdown, etc.).
	•	--config: Path to the configuration file.
	
#### Generating Protobuffs

./protomanager generate --packages git,docker --languages go,python --push --validate

#### Using External Registry

```go
// main.go

func main() {
    // Initialize logger
    logger := protomanager.NewLogger()

    // Initialize external registry
    extRegistry := registry.NewExternalRegistry()

    // Initialize ProtoManager with external ProtoRegistry
    pm, err := protomanager.NewProtoManager(extRegistry, "./proto/global.proto", "./proto/microservices", "./generated", logger)
    if err != nil {
        logger.Fatalf("Failed to initialize ProtoManager: %v", err)
    }

    // Rest of the code...
}
``` 

#### Handling Signals

The application listens for system signals such as SIGINT, SIGTERM, and SIGHUP to perform actions like shutdown and configuration reloads.

	•	Shutdown: Send SIGINT or SIGTERM to gracefully shut down the application.
	•	Reload Configuration: Send SIGHUP to reload the configuration file without restarting the application.

#### Testing

ProtoManager includes comprehensive tests to ensure functionality.

go test ./pkg/protomanager/...

#### CI/CD

ProtoManager uses GitHub Actions for continuous integration. The CI workflow is defined in .github/workflows/ci.yml.

#### Project Structure

CdaprodParentApp/
├── protomanager/
│   ├── cmd/
│   │   └── command.go
│   ├── cluster/
│   │   └── cluster.go
│   ├── tasks/
│   │   ├── protobuf_generation_task.go
│   │   ├── push_task.go
│   │   └── validation_task.go
│   ├── proto_registry.go
│   ├── internal_proto_registry.go
│   ├── protomanager.go
│   └── main.go
├── external/
│   └── registry/
│       └── registry.go
└── proto/
    └── global.proto

### Contributing

Contributions are welcome! Please follow these steps:

	1.	Fork the repository.
	2.	Create your feature branch (git checkout -b feature/NewFeature).
	3.	Commit your changes (git commit -m 'Add some feature').
	4.	Push to the branch (git push origin feature/NewFeature).
	5.	Open a Pull Request.

### License

This project is licensed under the MIT License - see the LICENSE file for details.

### Contact

For any inquiries or support, please contact Cdaprod.
