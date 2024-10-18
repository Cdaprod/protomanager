// protomanager/main.go
package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/Cdaprod/protomanager"
    "github.com/Cdaprod/protomanager/cmd"
    "github.com/Cdaprod/protomanager/cluster"
)

func main() {
    // Initialize logger
    logger := protomanager.NewLogger()

    // Parse flags
    globalProtoPath := flag.String("global-proto", "./proto/global.proto", "Path to the global.proto file")
    microserviceProtoDir := flag.String("proto-dir", "./proto/microservices", "Directory to store microservice proto files")
    outputDir := flag.String("output-dir", "./generated", "Directory for generated protobuf code")
    flag.Parse()

    // Initialize ProtoManager
    pm, err := protomanager.NewProtoManager(nil, *globalProtoPath, *microserviceProtoDir, *outputDir, logger)
    if err != nil {
        logger.Fatalf("Failed to initialize ProtoManager: %v", err)
    }

    // Initialize ClusterManager
    clusterManager := cluster.NewClusterManager[string]()

    // Subscribe to events
    pm.AddEventListener(func(event protomanager.Event) {
        switch event.Type {
        case "ServiceRegistered":
            logger.Infof("EVENT: %s", event.Message)
        case "Error":
            logger.Errorf("EVENT: %s", event.Message)
        default:
            logger.Infof("EVENT: %s", event.Message)
        }
    })

    // Initialize the signal handler
    signalHandler := protomanager.NewSignalHandler()

    // Listen for system signals in a separate goroutine
    go signalHandler.ListenForSignal()

    // Handle actions from signals
    go func() {
        for action := range signalHandler.ActionCh {
            switch action {
            case "shutdown":
                logger.Info("Shutting down ProtoManager...")
                // Perform any necessary cleanup here
                os.Exit(0)
            case "reload":
                logger.Info("Reloading configuration...")
                // Implement reload logic if applicable
            }
        }
    }()

    // Execute CLI commands
    go func() {
        cmd.Execute()
    }()

    // Keep the main goroutine alive
    select {}
}