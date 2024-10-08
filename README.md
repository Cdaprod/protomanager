# ProtoManager

ProtoManager is a robust and scalable tool designed to manage Protocol Buffer (protobuf) workflows within microservices architectures. It automates the registration of microservices, manages global proto files, and ensures consistent generation of protobuf-based code across services. Leveraging Go's concurrency features and Viper for configuration management, ProtoManager provides a flexible and efficient solution for large-scale microservices management.

## Features

- **Dynamic Microservice Registration**: Easily register new microservices and manage their protobuf definitions.
- **Global Proto File Management**: Maintain a central proto file that aggregates all service definitions for consistency.
- **Automated Code Generation**: Automatically generate language-specific protobuf code whenever the global proto file is updated.
- **Configuration Management**: Use Viper for flexible and dynamic configuration handling.
- **Signal Handling**: Gracefully handle system signals for shutdowns and configuration reloads.
- **Event Subscription**: Subscribe to system events for real-time monitoring and logging.

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

####Registering a Microservice

To register a new microservice and update the global proto file:

go run cmd/main.go --repo /path/to/repo --name myservice --domain mydomain --action register --config ./config.yaml

	•	--repo: Path to the microservice repository.
	•	--name: Name of the microservice.
	•	--domain: Domain under which the microservice is registered.
	•	--action: Action to perform (register, shutdown, etc.).
	•	--config: Path to the configuration file.

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

protomanager/
├── cmd/
│   └── main.go                 // CLI entry point
├── pkg/
│   └── protomanager/
│       ├── manager.go          // Core logic for managing proto files and microservices
│       ├── event.go            // Event definitions and subscription logic
│       ├── state.go            // State management (if any)
│       └── manager_test.go     // Tests for ProtoManager
├── proto/
│   ├── global.proto            // Global proto file aggregating all service definitions
│   └── microservices/          // Directory for individual microservice proto files
├── .github/
│   └── workflows/
│       └── ci.yml              // GitHub Actions CI workflow
├── config.yaml                 // Configuration file
├── go.mod                      // Go module dependencies
└── go.sum                      // Dependency checksum file

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
