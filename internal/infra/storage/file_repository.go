package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jasonli0226/ssh-connection-manager/internal/domain"
)

type FileConnectionRepository struct {
	filePath string
	mutex    sync.RWMutex
}

func NewFileConnectionRepository() *FileConnectionRepository {
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, ".ssh_manager_connections.json")
	return &FileConnectionRepository{filePath: filePath}
}

func (r *FileConnectionRepository) Add(conn domain.Connection) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	conns, err := r.readConnections()
	if err != nil {
		return err
	}

	conns = append(conns, conn)
	return r.writeConnections(conns)
}

func (r *FileConnectionRepository) Get(alias string) (domain.Connection, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	conns, err := r.readConnections()
	if err != nil {
		return domain.Connection{}, err
	}

	for _, conn := range conns {
		if conn.Alias == alias {
			return conn, nil
		}
	}

	return domain.Connection{}, fmt.Errorf("connection with alias %s not found", alias)
}

func (r *FileConnectionRepository) List() ([]domain.Connection, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.readConnections()
}

func (r *FileConnectionRepository) Delete(alias string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	conns, err := r.readConnections()
	if err != nil {
		return err
	}

	for i, conn := range conns {
		if conn.Alias == alias {
			conns = append(conns[:i], conns[i+1:]...)
			return r.writeConnections(conns)
		}
	}

	return fmt.Errorf("connection with alias %s not found", alias)
}

func (r *FileConnectionRepository) readConnections() ([]domain.Connection, error) {
	file, err := os.ReadFile(r.filePath)
	if os.IsNotExist(err) {
		return []domain.Connection{}, nil
	}
	if err != nil {
		return nil, err
	}

	var conns []domain.Connection
	err = json.Unmarshal(file, &conns)
	if err != nil {
		return nil, err
	}

	return conns, nil
}

func (r *FileConnectionRepository) writeConnections(conns []domain.Connection) error {
	data, err := json.MarshalIndent(conns, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0600)
}
