package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/Cdaprod/protomanager/pkg/protomanager"
)

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

// DynamicRunner is a generic function that runs any given function concurrently with dynamic arguments.
func DynamicRunner[T any](fn func(T), arg T) {
	go func() {
		fn(arg)
	}()
}

// handleArgsOrSignals dynamically processes actions with generic argument types.
func handleArgsOrSignals[T any](pm *protomanager.ProtoManager, action string, arg T) {
	switch action {
	case "register":
		// Register microservice using dynamic arguments passed via CLI or signal
		DynamicRunner(pm.RegisterMicroservice, arg) // Run in a Go routine
	case "shutdown":
		fmt.Println("Shutting down ProtoManager...")
		// Add shutdown logic here if necessary
		DynamicRunner(func(_ struct{}) {
			time.Sleep(1 * time.Second) // Simulate some work
			fmt.Println("ProtoManager shutdown complete")
		}, struct{}{}) // Use empty struct as no argument is needed
	}
}

// ProtoManagerArguments holds dynamic arguments for a generic registration.
type ProtoManagerArguments struct {
	ServiceRepoPath string
	ServiceName     string
	Domain          string
}

func main() {
	// Parse CLI arguments using the flag package
	serviceRepoPath := flag.String("repo", "", "Path to the microservice repository containing the proto file")
	serviceName := flag.String("name", "", "Microservice name")
	domain := flag.String("domain", "", "Domain to register the microservice under")
	action := flag.String("action", "register", "Action to perform (register, shutdown, etc.)")
	globalProtoPath := flag.String("global-proto", "./proto/global.proto", "Path to the global.proto file")
	microserviceProtoDir := flag.String("proto-dir", "./proto/microservices", "Directory to store microservice proto files")
	outputDir := flag.String("output-dir", "./generated", "Directory for generated protobuf code")
	flag.Parse()

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Check if necessary arguments are provided
	if *serviceRepoPath == "" || *serviceName == "" || *domain == "" {
		fmt.Println("Error: missing required arguments --repo, --name, and --domain")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize ProtoManager
	pm, err := protomanager.NewProtoManager(*globalProtoPath, *microserviceProtoDir, *outputDir, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize ProtoManager: %v", err)
	}

	// Subscribe to events
	eventCh := pm.SubscribeEvents()
	go func() {
		for event := range eventCh {
			switch event.Type {
			case protomanager.EventTypeSuccess:
				logger.Infof("SUCCESS: %s", event.Message)
			case protomanager.EventTypeError:
				logger.Errorf("ERROR: %s", event.Message)
			case protomanager.EventTypeInfo:
				logger.Infof("INFO: %s", event.Message)
			}
		}
	}()

	// Initialize the signal handler
	signalHandler := NewSignalHandler()

	// Listen for system signals (like SIGINT, SIGTERM) in a separate goroutine
	go signalHandler.ListenForSignal()

	// Dynamic arguments for registering a service
	args := ProtoManagerArguments{
		ServiceRepoPath: *serviceRepoPath,
		ServiceName:     *serviceName,
		Domain:          *domain,
	}

	// Handle CLI args or signal-based actions dynamically using generics
	go func() {
		handleArgsOrSignals(pm, *action, args)

		// Process actions from signals dynamically
		for action := range signalHandler.ActionCh {
			handleArgsOrSignals(pm, action, args)
		}
	}()

	// Simulate waiting for incoming signals or shutdown (run indefinitely until interrupted)
	for {
		time.Sleep(1 * time.Second)
	}
}