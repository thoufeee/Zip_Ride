package middleware

import "github.com/gin-gonic/gin"

// getting user id from token

func GetUserID(c *gin.Context) uint {
	if user_id, err := c.Get("user_id"); err {
		return user_id.(uint)
	}

	return 0
}
