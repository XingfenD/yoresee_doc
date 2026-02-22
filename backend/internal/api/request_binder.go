package api

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func BindUriParams(c *gin.Context, req interface{}) {
	v := reflect.ValueOf(req).Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		uriTag := fieldType.Tag.Get("uri")
		if uriTag == "" {
			continue
		}

		paramValue := c.Param(uriTag)
		if paramValue == "" {
			continue
		}

		if field.Kind() == reflect.String {
			field.SetString(paramValue)
		}
	}
}

func BindQueryParams(c *gin.Context, req interface{}) {
	v := reflect.ValueOf(req).Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := typ.Field(i)

		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			continue
		}

		paramValue := c.Query(formTag)
		if paramValue == "" {
			continue
		}

		if field.Kind() == reflect.String {
			field.SetString(paramValue)
		}
	}
}

func BindRequest(c *gin.Context, req interface{}) {
	BindUriParams(c, req)

	BindQueryParams(c, req)

	c.ShouldBindJSON(req)
}
