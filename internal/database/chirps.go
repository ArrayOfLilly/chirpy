package database

import (
	"os"
)

// Structure of a "Chirp" (limited length message) as a databasae entry
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp in the database.
//
// It takes a body string as a parameter, which represents the content of the chirp.
// It returns a Chirp struct and an error. 
// If there is an error while loading the database or writing to the database, 
// it returns an empty Chirp struct and the corresponding error. 
// Otherwise, it returns the created chirp and a nil error.
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}


// GetChirp retrieves a Chirp from the database based on its ID.
//
// Parameters:
// - id: the ID of the Chirp to retrieve.
//
// Returns:
// - Chirp: the Chirp with the given ID, or an empty Chirp if it does not exist.
// - error: an error if there was a problem loading the database or if the Chirp does not exist.
func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}


// GetChirps retrieves all the chirps from the database.
//
// It returns a slice of Chirp structs and an error if there was a problem loading the database.
// The slice contains all the chirps stored in the database.
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// createDB creates a new database structure and writes it to the database file.
//
// No parameters.
// Returns an error if there was an issue writing to the database.
func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

