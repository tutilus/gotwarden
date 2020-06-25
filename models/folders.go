package models

import (
	"log"
	"time"
)

// Folder to organize the f
type Folder struct {
	UUID     string    `db:"uuid" json:"Id"`
	UserUUID string    `db:"user_uuid" json:"-"`
	Name     []byte    `db:"name" json:"Name" binding:"required"`
	UpdateAt time.Time `db:"update_at" json:"RevisionDate"`
	Object   string    `db:"-" json:"Object"`
}

// AllFolders gets all the folders for this user
func (db *DB) AllFolders() (*[]Folder, error) {
	var folders []Folder
	_, err := db.Select(&folders, "SELECT * FROM folders")

	return &folders, err
}

// GetFoldersByUserUUID gets all the folders for this user
func (db *DB) GetFoldersByUserUUID(uuid string) (*[]Folder, error) {
	var folders []Folder
	_, err := db.Select(&folders, "SELECT * FROM folders WHERE user_uuid=?", uuid)

	return &folders, err
}

// GetFolder gets folder data from database
func (db *DB) GetFolder(uuid string) *Folder {
	obj, err := db.DbMap.Get(Folder{}, uuid)

	if err != nil {
		log.Printf("Failed to get the folder %s", uuid)
		return nil
	}
	f := obj.(*Folder)
	return f
}

// AddFolder saves a new f
func (db *DB) AddFolder(f *Folder) error {
	return db.Insert(f)
}

// SaveFolder updates existing f
func (db *DB) SaveFolder(f *Folder) error {
	_, err := db.Update(f)
	return err
}

// DeleteFolder deletes f provided
func (db *DB) DeleteFolder(f *Folder) error {

	// Remove all the reference into the ciphers
	ciphers, err := db.GetCiphersByFolderUUID(f.UUID)
	if err != nil {
		log.Printf("Impossible to get ciphers for this uuid and to deference it")
		return err
	}
	for _, cipher := range *ciphers {
		cipher.FolderUUID = ""
		db.SaveCipher(&cipher)
	}
	_, err = db.Delete(f)
	return err
}

// Jsonify creates object ready to send back
func (f *Folder) Jsonify() *Folder {
	// Copy folder
	fo := *f
	fo.Object = "folder"
	return &fo
}
