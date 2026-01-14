package cmd

import (
	"log"
	"net/http"
	"time"

	"ewallet-ums/helpers"

	"github.com/gin-gonic/gin"
)

func (d *Dependency) MiddlewareValidateAuth(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		log.Println("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	_, err := d.UserRepo.GetUserSessionByToken(c.Request.Context(), auth)
	if err != nil {
		log.Println("failed to get user session on db: ", err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	claim, err := helpers.ValidateToken(c.Request.Context(), auth)
	if err != nil {
		log.Println(err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired: ", claim.ExpiresAt)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	c.Set("token", claim)
	c.Next()
}
