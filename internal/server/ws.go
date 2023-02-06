package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/auth"
	"github.com/stenic/ledger/internal/pkg/versions"
)

func wsHandler(authValidator auth.LedgerValidator) gin.HandlerFunc {
	logger := logrus.WithFields(logrus.Fields{
		"scope": "websockets",
	})

	var count = *versions.CountTotal(nil)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				count = *versions.CountTotal(nil)
				logger.WithField("count", count).Trace("Refreshed version count")
			}
		}
	}()

	return func(c *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			logger.Warn(err)
		}
		logger.Debug("Client connected")

		go func() {
			defer conn.Close()
			var lastSend = count

			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				logger.Error(err)
			}
			if err := authValidator.ValidateToken(string(msg)); err != nil {
				logger.Error(err)
				return
			}

			for {
				if lastSend != count {
					logger.Debug("Sending refreshVersions")
					err = wsutil.WriteServerMessage(conn, ws.OpText, []byte("refreshVersions"))
					if err != nil {
						if _, ok := err.(wsutil.ClosedError); ok {
							logger.Debug("Client disconnected")
						} else {
							logger.Error(err)
						}
						return
					}
					lastSend = count
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}
}
