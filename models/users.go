package models

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// User data structure
type User struct {
	UUID             string    `db:"uuid"`
	Email            string    `db:"email"`
	EmailVerified    bool      `db:"email_verified"`
	Premium          bool      `db:"premium"`
	Name             string    `db:"name"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	PasswordHint     string    `db:"password_hint" json:"MasterPasswordHint"`
	Key              string    `db:"key_pass"`
	Culture          string    `db:"culture"`
	TwoFactorEnabled bool      `db:"-"`
	PrivateKey       []byte    `db:"private_key"`
	PublicKey        []byte    `db:"public_key" json:"-"`
	TotpSecret       string    `db:"totp_secret" json:"-"`
	SecurityStamp    string    `db:"security_stamp"`
	CreatedAt        time.Time `db:"created_at" json:"-"`
	Kdf              int       `db:"kdf"`
	KdfIterations    int       `db:"kdf_iterations" binding:"required"`
	Organizations    []string  `db:"-"`
	Object           string    `db:"-"`
}

// AllUsers get all the users
func (db *DB) AllUsers() (*[]User, error) {
	uu := []User{}

	_, err := db.Select(&uu, "SELECT * FROM users")

	return &uu, err
}

// GetUser get a user
func (db *DB) GetUser(uuid string) *User {
	obj, err := db.DbMap.Get(User{}, uuid)

	if obj == nil {
		log.Printf("Get User error %s", err)
		return nil
	}
	u := obj.(*User)
	// Add some fields
	u.Object = "profile"
	u.TwoFactorEnabled = u.TotpSecret != ""
	return u
}

// SaveUser updates uuid user with the value into user provided
func (db *DB) SaveUser(u *User) error {
	_, err := db.Update(u)
	return err
}

// GetUserFromEmail get a user
func (db *DB) GetUserFromEmail(email string) (*User, error) {
	u := User{}

	err := db.SelectOne(&u, "SELECT * FROM users WHERE email=?", email)
	log.Printf("GetUserFromEmail : %s", err)
	return &u, err
}

// AddUser persistes the User provided
func (db *DB) AddUser(user *User) error {
	return db.Insert(user)
}

// CheckPassword checks if password is valid
func (u *User) CheckPassword(password string) bool {
	return u.PasswordHash == password
}

// NewUser declares a new user
func NewUser(name, email, masterPasswordHash, masterPasswordHint, key string, kdf, kdfIterations int) *User {
	return &User{
		UUID:          uuid.New().String(),
		Name:          name,
		Email:         email,
		EmailVerified: true,
		Culture:       "en-US",
		Premium:       true,
		PasswordHash:  masterPasswordHash,
		PasswordHint:  masterPasswordHint,
		Key:           key,
		Kdf:           kdf,
		KdfIterations: kdfIterations,
	}
}
