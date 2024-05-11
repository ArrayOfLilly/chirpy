package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Database connection
type DB struct {
	path string

	// A RWMutex is a reader/writer mutual exclusion lock. 
	// The lock can be held by an arbitrary number of readers or a single writer. 
	mu   *sync.RWMutex
}

// Structure of the database
// Sample:
// {
// 	"chirps":
// 	{
// 		"1":
// 		{
// 			"id": 1,
// 			"body": "sample text 1"
// 		},
// 		"2": 
// 		{
// 			"id": 2,
// 			"body": "sample text 2"
// 		},
// 	},
// 	"users":
// 	{
// 		"1":
// 		{
// 			"id": 1,
// 			"eamil": "example1@mail.com"
// 		},
// 		"2": 
// 		{
// 			"id": 2,
// 			"email": "example2@mail.com"
// 		},
// 	},
// }
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

// NewDB creates a new database connection.
//
// It takes a path string as a parameter, which specifies the path to the database file.
// It returns a pointer to a DB struct and an error.
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// ensureDB checks if the database file exists and creates it if it doesn't.
//
// It reads the database file using os.ReadFile. 
// If the file doesn't exist, it calls the createDB function to create the database.
// The function returns an error if there was an issue reading the file or creating the database.
// Returns:
// - error: An error if there was an issue reading the file or creating the database.
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

// loadDB loads the database structure from the specified file path.
//
// It acquires a read lock on the database mutex to ensure thread safety.
// The function initializes an empty DBStructure and attempts to read the file
// specified by the database path. 
// If the file does not exist, it returns the empty DBStructure and the os.ErrNotExist error.
// If the file exists, it attempts to unmarshal the JSON data into the DBStructure. 
// If there is an error during unmarshaling, it returns the empty DBStructure and the error.
// If the file is successfully read and unmarshaled, it returns the DBStructure and a nil error.
//
// Parameters:
// - None
//
// Returns:
// - DBStructure: The loaded database structure.
// - error: An error if there was an issue reading the file or unmarshaling the data.
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

// writeDB writes the given database structure to the database file.
//
// It takes a DBStructure as a parameter, which represents the structure of the database.
// It returns an error if there was an issue writing to the database.
//
// The function acquires a write lock on the database mutex to ensure thread safety.
// It marshals the DBStructure into JSON format and writes it to the database file specified by the path.
// The file is created with read and write permissions for the owner only.
// If there is an error during marshaling or writing, it returns the corresponding error.
// Otherwise, it returns nil.
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}