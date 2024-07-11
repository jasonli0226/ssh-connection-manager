package domain

type Connection struct {
	Alias    string
	Host     string
	User     string
	Password string
	Port     int
}

type ConnectionRepository interface {
	Add(Connection) error
	Get(alias string) (Connection, error)
	List() ([]Connection, error)
	Delete(alias string) error
}
