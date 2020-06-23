package models

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/google/uuid"
)

// Device which access to GotWarden server
type Device struct {
	UUID           string    `db:"uuid"`
	Name           string    `db:"name" json:"deviceName"`
	Type           string    `db:"type" json:"deviceType"`
	PushToken      string    `db:"push_token"`
	AccessToken    string    `db:"access_token"`
	RefreshToken   string    `db:"refresh_token"`
	TokenExpiresAt time.Time `db:"token_expires_at"`
	UserUUID       string    `db:"user_uuid"`
	PrivateKey     []byte    `db:"private_key"`
}

// AllDevices get all the devices
func (db *DB) AllDevices() (*[]Device, error) {
	dd := []Device{}

	_, err := db.Select(&dd, "SELECT * FROM devices")

	return &dd, err
}

// GetDevice get device for specific uuid identifier
func (db *DB) GetDevice(uuid string) *Device {
	obj, err := db.DbMap.Get(Device{}, uuid)

	if obj == nil {
		log.Printf("Get Device error %s", err)
		return nil
	}
	return obj.(*Device)
}

// GetDeviceFromToken get device for specific token
func (db *DB) GetDeviceFromToken(token string) (*Device, error) {
	d := Device{}
	err := db.DbMap.SelectOne(&d, "SELECT FROM * WHERE access_token=?", token)

	return &d, err
}

// AddDevice persiste an object Device
func (db *DB) AddDevice(device *Device) error {
	return db.Insert(device)
}

// SaveDevice updates uuid device with the value into device provided
func (db *DB) SaveDevice(device *Device) error {
	_, err := db.Update(device)
	return err
}

// ---- Functions utilities ------- //

// NewDevice creates a new Device and persistes it
func NewDevice(deviceIdentifier, deviceName, deviceType, userID string) *Device {
	log.Printf("NewDevice : %s %s %s %s", deviceIdentifier, deviceName, deviceType, userID)
	return &Device{
		UUID:         deviceIdentifier,
		Name:         deviceName,
		Type:         deviceType,
		RefreshToken: NewTokenURLSafe(),
		UserUUID:     userID,
	}
}

// NewTokenURLSafe creates a token Url Safe
func NewTokenURLSafe() string {
	code := uuid.New()
	return base64.RawURLEncoding.EncodeToString([]byte(code.String()))
}
