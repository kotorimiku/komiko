package service

import (
	"context"
	"komiko/model"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskFunc func(ctx context.Context, task *model.Task)

type TaskManager struct {
	mu     sync.Mutex
	tasks  map[string]*model.Task
	cancel map[string]context.CancelFunc
}

var taskManager *TaskManager
var once sync.Once

func GetTaskManager() *TaskManager {
	once.Do(func() {
		taskManager = &TaskManager{
			tasks:  make(map[string]*model.Task),
			cancel: make(map[string]context.CancelFunc),
		}
	})
	return taskManager
}

func (tm *TaskManager) AddTask(name string, fn TaskFunc) string {
	id := uuid.New().String()
	task := &model.Task{
		ID:        id,
		Name:      name,
		Status:    model.TaskPending,
		Progress:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	tm.mu.Lock()
	for _, t := range tm.tasks {
		if t.Name == name && (t.Status == model.TaskPending || t.Status == model.TaskRunning) {
			return t.ID
		}
	}
	tm.tasks[id] = task
	tm.mu.Unlock()
	go tm.runTask(task, fn)
	return id
}

func (tm *TaskManager) runTask(task *model.Task, fn TaskFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	tm.mu.Lock()
	tm.cancel[task.ID] = cancel
	tm.mu.Unlock()
	task.Status = model.TaskRunning
	task.UpdatedAt = time.Now()

	defer func() {
		tm.mu.Lock()
		delete(tm.cancel, task.ID)
		tm.mu.Unlock()
	}()

	fn(ctx, task)
	task.Status = model.TaskCompleted
	task.Result = "任务完成"
	task.UpdatedAt = time.Now()
}

func (tm *TaskManager) StopTask(id string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if cancel, ok := tm.cancel[id]; ok {
		cancel()
		if task, ok := tm.tasks[id]; ok {
			task.Status = model.TaskFailed
			task.Error = "stopped by user"
			task.UpdatedAt = time.Now()
		}
	}
}

func (tm *TaskManager) GetTask(id string) *model.Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if task, ok := tm.tasks[id]; ok {
		return task
	}
	return nil
}

func (tm *TaskManager) ListTasks() []*model.Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tasks := make([]*model.Task, 0, len(tm.tasks))
	for _, t := range tm.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}
