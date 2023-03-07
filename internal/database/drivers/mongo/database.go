package mongo

import (
	"context"
	"fmt"
	"github.com/akhmettolegen/proxy/internal/database/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 10 * time.Second
	retries           = 1
	collectionTask    = "task"
)

type Mongo struct {
	MongoURL string
	client   *mongo.Client
	dbname   string

	DB      *mongo.Database
	Context context.Context

	taskRepository *TaskRepository

	retries           int
	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration
}

func (m *Mongo) Name() string { return "Mongo" }

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	return &Mongo{
		MongoURL:          conf.URL,
		dbname:            conf.DataBaseName,
		retries:           retries,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	fmt.Printf("Connecting to: %s at %s\n", m.dbname, m.MongoURL)
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.MongoURL))
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbname)

	// убеждаемся что созданы все необходимые индексы
	return m.ensureIndexes()
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(m.Context)
}

func (m *Mongo) Task() drivers.TaskRepository {
	if m.taskRepository == nil {
		m.taskRepository = &TaskRepository{
			collection: m.DB.Collection(collectionTask),
		}
	}
	return m.taskRepository
}

func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if err := m.ensureTaskIndexes(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mongo) ensureTaskIndexes(ctx context.Context) error {
	col := m.DB.Collection(collectionTask)

	models := []mongo.IndexModel{
		{Keys: bson.M{"id": 1}},
		{Keys: bson.M{"createdAt": -1}},
	}

	opts := options.CreateIndexes().SetMaxTime(m.ensureIdxTimeout)
	_, err := col.Indexes().CreateMany(ctx, models, opts)

	return err
}
