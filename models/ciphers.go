package models

import (
	"gotwarden/util"
	"log"
	"time"
)

// CipherData contains the data to store for a cipher
type CipherData struct {
	UUID             string    `db:"uuid"`
	UserUUID         string    `db:"user_uuid"`
	FolderUUID       string    `db:"folder_uuid"`
	OrganizationUUID string    `db:"organization_uuid"`
	Type             int       `db:"type"`
	Data             []byte    `db:"data"`
	Favorite         bool      `db:"favorite"`
	Name             string    `db:"name"`
	Notes            []byte    `db:"notes"`
	Fields           []byte    `db:"fields"`
	Login            []byte    `db:"login"`
	Card             []byte    `db:"card"`
	Identity         []byte    `db:"identity"`
	SecureNote       []byte    `db:"securenote"`
	PasswordHistory  []byte    `db:"passwordhistory"`
	UpdateAt         time.Time `db:"update_at"`
}

// CipherObject is a components into Warden server
type CipherObject struct {
	UUID                string `json:"Id"`
	FolderUUID          string `json:"FolderId"`
	OrganizationUUID    string `json:"OrganizationId"`
	OrganizationUseTotp bool
	Type                int
	Favorite            bool
	Attachments         []interface{}
	Name                string
	Totp                interface{}
	Notes               interface{}
	Fields              []interface{}
	Login               interface{}
	Card                interface{}
	Identity            interface{}
	SecureNote          interface{}
	PasswordHistory     []interface{}
	RevisionDate        string
	Edit                bool
	Object              string
}

// AllCiphers gets all the ciphers for this user
func (db *DB) AllCiphers() (*[]CipherData, error) {
	var ciphers []CipherData
	_, err := db.Select(&ciphers, "SELECT * FROM ciphers")

	return &ciphers, err
}

// GetCiphersByUserUUID gets all the ciphers for this user
func (db *DB) GetCiphersByUserUUID(uuid string) (*[]CipherData, error) {
	var ciphers []CipherData
	_, err := db.Select(&ciphers, "SELECT * FROM ciphers WHERE user_uuid=?", uuid)

	return &ciphers, err
}

// GetCipher gets a specific cipher
func (db *DB) GetCipher(uuid string) *CipherData {
	obj, err := db.DbMap.Get(CipherData{}, uuid)

	if obj == nil {
		log.Printf("Get User error %s", err)
		return nil
	}
	c := obj.(*CipherData)
	return c
}

// AddCipher saves a new cipher
func (db *DB) AddCipher(cipher *CipherData) error {
	return db.Insert(cipher)
}

// SaveCipher updates existing cipher
func (db *DB) SaveCipher(cipher *CipherData) error {
	_, err := db.Update(cipher)
	return err
}

// DeleteCipher deletes cipher provided
func (db *DB) DeleteCipher(cipher *CipherData) error {
	_, err := db.Delete(cipher)
	return err
}

// Jsonify provides a Struct to send back as Json
func (cd *CipherData) Jsonify() *CipherObject {
	return &CipherObject{
		UUID:             cd.UUID,
		FolderUUID:       cd.FolderUUID,
		Favorite:         cd.Favorite,
		Type:             cd.Type,
		OrganizationUUID: cd.OrganizationUUID,
		Login:            util.UnmarshalObject(cd.Login),
		Fields:           util.UnmarshalArray(cd.Fields),
		PasswordHistory:  util.UnmarshalArray(cd.PasswordHistory),
		RevisionDate:     cd.UpdateAt.Format("2020-12-31T12:01:10.000000Z"),
		//TODO		Attachments: db.GetAttachments(cd.UUID),
		OrganizationUseTotp: false,
		Name:                cd.Name,
		Notes:               util.UnmarshalObject(cd.Notes),
		Card:                util.UnmarshalObject(cd.Card),
		Identity:            util.UnmarshalObject(cd.Identity),
		SecureNote:          util.UnmarshalObject(cd.SecureNote),
		Edit:                true,
		Object:              "cipher",
	}
}
