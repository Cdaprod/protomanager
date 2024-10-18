// protomanager/tasks/push_task.go
package tasks

import (
    "fmt"
    "os/exec"
    "path/filepath"

    "github.com/Cdaprod/protomanager"
)

// PushTask pushes generated protobufs to the repository.
type PushTask struct {
    ProtoManager *protomanager.ProtoManager
    PackageName  string
    CommitMsg    string
}

// Execute runs the push task.
func (pt *PushTask) Execute() (string, error) {
    pm := pt.ProtoManager
    repoPath := filepath.Join(pm.MicroserviceProtoDir, pt.PackageName)

    pm.Logger.Infof("Pushing protobufs for package '%s' to repository at '%s'", pt.PackageName, repoPath)

    // Add changes
    cmdAdd := exec.Command("git", "add", ".")
    cmdAdd.Dir = repoPath
    if output, err := cmdAdd.CombinedOutput(); err != nil {
        pm.Logger.Errorf("Failed to add changes in '%s': %v\nOutput: %s", repoPath, err, string(output))
        return "", err
    }

    // Commit changes
    cmdCommit := exec.Command("git", "commit", "-m", pt.CommitMsg)
    cmdCommit.Dir = repoPath
    if output, err := cmdCommit.CombinedOutput(); err != nil {
        pm.Logger.Errorf("Failed to commit changes in '%s': %v\nOutput: %s", repoPath, err, string(output))
        return "", err
    }

    // Push changes
    cmdPush := exec.Command("git", "push")
    cmdPush.Dir = repoPath
    if output, err := cmdPush.CombinedOutput(); err != nil {
        pm.Logger.Errorf("Failed to push changes in '%s': %v\nOutput: %s", repoPath, err, string(output))
        return "", err
    }

    pm.Logger.Infof("Successfully pushed protobufs for package '%s'", pt.PackageName)
    return fmt.Sprintf("Pushed protobufs for package '%s'", pt.PackageName), nil
}