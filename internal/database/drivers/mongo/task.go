package mongo

import (
	"context"
	"github.com/akhmettolegen/proxy/internal/database/drivers"
	"github.com/akhmettolegen/proxy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	if task == nil {
		return drivers.ErrTaskEmpty
	}

	task.CreatedAt = time.Now().In(time.UTC)

	if _, err := r.collection.InsertOne(ctx, task); err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {

	filter := bson.D{
		{Key: "id", Value: task.Id},
	}

	updateFields := bson.D{
		{"status", task.Status},
		{"httpStatusCode", task.HttpStatusCode},
		{"updatedAt", time.Now().In(time.UTC)},
		{"length", task.Length},
		{"headers", task.Headers},
	}

	update := bson.D{
		{Key: "$set", Value: updateFields},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	switch err {
	case nil:
		return nil
	case mongo.ErrNoDocuments:
		return drivers.ErrTaskNotFound
	default:
		return err
	}
}

func (r *TaskRepository) TaskById(ctx context.Context, id string) (*models.Task, error) {
	task := new(models.Task)

	filter := bson.D{
		{"id", id},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&task)

	switch err {
	case nil:
		return task, nil
	case mongo.ErrNoDocuments:
		return nil, drivers.ErrTaskNotFound
	default:
		return nil, err
	}

}
