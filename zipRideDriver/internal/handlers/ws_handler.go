package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type WSHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewWSHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *WSHandler {
	return &WSHandler{cfg: cfg, db: db, rdb: rdb, log: log}
}

func (h *WSHandler) DriverWS(c *gin.Context) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(uint)
	idStr := strings.TrimSpace(c.Param("driverId"))
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || uid == 0 || uint(id64) != uid {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	channel := "driver:events:" + strconv.Itoa(int(uid))
	ctx := c.Request.Context()
	pubsub := h.rdb.Subscribe(ctx, channel)
	defer pubsub.Close()
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			return
		}
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			return
		}
	}
}
