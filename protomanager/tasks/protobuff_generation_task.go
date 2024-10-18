// protomanager/tasks/protobuf_generation_task.go
package tasks

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/Cdaprod/protomanager"
)

// ProtobufGenerationTask generates protobuf code for a specific package and language.
type ProtobufGenerationTask struct {
    ProtoManager *protomanager.ProtoManager
    PackageName  string
    Languages    []string
}

// Execute runs the protobuf generation task.
func (pt *ProtobufGenerationTask) Execute() (string, error) {
    pm := pt.ProtoManager
    pm.Logger.Infof("Generating protobufs for package '%s'", pt.PackageName)

    protoPath := filepath.Join(pm.MicroserviceProtoDir, pt.PackageName, "proto")
    outDir := filepath.Join(pm.OutputDir, pt.PackageName)

    for _, lang := range pt.Languages {
        outLangDir := filepath.Join(outDir, lang)
        if err := os.MkdirAll(outLangDir, os.ModePerm); err != nil {
            pm.Logger.Errorf("Failed to create output directory '%s': %v", outLangDir, err)
            return "", err
        }

        // Build protoc command
        protocArgs := []string{
            fmt.Sprintf("--%s_out=%s", lang, outLangDir),
            fmt.Sprintf("--%s-grpc_out=%s", lang, outLangDir),
            "--proto_path", protoPath,
            filepath.Join(protoPath, "*.proto"),
        }

        cmd := exec.Command("protoc", protocArgs...)
        output, err := cmd.CombinedOutput()
        if err != nil {
            pm.Logger.Errorf("Failed to generate protobufs for '%s' in '%s': %v\nOutput: %s", pt.PackageName, lang, err, string(output))
            return "", err
        }

        pm.Logger.Infof("Successfully generated protobufs for '%s' in '%s'", pt.PackageName, lang)
    }

    return fmt.Sprintf("Protobufs generated for package '%s'", pt.PackageName), nil
}