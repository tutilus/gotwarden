package models

import "time"

type Attachment struct {
	UUID       string    `db:"uuid" json:"Id"`
	CipherUUID string    `db:"cipher_uuid" json:""`
	URL        string    `db:"url" json:"Url"`
	Filename   string    `db:"filename" json:"Filename"`
	Size       int       `dv:"size" json:"Size"`
	File       []byte    `db:"file" json:"File"`
	UpateAt    time.Time `db:"update_at" json:"RevisionDate"`
}

//TODO: Créer la table attachment + le model
