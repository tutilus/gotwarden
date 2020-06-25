package handlers

import (
	"gotwarden/models"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Identity contains all fields provided to Authentificator JWT function
type Identity struct {
	ClientID         string `form:"client_id"`
	GrantType        string `form:"grant_type"`
	DeviceIdentifier string `form:"deviceIdentifier"`
	DeviceName       string `form:"deviceName"`
	DeviceType       string `form:"deviceType"`
	Password         string `form:"password"`
	Scope            string `form:"scope"`
	Username         string `form:"username"`
	PushToken        string `form:"devicePushToken"`
}

// AccessToken is all the fields needed by the official clients
type AccessToken struct {
	Sub           string   `json:"sub"`
	Premium       string   `json:"premium"`
	Name          string   `json:"name"`
	Email         string   `json:"email"`
	EmailVerified string   `json:"email_verified"`
	Sstamp        string   `json:"sstamp"`
	Device        string   `json:"device"`
	Scope         []string `json:"scope"`
}

// JwtMiddleware is the middleware to manage Jwt
func JwtMiddleware(ctx *WardenCtx) *jwt.GinJWTMiddleware {

	return &jwt.GinJWTMiddleware{
		Realm:       "gotwarden JWT",
		Key:         []byte(ctx.SecretPhrase),
		IdentityKey: "sub",
		Timeout:     ctx.Validity,
		MaxRefresh:  ctx.RefeshValidity,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*AccessToken); ok {
				return jwt.MapClaims{
					"nbf":            time.Now().Add(-time.Minute * time.Duration(2)).Unix(),
					"iss":            ctx.IdentityURL,
					"exp":            time.Now().Add(ctx.Validity).Unix(),
					"sub":            v.Sub,
					"premium":        v.Premium,
					"name":           v.Name,
					"email":          v.Email,
					"email_verified": v.EmailVerified,
					"sstamp":         v.Sstamp,
					"device":         v.Device,
					"scope":          v.Scope,
					"amr":            []string{"Application"},
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var identity Identity

			if err := c.ShouldBind(&identity); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if identity.Scope != "api offline_access" {
				return "", jwt.ErrMissingLoginValues
			}

			// Check is user exists
			if user, err := ctx.Db.GetUserFromEmail(identity.Username); err == nil {
				// Verify Password
				if user.CheckPassword(identity.Password) {
					// TODO: Two-factor

					// Get the Device for the DeviceIdentifier attach to the user
					d := ctx.Db.GetDevice(identity.DeviceIdentifier)

					if d == nil {
						// If Device not found, create one
						d = models.NewDevice(identity.DeviceIdentifier, identity.DeviceName, identity.DeviceType, user.UUID)
						err = ctx.Db.AddDevice(d)
						if err != nil {
							log.Printf("Cannot insert Device %s : %s", identity.DeviceIdentifier, err)
						}
					} else {
						d.Type = identity.DeviceType
						d.Name = identity.DeviceName
						if identity.PushToken != "" {
							d.PushToken = identity.PushToken
						}
						if err = ctx.Db.SaveDevice(d); err != nil {
							return nil, jwt.ErrFailedAuthentication
						}
					}

					return &AccessToken{
						Sub:           user.UUID,
						Name:          user.Name,
						Email:         user.Email,
						EmailVerified: strconv.FormatBool(user.EmailVerified),
						Sstamp:        user.SecurityStamp,
						Device:        d.UUID,
						Scope:         strings.Split(identity.Scope, " "),
					}, nil
				}
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"access_token":  token,
				"expires_in":    int(ctx.Validity.Seconds()),
				"token_type":    "Bearer",
				"refresh_token": c.PostForm("refresh_token"),
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// Get the user info
			var identity Identity

			// Récupérer depuis le contexte et la Db les infos nécessaires
			if err := c.ShouldBind(&identity); err == nil {
				if user, err := ctx.Db.GetUserFromEmail(identity.Username); err == nil {
					d := ctx.Db.GetDevice(identity.DeviceIdentifier)
					if d != nil {
						// Save token
						d.AccessToken = token
						ctx.Db.SaveDevice(d)
						c.JSON(code, gin.H{
							"access_token":  token,
							"expire_in":     int(ctx.Validity.Seconds()),
							"token_type":    "Bearer",
							"refresh_token": d.RefreshToken,
							"Key":           user.Key,
						})
						return
					}
				}
			}
			c.AbortWithStatusJSON(500, FormattedError("Inner database access issue"))
		},
	}
}
