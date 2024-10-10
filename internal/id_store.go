package internal

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

const (
	fileName = "id.data"
)

type IDStore struct {
	data map[string]struct{}
	m    sync.Mutex
	file *os.File
}

func NewIDStore() (*IDStore, error) {
	store := &IDStore{
		data: make(map[string]struct{}),
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	filePath := filepath.Join(dir, fileName)

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	store.file = file

	if err := lockFile(store.file); err != nil {
		log.Println(err)
		return nil, err
	}
	defer unlockFile(store.file)

	// load existing IDs from the file
	if err := store.loadFromFile(); err != nil {
		return nil, err
	}

	return store, nil
}

func (store *IDStore) GetId() string {
	store.m.Lock()
	defer store.m.Unlock()

	if err := lockFile(store.file); err != nil {
		log.Println(err)
		return ""
	}
	defer unlockFile(store.file)

	var id string

	for {
		id = uuid.NewString()

		// check if the ID is unique in the map
		if _, exists := store.data[id]; !exists {
			store.data[id] = struct{}{}
			break
		}
	}

	if err := store.saveToFile(id); err != nil {
		log.Println(err)
		return ""
	}

	return id
}

func (store *IDStore) FreeId(id string) error {
	store.m.Lock()
	defer store.m.Unlock()

	if err := lockFile(store.file); err != nil {
		return err
	}
	defer unlockFile(store.file)

	if _, exists := store.data[id]; exists {
		delete(store.data, id)

		if err := store.rewriteFile(); err != nil {
			return err
		}

		return nil
	}

	return errors.New("ID not found")
}

func (store *IDStore) Close() {
	err := store.file.Close()
	if err != nil {
		log.Println(err)
		return
	}
}

func (store *IDStore) loadFromFile() error {
	store.m.Lock()
	defer store.m.Unlock()

	if err := lockFile(store.file); err != nil {
		log.Println(err)
		return err
	}
	defer unlockFile(store.file)

	_, err := store.file.Seek(0, 0)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(store.file)
	for scanner.Scan() {
		id := strings.TrimSpace(scanner.Text())
		if id != "" {
			store.data[id] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (store *IDStore) saveToFile(id string) error {
	if _, err := store.file.WriteString(id + "\n"); err != nil {
		return err
	}

	return nil
}

// TODO: potential performance bottleneck
func (store *IDStore) rewriteFile() error {
	if err := store.file.Truncate(0); err != nil {
		return err
	}

	if _, err := store.file.Seek(0, 0); err != nil {
		return err
	}

	for id := range store.data {
		_, err := store.file.WriteString(id + "\n")
		if err != nil {
			return err
		}
	}

	// ensure the file is flushed
	if err := store.file.Sync(); err != nil {
		return err
	}

	return nil
}

// Lock the file before performing any operation
func lockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_EX)
}

// Unlock the file after operation is done
func unlockFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
