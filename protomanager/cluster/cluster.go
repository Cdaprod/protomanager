// protomanager/cluster/cluster.go
package cluster

import (
    "fmt"
    "sync"
)

// Task is the generic interface for tasks.
type Task[T any] interface {
    Execute() (T, error)
}

// ClusterManager manages and runs tasks concurrently.
type ClusterManager[T any] struct {
    tasks []Task[T]
    mu    sync.Mutex
}

// NewClusterManager initializes a new ClusterManager.
func NewClusterManager[T any]() *ClusterManager[T] {
    return &ClusterManager[T]{}
}

// AddTask adds a new task to the cluster.
func (cm *ClusterManager[T]) AddTask(task Task[T]) {
    cm.mu.Lock()
    defer cm.mu.Unlock()
    cm.tasks = append(cm.tasks, task)
}

// RunTasks executes all tasks concurrently and returns their results.
// If any task fails, it returns an aggregated error.
func (cm *ClusterManager[T]) RunTasks() ([]T, error) {
    var wg sync.WaitGroup
    results := make([]T, len(cm.tasks))
    errors := make([]error, len(cm.tasks))

    for i, task := range cm.tasks {
        wg.Add(1)
        go func(i int, task Task[T]) {
            defer wg.Done()
            result, err := task.Execute()
            results[i] = result
            errors[i] = err
        }(i, task)
    }

    wg.Wait()

    // Aggregate errors
    var aggErr string
    for _, err := range errors {
        if err != nil {
            aggErr += err.Error() + "; "
        }
    }

    if aggErr != "" {
        return results, fmt.Errorf("errors occurred: %s", aggErr)
    }

    return results, nil
}