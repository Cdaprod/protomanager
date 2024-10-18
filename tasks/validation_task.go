// protomanager/tasks/validation_task.go
package tasks

import (
    "fmt"
    "os/exec"
    "path/filepath"

    "github.com/Cdaprod/protomanager"
)

// ValidationTask validates protobuf files for a specific package.
type ValidationTask struct {
    ProtoManager *protomanager.ProtoManager
    PackageName  string
}

// Execute runs the validation task.
func (vt *ValidationTask) Execute() (string, error) {
    pm := vt.ProtoManager
    pkgPath := filepath.Join(pm.MicroserviceProtoDir, vt.PackageName, "proto")
    cmd := exec.Command("protoc", "--lint_out=.", filepath.Join(pkgPath, "*.proto"))

    output, err := cmd.CombinedOutput()
    if err != nil {
        pm.Logger.Errorf("Validation failed for package '%s': %v\nOutput: %s", vt.PackageName, err, string(output))
        return "", err
    }

    pm.Logger.Infof("Validation successful for package '%s'", vt.PackageName)
    return fmt.Sprintf("Validation successful for package '%s'", vt.PackageName), nil
}