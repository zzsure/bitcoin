package response

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, code int, msg string, data interface{}) {
	requestId := c.MustGet("requestId")
	c.JSON(200, map[string]interface{}{
		"data":       data,
		"error_no":   code,
		"error_msg":  msg,
		"request_id": requestId,
	})
}

func ClientErr(c *gin.Context, msg string) {
	Response(c, 400, msg, nil)
}

func ServerErr(c *gin.Context, msg string) {
	Response(c, 500, msg, nil)
}

func Success(c *gin.Context) {
	Response(c, 0, "success", nil)
}

