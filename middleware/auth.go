package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Auth header not present"))
			return
		}

		authContent := strings.Split(authHeader, "Bearer ")
		if len(authContent) != 2 {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Bearer token not properly formatted"))
			return
		}

		url := fmt.Sprintf(
			"https://www.googleapis.com/oauth2/v1/tokeninfo?access_token=%s",
			authContent[1],
		)

		res, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}

		if res.StatusCode != http.StatusOK {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
			return
		}

		c.Next()
	}
}
