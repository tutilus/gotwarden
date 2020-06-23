package handlers

import (
	"gotwarden/models"
	"gotwarden/util"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login contains data when log in
type Login struct {
	Email         string `json:"email" binding:"required"`
	Name          string `json:"name"`
	PasswordHash  string `json:"masterPasswordHash" binding:"required"`
	PasswordHint  string `json:"masterPasswordHint"`
	Key           string `json:"key" binding:"required"`
	Kdf           int    `json:"kdf"`
	KdfIterations int    `json:"kdfIterations" binding:"required"`
}

// PreLogin contains data needed to prepare login
type PreLogin struct {
	Email string `json:"email" binding:"required"`
}

// Cipher object
type Cipher struct {
	Type            int           `json:"type" binding:"required"`
	FolderID        string        `json:"folderId"`
	OrganisationID  string        `json:"organisationId"`
	Name            string        `json:"name" binding:"required"`
	Notes           string        `json:"notes"`
	Favorite        bool          `json:"favorite"`
	Fields          []interface{} `json:"fields"`
	Login           interface{}   `json:"login"`
	PasswordHistory []interface{} `json:"passwordhistory"`
	SecureNote      interface{}   `json:"securenote"`
	Card            interface{}   `json:"card"`
	Identity        interface{}   `json:"identity"`
}

// SignUp create a new User if doesn't exist (based on email)
func (ctx *WardenCtx) SignUp(c *gin.Context) {
	var l Login

	err := c.BindJSON(&l)
	if err != nil {
		log.Printf("Binding error %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("Bad data request provided"))
		return
	}

	// Check if a user with this profile exists (based on email)
	_, err = ctx.Db.GetUserFromEmail(l.Email)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("A user with this email already exists"))
		return
	}

	// Create new user based on the profile data
	u := models.NewUser(l.Name, l.Email, l.PasswordHash, l.PasswordHint, l.Key, l.Kdf, l.KdfIterations)
	err = ctx.Db.AddUser(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Cannot register this new user"))
		return
	}
}

// PreLogin gets info needed for login
func (ctx *WardenCtx) PreLogin(c *gin.Context) {
	var login PreLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, FormattedError("Bad data request provided"))
	} else {
		u, err := ctx.Db.GetUserFromEmail(login.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("Bad data request provided"))
		} else {
			// Response Kdf and KdfIterations used
			c.JSON(200, gin.H{
				"Kdf":           u.Kdf,
				"KdfIterations": u.KdfIterations,
			})
		}
	}
}

// Synchronize asks server to provide all data
func (ctx *WardenCtx) Synchronize(c *gin.Context) {
	claim := jwt.ExtractClaims(c)

	if claim["sub"] == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("JWT fail to reach user"))
		return
	}

	// Get user from the id into the token
	u := ctx.Db.GetUser(claim["sub"].(string))

	ciphers, _ := ctx.Db.GetCiphersByUserUUID(u.UUID)
	var cj []interface{}
	for _, cipher := range *ciphers {
		cj = append(cj, cipher.Jsonify())
	}

	folders, _ := ctx.Db.GetFoldersByUserUUID(u.UUID)

	c.JSON(http.StatusOK, gin.H{
		"Profile": u,
		"Folders": folders,
		"Ciphers": cj,
		"Domains": models.Domains{
			EquivalentDomains: nil,
			Object:            "domains",
		},
		"Object": "sync",
	})
}

// GetKeys provides public and encrypted private keys
func (ctx *WardenCtx) GetKeys(c *gin.Context) {
	claim := jwt.ExtractClaims(c)
	if claim["sub"] == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Middleware JWT failed"))
		return
	}
	// Get user from the id into the token
	u := ctx.Db.GetUser(claim["sub"].(string))
	u.PrivateKey = []byte(c.PostForm("encryptedPrivateKey"))
	u.PublicKey = []byte(c.PostForm("publicKey"))

	ctx.Db.SaveUser(u)
}

// SaveCipher creates a new Cipher component
func (ctx *WardenCtx) SaveCipher(c *gin.Context) {
	claim := jwt.ExtractClaims(c)

	// Get user from the id into the token
	userUUID := claim["sub"].(string)

	cipher := Cipher{}

	if err := c.ShouldBindJSON(&cipher); err != nil {
		log.Printf("A : %s", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError(err.Error()))
		return
	}
	// Check if the folder id exist and is belonging to the user
	if cipher.FolderID != "" {
		f := ctx.Db.GetFolder(cipher.FolderID)
		if f == nil {
			log.Print("B")
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Database failed to find a folder for this id"))
			return
		}
		if f.UserUUID != userUUID {
			log.Printf("C")
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Invalid folder"))
			return
		}
	}

	// Save the cipher
	var cd *models.CipherData

	switch c.Request.Method {
	case "POST":
		// Create a new cipher with uuid.New() UUID
		cd = cipher.ToCipherData(userUUID, uuid.New().String())
		err := ctx.Db.AddCipher(cd)
		if err != nil {
			log.Printf("D1 : %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Database failed to add cipher"))
			return
		}
	case "PUT":
		// Update existing cipher with uuid from parameter
		cd = cipher.ToCipherData(userUUID, c.Param("uuid"))
		err := ctx.Db.SaveCipher(cd)
		if err != nil {
			log.Printf("D2 : %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Database failed to save cipher"))
			return
		}
	default:
		log.Print("default")
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("Unexpected method found"))
		return
	}

	c.JSON(http.StatusOK, cd.Jsonify())

}

// DeleteCipher action
func (ctx *WardenCtx) DeleteCipher(c *gin.Context) {
	claim := jwt.ExtractClaims(c)
	// Get user from the id into the token
	userUUID := claim["sub"].(string)
	log.Printf("Delete Cipher for user %s", userUUID)

	uuid := c.Param("uuid")
	log.Printf("Cipher UUID %s", uuid)

	cipher := ctx.Db.GetCipher(uuid)
	if cipher == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to get the designed cipher"))
		return
	}

	ctx.Db.DeleteCipher(cipher)
}

// SaveFolder creates or updates a Folder
func (ctx *WardenCtx) SaveFolder(c *gin.Context) {
	claim := jwt.ExtractClaims(c)

	f := models.Folder{}

	// Get user from the id into the token
	f.UserUUID = claim["sub"].(string)

	if err := c.ShouldBindJSON(&f); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError(err.Error()))
		return
	}

	switch c.Request.Method {
	case "POST":
		// Create a new folder with uuid.New() UUID
		f.UUID = uuid.New().String()
		err := ctx.Db.AddFolder(&f)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Database failed to add cipher"))
			return
		}
	case "PUT":
		// Update existing cipher with uuid from parameter
		f.UUID = c.Param("uuid")
		err := ctx.Db.SaveFolder(&f)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Database failed to save cipher"))
			return
		}
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("Unexpected method found"))
		return
	}

	c.JSON(http.StatusOK, f.Jsonify())

}

// DeleteFolder action
func (ctx *WardenCtx) DeleteFolder(c *gin.Context) {
	claim := jwt.ExtractClaims(c)
	// Get user from the id into the token
	userUUID := claim["sub"].(string)
	log.Printf("Delete Folder for user %s", userUUID)

	uuid := c.Param("uuid")
	log.Printf("Folder UUID %s", uuid)

	f := ctx.Db.GetFolder(uuid)
	if f == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to get the designed cipher"))
		return
	}

	ctx.Db.DeleteFolder(f)
}

// ClearToken clears token for the device
func (ctx *WardenCtx) ClearToken(c *gin.Context) {
	claim := jwt.ExtractClaims(c)
	// Get device UUID from the token
	deviceUUID := claim["device"].(string)

	device := ctx.Db.GetDevice(deviceUUID)

	if device == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to get data for the device"))
		return
	}
	device.PushToken = ""

	err := ctx.Db.SaveDevice(device)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to update device"))
		return
	}
}

// UpdateToken updates device with provided token
func (ctx *WardenCtx) UpdateToken(c *gin.Context) {
	claim := jwt.ExtractClaims(c)
	// Get device UUID from the token
	deviceUUID := claim["device"].(string)

	device := ctx.Db.GetDevice(deviceUUID)

	if device == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to get data for the device"))
		return
	}
	device.PushToken = c.PostForm("pushtoken")

	err := ctx.Db.SaveDevice(device)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, FormattedError("Failed to update device"))
		return
	}
}

// ToCipherData populates data fields for Cipher
func (c *Cipher) ToCipherData(uu string, uuid string) *models.CipherData {

	return &models.CipherData{
		UUID:             uuid,
		UserUUID:         uu,
		FolderUUID:       c.FolderID,
		OrganizationUUID: c.OrganisationID,
		Favorite:         c.Favorite,
		Type:             c.Type,
		Name:             c.Name,
		Notes:            util.MarshalObject(c.Notes),
		Fields:           util.MarshalArray(c.Fields),
		Login:            util.MarshalObject(c.Login),
		PasswordHistory:  util.MarshalArray(c.PasswordHistory),
		SecureNote:       util.MarshalObject(c.SecureNote),
		Card:             util.MarshalObject(c.Card),
		Identity:         util.MarshalObject(c.Identity),
		UpdateAt:         time.Now(),
	}
}
