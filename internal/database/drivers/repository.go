package drivers

import (
	"context"
	"github.com/akhmettolegen/proxy/internal/models"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	Update(ctx context.Context, task *models.Task) error
	TaskById(ctx context.Context, id string) (*models.Task, error)
}
