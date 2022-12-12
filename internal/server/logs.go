package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbecker/gin-api-demo/util"
	log "github.com/sirupsen/logrus"
)

func JSONLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := util.GetDurationInMillseconds(start)

		entry := log.WithFields(log.Fields{
			"client_ip":  util.GetClientIP(c),
			"duration":   duration,
			"method":     c.Request.Method,
			"path":       c.Request.RequestURI,
			"status":     c.Writer.Status(),
			"user_id":    util.GetUserID(c),
			"referrer":   c.Request.Referer(),
			"request_id": c.Writer.Header().Get("Request-Id"),
		})

		msg := strings.Join(delete_empty([]string{http.StatusText(c.Writer.Status()), c.Errors.String()}), " - ")
		switch status := c.Writer.Status(); {
		case status >= http.StatusInternalServerError:
			entry.Error(msg)
		case status >= http.StatusBadRequest:
			entry.Info(msg)
		default:
			entry.Debug(msg)
		}
	}
}

func delete_empty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
