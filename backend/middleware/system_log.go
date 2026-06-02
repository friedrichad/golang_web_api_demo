package middleware

import (
	"bytes"
	"log"
	"strconv"

	"github.com/friedrichad/golang_web_api_demo/backend/model"
	"github.com/friedrichad/golang_web_api_demo/backend/rabbitmq"
	"github.com/gin-gonic/gin"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *ResponseBodyWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
func SystemLogMiddleware(rmq *rabbitmq.RabbitMQ) gin.HandlerFunc {

	return func(c *gin.Context) {
		writer := &ResponseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = writer

		c.Next()
		userID := 0

		if v, exists := c.Get("user_id"); exists {
			if idStr, ok := v.(string); ok {
				if id, err := strconv.Atoi(idStr); err == nil {
					userID = id
				}
			}
		}

		respBody := writer.body.String()

		if len(respBody) > 5000 {
			respBody = respBody[:5000]
		}

		log.Printf("[SYSTEM_LOG] Route: %s, Status: %d, ResponseBodyLen: %d, BodyPreview: %.100s",
			c.Request.URL.Path,
			c.Writer.Status(),
			len(respBody),
			respBody)

		logData := model.SystemLogMessage{
			UserID: userID,

			HTTPMethod: c.Request.Method,

			Route: c.Request.URL.Path,

			StatusInt: c.Writer.Status(),

			IPAddress: c.ClientIP(),

			ResponseBody: respBody,
		}

		if rmq != nil {
			go rmq.PublishSystemLog(logData)
		}
	}
}
