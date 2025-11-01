package task

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"dev.helix.code/internal/database"
)

// TaskType represents different types of tasks
type TaskType string

const (
	TaskTypePlanning    TaskType = "planning"
	TaskTypeBuilding    TaskType = "building"
	TaskTypeTesting     TaskType = "testing"
	TaskTypeRefactoring TaskType = "refactoring"
	TaskTypeDebugging   TaskType = "debugging"
	TaskTypeDesign      TaskType = "design"
	TaskTypeDiagram     TaskType = "diagram"
	TaskTypeDeployment  TaskType = "deployment"
	TaskTypePorting     TaskType = "porting"
)

// TaskPriority represents task priority levels
type TaskPriority int

const (
	PriorityLow      TaskPriority = 1
	PriorityNormal   TaskPriority = 5
	PriorityHigh     TaskPriority = 10
	PriorityCritical TaskPriority = 20
)

// TaskCriticality represents task criticality levels
type TaskCriticality string

const (
	CriticalityLow      TaskCriticality = "low"
	CriticalityNormal   TaskCriticality = "normal"
	CriticalityHigh     TaskCriticality = "high"
	CriticalityCritical TaskCriticality = "critical"
)

// TaskStatus represents task status
type TaskStatus string

const (
	TaskStatusPending            TaskStatus = "pending"
	TaskStatusAssigned           TaskStatus = "assigned"
	TaskStatusRunning            TaskStatus = "running"
	TaskStatusCompleted          TaskStatus = "completed"
	TaskStatusFailed             TaskStatus = "failed"
	TaskStatusPaused             TaskStatus = "paused"
	TaskStatusWaitingForWorker   TaskStatus = "waiting_for_worker"
	TaskStatusWaitingForDeps     TaskStatus = "waiting_for_deps"
)

// ComplexityLevel represents task complexity
type ComplexityLevel string

const (
	ComplexityLow    ComplexityLevel = "low"
	ComplexityMedium ComplexityLevel = "medium"
	ComplexityHigh   ComplexityLevel = "high"
)

// Task represents a distributed task
type Task struct {
	ID              uuid.UUID       `json:"id"`
	Type            TaskType        `json:"type"`
	Data            map[string]interface{} `json:"data"`
	Status          TaskStatus      `json:"status"`
	Priority        TaskPriority    `json:"priority"`
	Criticality     TaskCriticality `json:"criticality"`
	AssignedWorker  *uuid.UUID      `json:"assigned_worker"`
	OriginalWorker  *uuid.UUID      `json:"original_worker"`
	Dependencies    []uuid.UUID     `json:"dependencies"`
	RetryCount      int             `json:"retry_count"`
	MaxRetries      int             `json:"max_retries"`
	ErrorMessage    string          `json:"error_message"`
	ResultData      map[string]interface{} `json:"result_data"`
	CheckpointData  map[string]interface{} `json:"checkpoint_data"`
	EstimatedDuration time.Duration `json:"estimated_duration"`
	StartedAt       *time.Time      `json:"started_at"`
	CompletedAt     *time.Time      `json:"completed_at"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// TaskManager manages distributed tasks
type TaskManager struct {
	db            *database.Database
	mu            sync.RWMutex
	tasks         map[uuid.UUID]*Task
	workers       map[uuid.UUID]*Worker
	queue         *TaskQueue
	checkpointMgr *CheckpointManager
	dependencyMgr *DependencyManager
}

// Worker represents a worker node
type Worker struct {
	ID                  uuid.UUID       `json:"id"`
	Hostname            string          `json:"hostname"`
	DisplayName         string          `json:"display_name"`
	SSHConfig           map[string]interface{} `json:"ssh_config"`
	Capabilities        []string        `json:"capabilities"`
	Resources           map[string]interface{} `json:"resources"`
	Status              string          `json:"status"`
	HealthStatus        string          `json:"health_status"`
	LastHeartbeat       *time.Time      `json:"last_heartbeat"`
	CPUUsagePercent     float64         `json:"cpu_usage_percent"`
	MemoryUsagePercent  float64         `json:"memory_usage_percent"`
	DiskUsagePercent    float64         `json:"disk_usage_percent"`
	CurrentTasksCount   int             `json:"current_tasks_count"`
	MaxConcurrentTasks  int             `json:"max_concurrent_tasks"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// TaskQueue manages task prioritization
type TaskQueue struct {
	highPriority   []*Task
	normalPriority []*Task
	lowPriority    []*Task
	mu             sync.RWMutex
}

// CheckpointManager manages task checkpoints
type CheckpointManager struct {
	db *database.Database
}

// DependencyManager manages task dependencies
type DependencyManager struct {
	db *database.Database
}

// TaskAnalysis represents analysis of a task for splitting
type TaskAnalysis struct {
	TaskID      uuid.UUID
	TaskType    TaskType
	Complexity  ComplexityLevel
	DataSize    int64
	Dependencies int
}

// TaskProgress represents task progress information
type TaskProgress struct {
	TaskID    uuid.UUID
	Status    TaskStatus
	Progress  float64
	StartedAt *time.Time
	UpdatedAt time.Time
}

// SplitStrategy defines interface for task splitting strategies
type SplitStrategy interface {
	GenerateSubtasks(parent *Task, analysis *TaskAnalysis) ([]SubtaskData, error)
}

// SubtaskData represents data for a subtask
type SubtaskData struct {
	Data         map[string]interface{}
	Dependencies []uuid.UUID
}

// NewTaskManager creates a new task manager
func NewTaskManager(db *database.Database) *TaskManager {
	return &TaskManager{
		db:            db,
		tasks:         make(map[uuid.UUID]*Task),
		workers:       make(map[uuid.UUID]*Worker),
		queue:         NewTaskQueue(),
		checkpointMgr: NewCheckpointManager(db),
		dependencyMgr: NewDependencyManager(db),
	}
}

// CreateTask creates a new task
func (tm *TaskManager) CreateTask(taskType TaskType, data map[string]interface{}, 
	priority TaskPriority, criticality TaskCriticality, dependencies []uuid.UUID) (*Task, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task := &Task{
		ID:              uuid.New(),
		Type:            taskType,
		Data:            data,
		Status:          TaskStatusPending,
		Priority:        priority,
		Criticality:     criticality,
		Dependencies:    dependencies,
		MaxRetries:      3,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validate dependencies
	if err := tm.dependencyMgr.ValidateDependencies(dependencies); err != nil {
		return nil, fmt.Errorf("invalid dependencies: %v", err)
	}

	// Store in memory
	tm.tasks[task.ID] = task

	// Add to database
	if err := tm.storeTaskInDB(task); err != nil {
		delete(tm.tasks, task.ID)
		return nil, fmt.Errorf("failed to store task in database: %v", err)
	}

	// Add to appropriate queue
	tm.queue.AddTask(task)

	log.Printf("✅ Task created: %s (type: %s, priority: %d)", task.ID, taskType, priority)
	return task, nil
}