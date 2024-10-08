package protomanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"github.com/sirupsen/logrus"
)

type ProtoManager struct {
	GlobalProtoPath       string
	MicroserviceProtoDir  string
	OutputDir             string
	Logger                *logrus.Logger
	mu                    sync.Mutex // To ensure safe concurrent access
}

func NewProtoManager(globalProtoPath, microserviceProtoDir, outputDir string, logger *logrus.Logger) (*ProtoManager, error) {
	return &ProtoManager{
		GlobalProtoPath:      globalProtoPath,
		MicroserviceProtoDir: microserviceProtoDir,
		OutputDir:            outputDir,
		Logger:               logger,
	}, nil
}

// RegisterMicroservice appends a new service definition to the global proto file and regenerates code.
func (pm *ProtoManager) RegisterMicroservice(serviceRepoPath, serviceName, domain string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Create service proto definition
	serviceProto := fmt.Sprintf(`
service %sService {
  rpc ExampleRPC (ExampleRequest) returns (ExampleResponse) {}
}

message ExampleRequest {
  string message = 1;
}

message ExampleResponse {
  string message = 1;
}
`, serviceName)

	// Read the existing global proto file
	globalProtoContent, err := ioutil.ReadFile(pm.GlobalProtoPath)
	if err != nil {
		pm.Logger.Errorf("Failed to read global proto file: %v", err)
		return err
	}

	// Append the new service to the global proto content
	globalProtoContent = append(globalProtoContent, []byte(serviceProto)...)

	// Write the updated proto file back to disk
	err = ioutil.WriteFile(pm.GlobalProtoPath, globalProtoContent, 0644)
	if err != nil {
		pm.Logger.Errorf("Failed to update global proto file: %v", err)
		return err
	}

	pm.Logger.Infof("Successfully registered microservice '%s' in global proto file.", serviceName)

	// Regenerate the proto code
	return pm.GenerateProtoCode()
}

// GenerateProtoCode regenerates Go code from the updated proto file.
func (pm *ProtoManager) GenerateProtoCode() error {
	pm.Logger.Info("Regenerating Go code from proto file...")

	// Command to regenerate the Go code from proto file
	cmd := exec.Command("protoc", "--go_out="+pm.OutputDir, "--go-grpc_out="+pm.OutputDir, pm.GlobalProtoPath)

	// Capture the output of the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.Logger.Errorf("Failed to regenerate Go code: %v", err)
		return err
	}

	pm.Logger.Infof("Proto code regenerated successfully. Output: %s", string(output))
	return nil
}