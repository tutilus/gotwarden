package models

import (
	"log"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// AttachmentData is the data struct for the model
type AttachmentData struct {
	UUID       string    `db:"uuid"`
	CipherUUID string    `db:"cipher_uuid"`
	Filename   string    `db:"filename"`
	Size       int       `dv:"size"`
	File       []byte    `db:"file"`
	UpdateAt   time.Time `db:"update_at"`
}

// AttachmentObject is the struct to manage internally the Attachment (ie Response)
type AttachmentObject struct {
	UUID         string `json:"Id"`
	URL          string `json:"Url"`
	Filename     string
	Size         int
	SizeName     string
	File         []byte
	RevisionDate string
	Object       string
}

// AllAttachments gets all the Attachments for this user
func (db *DB) AllAttachments() (*[]AttachmentData, error) {
	var attachments []AttachmentData
	_, err := db.Select(&attachments, "SELECT * FROM attachments")

	return &attachments, err
}

// GetAttachmentsByCypherUUID gets all the Attachments for this cipher
func (db *DB) GetAttachmentsByCypherUUID(uuid string) (*[]AttachmentData, error) {
	var attachments []AttachmentData
	_, err := db.Select(&attachments, "SELECT * FROM Attachments WHERE cypher_uuid=?", uuid)

	return &attachments, err
}

// GetAttachment gets Attachment data from database
func (db *DB) GetAttachment(uuid string) *AttachmentData {
	obj, err := db.DbMap.Get(AttachmentData{}, uuid)

	if err != nil {
		log.Printf("Failed to get the Attachment %s", uuid)
		return nil
	}
	a := obj.(*AttachmentData)
	return a
}

// AddAttachment saves a new f
func (db *DB) AddAttachment(a *AttachmentData) error {
	return db.Insert(a)
}

// DeleteAttachment deletes f provided
func (db *DB) DeleteAttachment(a *AttachmentData) error {
	_, err := db.Delete(a)
	return err
}

// Jsonify creates object ready to send back
func (a *AttachmentData) Jsonify() *AttachmentObject {
	return &AttachmentObject{
		UUID:         a.UUID,
		Filename:     a.Filename,
		File:         a.File,
		Size:         a.Size,
		SizeName:     humanize.Bytes(uint64(a.Size)),
		RevisionDate: a.UpdateAt.Format(time.RFC3339),
		Object:       "attachment",
	}
}
