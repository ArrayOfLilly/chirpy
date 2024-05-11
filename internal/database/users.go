package database

// Structure of a "User" as a databasae entry
type User struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
}

// CreateUser creates a new user in the database.
//
// It takes an email string as a parameter, which specifies the email address of the user.
// It returns a User struct and an error. 
// If there is an error while loading the database or writing to the database, or the specified email is invalid 
// it returns an empty User struct and the corresponding error. Otherwise, it returns the created user and a nil error.
func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	
	user := User{
		ID:   id,
		Email: email,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}