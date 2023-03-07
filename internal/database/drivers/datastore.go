package drivers

type DataStore interface {
	Name() string
	Close() error
	Connect() error

	Task() TaskRepository
	Ping() error
}
