package models

import (

	// Database drivers

	"database/sql"
	"log"

	// SQLite3 drivers
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

// Datastore functions to manage User
type Datastore interface {
	AllUsers() (*[]User, error)
	AddUser(user *User) error
	GetUserFromEmail(email string) (*User, error)
	GetUser(uuid string) *User
	SaveUser(u *User) error
	AllDevices() (*[]Device, error)
	GetDevice(uuid string) *Device
	GetDeviceFromToken(token string) (*Device, error)
	AddDevice(device *Device) error
	SaveDevice(device *Device) error
	GetFolder(uuid string) *Folder
	GetFoldersByUserUUID(uuid string) (*[]Folder, error)
	AllFolders() (*[]Folder, error)
	AddFolder(f *Folder) error
	SaveFolder(f *Folder) error
	DeleteFolder(f *Folder) error
	AllCiphers() (*[]CipherData, error)
	AddCipher(cipher *CipherData) error
	GetCiphersByUserUUID(uuid string) (*[]CipherData, error)
	SaveCipher(cipher *CipherData) error
	DeleteCipher(cipher *CipherData) error
	GetCipher(uuid string) *CipherData
}

// DB injector
type DB struct {
	*gorp.DbMap
}

// NewDB create a new DB for the dataSource provided
func NewDB(typeDb, connectDb string) (*DB, error) {
	db, err := sql.Open(typeDb, connectDb)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(User{}, "users").SetKeys(false, "UUID")
	dbmap.AddTableWithName(Device{}, "devices").SetKeys(false, "UUID")
	dbmap.AddTableWithName(Folder{}, "folders").SetKeys(false, "UUID")
	dbmap.AddTableWithName(CipherData{}, "ciphers").SetKeys(false, "UUID")
	dbmap.AddTableWithName(Attachment{}, "attachments").SetKeys(false, "UUID")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalf("Database init failed with error %s", err)
	}
	return &DB{dbmap}, nil
}
