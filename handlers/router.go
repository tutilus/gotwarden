package handlers

import (
	"gotwarden/models"
	"gotwarden/util"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// WardenCtx is the db datastore functions
type WardenCtx struct {
	Db              models.Datastore
	Port            string
	SecretPhrase    string
	Validity        time.Duration
	ValidityRefresh time.Duration
	IdentityURL     string
}

// Init is the constructor for WardenCtx
func Init(conf *util.Config) (*WardenCtx, error) {

	db, err := models.NewDB(conf.Db.GetType(), conf.Db.GetConnect())
	if err != nil {
		return nil, err
	}

	return &WardenCtx{
		db,
		conf.Port,
		"Phrase",
		time.Hour,
		// MaxRefreshTime arbitrary at 5 min after token expires
		time.Hour + time.Duration(5*60),
		"/identity",
	}, nil
}

// Router define all the handles
func (ctx *WardenCtx) Router() http.Handler {

	authMiddleware, err := jwt.New(JwtMiddleware(ctx))

	if err != nil {
		log.Fatal("JWT Error: " + err.Error())
	}

	r := gin.Default()

	accounts := r.Group("/api/accounts")
	{
		accounts.POST("/register", ctx.SignUp)
		accounts.POST("/prelogin", ctx.PreLogin)
		accounts.Use(authMiddleware.MiddlewareFunc())
		{
			accounts.POST("/keys", ctx.GetKeys).Use(authMiddleware.MiddlewareFunc())
		}
	}

	auth := r.Group("/api")
	// Middleware JWT
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/sync", ctx.Synchronize)
		auth.POST("/ciphers", ctx.SaveCipher)
		auth.PUT("/ciphers/:uuid", ctx.SaveCipher)
		auth.PUT("/ciphers/:uuid/delete", ctx.DeleteCipher)
		auth.POST("/folders", ctx.SaveFolder)
		auth.PUT("/folders/:uuid", ctx.SaveFolder)
		auth.DELETE("/folders/:uuid", ctx.DeleteFolder)
		auth.PUT("/devices/identifier/:uuid/clear-token", ctx.ClearToken)
		auth.PUT("/devices/identifier/:uuid/token", ctx.UpdateToken)
	}

	identity := r.Group("/identity")
	{
		identity.POST("/connect/token", func(c *gin.Context) {

			log.Printf("grant_type: %s", c.PostForm("grant_type"))

			switch c.PostForm("grant_type") {
			case "refresh_token":
				// Refresh the token
				refreshToken := c.PostForm("refresh_token")

				if refreshToken == "" {
					c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("'refresh_token' cannot be blank"))
				}

				// Device findByRefreshToken()
				authMiddleware.RefreshHandler(c)

			case "password":
				// Loginn Handler
				authMiddleware.LoginHandler(c)
			default:
				c.AbortWithStatusJSON(http.StatusBadRequest, FormattedError("grant_type should be 'password' or 'refresh_token'"))
			}
		})
	}

	r.GET("/icons/:domain/icon.png", func(c *gin.Context) {
		c.Redirect(http.StatusOK, "http://"+c.Param("domain")+"/favicon.ico")
	})

	notif := r.Group("/notifications")
	{
		notif.GET("/hub", ctx.NotifHub)
	}
	// Routes to help to build the API. Should be disabled at the end.
	admin := r.Group("/admin")
	{
		admin.GET("/users", ctx.ShowUsers)
		admin.GET("/devices", ctx.ShowDevices)
		admin.GET("/folders", ctx.ShowFolders)
		admin.GET("/ciphers", ctx.ShowCiphers)
	}

	return r
}

// WardenError is a nested msg error
type WardenError struct {
	Message []string `binding:"required"`
}

// FormattedError formats error as BitWarden expected
func FormattedError(msg string) map[string]interface{} {
	return gin.H{
		"ValidationErrors": WardenError{
			Message: []string{msg},
		},
		"Object": "error",
	}
}
