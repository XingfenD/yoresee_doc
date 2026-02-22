package api

import (
	"context"
	"net/http"
	"reflect"

	"github.com/XingfenD/yoresee_doc/internal/i18n"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct{}

type HandlerFunc func(ctx context.Context, req Request) (resp Response, err error)

func (h *BaseHandler) GinHandle(reqType reflect.Type, handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if userExternalID, exists := c.Get("user_external_id"); exists {
			ctx = context.WithValue(ctx, "user_external_id", userExternalID)
		}

		req := reflect.New(reqType).Interface().(Request)
		BindRequest(c, req)

		resp, err := handler(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GenBaseRespWithErrAndCtx(c, err))
			return
		}

		translateResponseMessage(c, resp)

		c.JSON(http.StatusOK, resp)
	}
}

type responseWithMessage interface {
	GetMessage() string
	SetMessage(string)
}

func translateResponseMessage(c *gin.Context, resp Response) {
	v := reflect.ValueOf(resp).Elem()

	messageField := v.FieldByName("Message")
	if messageField.IsValid() && messageField.CanSet() && messageField.Kind() == reflect.String {
		currentMsg := messageField.String()
		translatedMsg := i18n.Translate(c, currentMsg)
		messageField.SetString(translatedMsg)
	}

	baseField := v.FieldByName("BaseResponse")
	if baseField.IsValid() && baseField.Kind() == reflect.Struct {
		baseMessageField := baseField.FieldByName("Message")
		if baseMessageField.IsValid() && baseMessageField.CanSet() && baseMessageField.Kind() == reflect.String {
			currentMsg := baseMessageField.String()
			translatedMsg := i18n.Translate(c, currentMsg)
			baseMessageField.SetString(translatedMsg)
		}
	}
}
