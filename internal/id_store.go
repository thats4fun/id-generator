package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	fileName = "id_store.db"
)

type IDStore struct {
	db *sql.DB
	m  sync.Mutex
}

func NewIDStore() (*IDStore, error) {
	store := &IDStore{}

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// TODO: can be configurable
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	store.db = db

	err = store.createTableIfNotExists()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return store, nil
}

func (store *IDStore) GetId() string {
	store.m.Lock()
	defer store.m.Unlock()

	var id string

	for {
		id = uuid.NewString()

		// Check if the ID is unique in the table
		exists, err := store.exists(id)
		if err != nil {
			log.Println(err)
			return ""
		}
		if !exists {
			break
		}
	}

	// Save the new ID to the database
	err := store.saveToDB(id)
	if err != nil {
		log.Println(err)
		return ""
	}

	return id
}

func (store *IDStore) FreeId(id string) error {
	store.m.Lock()
	defer store.m.Unlock()

	exists, err := store.exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("ID not found")
	}

	// Delete the ID from the database
	err = store.deleteFromDB(id)
	if err != nil {
		return err
	}

	return nil
}

func (store *IDStore) Close() {
	err := store.db.Close()
	if err != nil {
		log.Println(err)
	}
}

func (store *IDStore) createTableIfNotExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS ids (
		id TEXT PRIMARY KEY
	)`
	_, err := store.db.Exec(query)
	return err
}

func (store *IDStore) exists(id string) (bool, error) {
	query := `SELECT 1 FROM ids WHERE id = ? LIMIT 1`

	var exists int
	err := store.db.QueryRow(query, id).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}

	return exists == 1, nil
}

func (store *IDStore) saveToDB(id string) error {
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO ids (id) VALUES (?)`
	_, err = tx.Exec(query, id)
	if err != nil {
		rollbackErr := tx.Rollback()
		err = fmt.Errorf("transaction exec failed: %w", rollbackErr)
		return err
	}

	return tx.Commit()
}

func (store *IDStore) deleteFromDB(id string) error {
	query := `DELETE FROM ids WHERE id = ?`
	_, err := store.db.Exec(query, id)
	return err
}
