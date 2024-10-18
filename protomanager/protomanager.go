// protomanager/protomanager.go
package protomanager

import (
    "fmt"
    "os"
    "os/exec"
    "os/signal"
    "sync"
    "syscall"

    "github.com/sirupsen/logrus"
)

// Event represents an event emitted by the ProtoManager.
type Event struct {
    Type    string
    Message string
}

// EventListener is a function that handles events.
type EventListener func(Event)

// ProtoManager handles protobuf file management and task execution.
type ProtoManager struct {
    ProtoRegistry         ProtoRegistry
    GlobalProtoPath       string
    MicroserviceProtoDir  string
    OutputDir             string
    Logger                *logrus.Logger
    eventListeners        []EventListener
    eventListenersMutex   sync.Mutex
    mu                    sync.Mutex // Ensures safe concurrent access
}

// NewProtoManager initializes a new ProtoManager.
func NewProtoManager(protoRegistry ProtoRegistry, globalProtoPath, microserviceProtoDir, outputDir string, logger *logrus.Logger) (*ProtoManager, error) {
    if protoRegistry == nil {
        protoRegistry = NewInternalProtoRegistry()
    }
    return &ProtoManager{
        ProtoRegistry:        protoRegistry,
        GlobalProtoPath:      globalProtoPath,
        MicroserviceProtoDir: microserviceProtoDir,
        OutputDir:            outputDir,
        Logger:               logger,
        eventListeners:       []EventListener{},
    }, nil
}

// NewLogger initializes and returns a new logger.
func NewLogger() *logrus.Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })
    logger.SetLevel(logrus.InfoLevel)
    return logger
}

// AddEventListener allows external entities to subscribe to ProtoManager events.
func (pm *ProtoManager) AddEventListener(listener EventListener) {
    pm.eventListenersMutex.Lock()
    defer pm.eventListenersMutex.Unlock()
    pm.eventListeners = append(pm.eventListeners, listener)
}

// emitEvent sends an event to all listeners.
func (pm *ProtoManager) emitEvent(event Event) {
    pm.eventListenersMutex.Lock()
    listeners := make([]EventListener, len(pm.eventListeners))
    copy(listeners, pm.eventListeners)
    pm.eventListenersMutex.Unlock()

    for _, listener := range listeners {
        go listener(event)
    }
}

// RegisterMicroservice registers a new microservice.
func (pm *ProtoManager) RegisterMicroservice(serviceName, domain, version string) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    metadata := ServiceMetadata{
        Domain:  domain,
        Version: version,
    }

    // Register service in ProtoRegistry
    if err := pm.ProtoRegistry.RegisterService(serviceName, metadata); err != nil {
        pm.Logger.Errorf("Failed to register service '%s': %v", serviceName, err)
        pm.emitEvent(Event{Type: "Error", Message: fmt.Sprintf("Failed to register service '%s': %v", serviceName, err)})
        return err
    }

    pm.Logger.Infof("Successfully registered service '%s'", serviceName)
    pm.emitEvent(Event{Type: "ServiceRegistered", Message: fmt.Sprintf("Service '%s' registered", serviceName)})

    // Update global proto file
    if err := pm.updateGlobalProto(serviceName, metadata); err != nil {
        return err
    }

    // Generate code
    if err := pm.GenerateProtoCode(); err != nil {
        return err
    }

    return nil
}

// updateGlobalProto updates the global proto file with the new service.
func (pm *ProtoManager) updateGlobalProto(serviceName string, metadata ServiceMetadata) error {
    pm.Logger.Infof("Updating global proto file for service '%s'", serviceName)

    serviceDefinition := fmt.Sprintf(`
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

    // Append service definition to global proto file
    file, err := os.OpenFile(pm.GlobalProtoPath, os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        pm.Logger.Errorf("Failed to open global proto file: %v", err)
        pm.emitEvent(Event{Type: "Error", Message: fmt.Sprintf("Failed to open global proto file: %v", err)})
        return err
    }
    defer file.Close()

    if _, err := file.WriteString(serviceDefinition); err != nil {
        pm.Logger.Errorf("Failed to write to global proto file: %v", err)
        pm.emitEvent(Event{Type: "Error", Message: fmt.Sprintf("Failed to write to global proto file: %v", err)})
        return err
    }

    pm.Logger.Infof("Global proto file updated for service '%s'", serviceName)
    pm.emitEvent(Event{Type: "GlobalProtoUpdated", Message: fmt.Sprintf("Global proto updated for service '%s'", serviceName)})
    return nil
}

// GenerateProtoCode regenerates code from the updated proto files.
func (pm *ProtoManager) GenerateProtoCode() error {
    pm.Logger.Info("Regenerating code from proto files...")

    cmd := exec.Command("protoc",
        "--go_out="+pm.OutputDir,
        "--go-grpc_out="+pm.OutputDir,
        "--proto_path="+pm.MicroserviceProtoDir,
        pm.GlobalProtoPath,
    )

    output, err := cmd.CombinedOutput()
    if err != nil {
        pm.Logger.Errorf("Failed to regenerate code: %v\nOutput: %s", err, string(output))
        pm.emitEvent(Event{Type: "Error", Message: fmt.Sprintf("Failed to regenerate code: %v", err)})
        return err
    }

    pm.Logger.Infof("Code regenerated successfully. Output: %s", string(output))
    pm.emitEvent(Event{Type: "CodeGenerated", Message: "Code regenerated successfully"})
    return nil
}

// SignalHandler defines the structure for receiving signals to trigger actions.
type SignalHandler struct {
    SignalCh chan os.Signal
    ActionCh chan string
}

// NewSignalHandler initializes and returns a SignalHandler.
func NewSignalHandler() *SignalHandler {
    return &SignalHandler{
        SignalCh: make(chan os.Signal, 1),
        ActionCh: make(chan string, 1),
    }
}

// ListenForSignal waits for system signals and passes corresponding actions to the action channel.
func (sh *SignalHandler) ListenForSignal() {
    signal.Notify(sh.SignalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
    for {
        select {
        case sig := <-sh.SignalCh:
            switch sig {
            case syscall.SIGINT, syscall.SIGTERM:
                fmt.Println("\nReceived shutdown signal")
                sh.ActionCh <- "shutdown"
                return
            case syscall.SIGHUP:
                fmt.Println("\nReceived reload signal")
                sh.ActionCh <- "reload"
            }
        }
    }
}