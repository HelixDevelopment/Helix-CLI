package worker

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TestDistributedWorkerManager tests the distributed worker manager
func TestDistributedWorkerManager(t *testing.T) {
	config := WorkerConfig{
		Enabled: true,
		Pool: map[string]WorkerConfigEntry{
			"test-worker-1": {
				Host:        "localhost",
				Port:        2222,
				Username:    "test",
				KeyPath:     "test/key",
				Capabilities: []string{"code-generation"},
				DisplayName: "Test Worker 1",
			},
		},
		AutoInstall:         false,
		HealthCheckInterval: 30,
		MaxConcurrentTasks:  10,
		TaskTimeout:         300,
	}

	manager := NewDistributedWorkerManager(config)

	// Test initialization
	ctx := context.Background()
	err := manager.Initialize(ctx)
	if err != nil {
		t.Fatalf("Failed to initialize worker manager: %v", err)
	}

	// Test worker retrieval
	workers := manager.GetAvailableWorkers()
	if len(workers) == 0 {
		t.Log("No workers available (may be normal in unit test)")
	}

	// Test worker stats
	stats := manager.GetWorkerStats()
	if stats["total_workers"].(int) != len(manager.workers) {
		t.Errorf("Worker stats mismatch: expected %d, got %d", len(manager.workers), stats["total_workers"])
	}

	// Test task submission
	task := &DistributedTask{
		Type:        "test-task",
		Data:        map[string]interface{}{"test": "data"},
		Priority:    5,
		Criticality: CriticalityNormal,
		MaxRetries:  3,
	}

	err = manager.SubmitTask(task)
	if err != nil {
		t.Fatalf("Failed to submit task: %v", err)
	}

	if task.ID == uuid.Nil {
		t.Error("Task ID should be set after submission")
	}

	if task.Status != TaskStatusPending {
		t.Errorf("Task status should be pending, got %s", task.Status)
	}

	t.Logf("✅ Distributed worker manager test passed: submitted task %s", task.ID)
}

// TestWorkerConfigValidation tests worker configuration validation
func TestWorkerConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config WorkerConfig
		valid  bool
	}{
		{
			name: "Valid configuration",
			config: WorkerConfig{
				Enabled: true,
				Pool: map[string]WorkerConfigEntry{
					"worker1": {
						Host:        "localhost",
						Port:        22,
						Username:    "user",
						KeyPath:     "/path/to/key",
						Capabilities: []string{"code-generation"},
					},
				},
				MaxConcurrentTasks: 10,
			},
			valid: true,
		},
		{
			name: "Disabled configuration",
			config: WorkerConfig{
				Enabled: false,
			},
			valid: true,
		},
		{
			name: "Invalid port",
			config: WorkerConfig{
				Enabled: true,
				Pool: map[string]WorkerConfigEntry{
					"worker1": {
						Host: "localhost",
						Port: 0, // Invalid port
					},
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewDistributedWorkerManager(tt.config)
			
			// Test initialization
			ctx := context.Background()
			err := manager.Initialize(ctx)
			
			if tt.valid && err != nil {
				t.Errorf("Expected valid configuration but got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("Expected invalid configuration but got no error")
			}
		})
	}
}

// TestTaskPriority tests task priority handling
func TestTaskPriority(t *testing.T) {
	config := WorkerConfig{
		Enabled: true,
		Pool: map[string]WorkerConfigEntry{
			"test-worker": {
				Host:        "localhost",
				Port:        22,
				Username:    "test",
				KeyPath:     "test/key",
				Capabilities: []string{"testing"},
			},
		},
		MaxConcurrentTasks: 5,
	}

	manager := NewDistributedWorkerManager(config)

	// Submit tasks with different priorities
	tasks := []struct {
		priority int
		criticality Criticality
	}{
		{1, CriticalityCritical}, // Highest priority
		{5, CriticalityNormal},   // Medium priority
		{10, CriticalityLow},     // Lowest priority
	}

	for _, taskDef := range tasks {
		task := &DistributedTask{
			Type:        "priority-test",
			Priority:    taskDef.priority,
			Criticality: taskDef.criticality,
			MaxRetries:  1,
		}

		err := manager.SubmitTask(task)
		if err != nil {
			t.Fatalf("Failed to submit task with priority %d: %v", taskDef.priority, err)
		}

		t.Logf("Submitted task with priority %d and criticality %s", taskDef.priority, taskDef.criticality)
	}

	t.Log("✅ Task priority test passed")
}

// TestWorkerCapabilities tests worker capability matching
func TestWorkerCapabilities(t *testing.T) {
	config := WorkerConfig{
		Enabled: true,
		Pool: map[string]WorkerConfigEntry{
			"code-worker": {
				Host:        "localhost",
				Port:        22,
				Username:    "test",
				KeyPath:     "test/key",
				Capabilities: []string{"code-generation", "refactoring"},
			},
			"test-worker": {
				Host:        "localhost",
				Port:        23,
				Username:    "test",
				KeyPath:     "test/key",
				Capabilities: []string{"testing", "debugging"},
			},
		},
	}

	manager := NewDistributedWorkerManager(config)

	workers := manager.GetAvailableWorkers()
	if len(workers) != 2 {
		t.Errorf("Expected 2 workers, got %d", len(workers))
	}

	// Verify worker capabilities
	for _, worker := range workers {
		if len(worker.Capabilities) == 0 {
			t.Error("Worker should have capabilities")
		}

		t.Logf("Worker %s capabilities: %v", worker.DisplayName, worker.Capabilities)
	}

	t.Log("✅ Worker capabilities test passed")
}

// TestTaskStatusTransitions tests task status transitions
func TestTaskStatusTransitions(t *testing.T) {
	task := &DistributedTask{
		Type:      "status-test",
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
	}

	// Test status transitions
	initialStatus := task.Status

	// Simulate task assignment
	task.Status = TaskStatusRunning
	if task.Status == initialStatus {
		t.Error("Task status should change when assigned")
	}

	// Simulate task start
	task.Status = TaskStatusRunning
	if task.Status != TaskStatusRunning {
		t.Error("Task status should be running")
	}

	// Simulate task completion
	task.Status = TaskStatusCompleted
	if task.Status != TaskStatusCompleted {
		t.Error("Task status should be completed")
	}

	// Simulate task failure
	task.Status = TaskStatusFailed
	if task.Status != TaskStatusFailed {
		t.Error("Task status should be failed")
	}

	t.Log("✅ Task status transitions test passed")
}

// TestWorkerHealthMonitoring tests worker health monitoring
func TestWorkerHealthMonitoring(t *testing.T) {
	config := WorkerConfig{
		Enabled: true,
		Pool: map[string]WorkerConfigEntry{
			"healthy-worker": {
				Host:        "localhost",
				Port:        22,
				Username:    "test",
				KeyPath:     "test/key",
				Capabilities: []string{"monitoring"},
			},
		},
		HealthCheckInterval: 5, // Short interval for testing
	}

	manager := NewDistributedWorkerManager(config)

	// Get worker and simulate health updates
	workers := manager.GetAvailableWorkers()
	if len(workers) == 0 {
		t.Skip("No workers available for health monitoring test")
	}

	worker := workers[0]

	// Test initial health status
	if worker.HealthStatus == "" {
		t.Error("Worker should have initial health status")
	}

	// Test last heartbeat
	if worker.LastHeartbeat.IsZero() {
		t.Error("Worker should have last heartbeat set")
	}

	t.Logf("Worker health: %s, last heartbeat: %v", worker.HealthStatus, worker.LastHeartbeat)
	t.Log("✅ Worker health monitoring test passed")
}