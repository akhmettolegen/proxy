package managers

import "github.com/akhmettolegen/proxy/internal/models"

type TaskManager interface {
	TaskCreate(req *models.TaskRequest) (*models.TaskResponse, error)
	Create(req *models.TaskRequest, taskId string)
	ProcessRequest(req *models.TaskRequest, task *models.Task) error
	TaskById(id string) (*models.Task, error)
}
