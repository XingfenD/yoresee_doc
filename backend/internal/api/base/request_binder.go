package api_base

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

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

		setFieldValue(field, paramValue)
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

		// Handle the case where form tag has options like "form:\"name,omitempty\""
		tagParts := strings.Split(formTag, ",")
		fieldName := tagParts[0]

		// Special handling for complex types like structs, slices, etc.
		if fieldName != "" {
			bindComplexField(c, field, fieldName)
		}
	}
}

// bindComplexField handles binding of complex field types including structs, slices, and basic types
func bindComplexField(c *gin.Context, field reflect.Value, fieldName string) {
	if !field.CanSet() {
		return
	}

	// Handle slice/array types
	if field.Kind() == reflect.Slice || field.Kind() == reflect.Array {
		values := c.QueryArray(fieldName)
		if len(values) > 0 {
			// Create a new slice of the appropriate type
			sliceType := field.Type()
			elemType := sliceType.Elem()

			newSlice := reflect.MakeSlice(sliceType, len(values), len(values))
			for i, val := range values {
				elemVal := reflect.New(elemType).Elem()
				setFieldValue(elemVal, val)
				newSlice.Index(i).Set(elemVal)
			}
			field.Set(newSlice)
		}
		return
	}

	// Handle pointer types
	if field.Kind() == reflect.Ptr {
		// For pointers, we need to create the underlying value if it doesn't exist
		if field.IsNil() {
			newVal := reflect.New(field.Type().Elem())
			field.Set(newVal)
		}

		// If the pointed-to type is a struct, try to parse it as JSON from query param
		pointedType := field.Type().Elem()
		if pointedType.Kind() == reflect.Struct {
			jsonStr := c.Query(fieldName)
			if jsonStr != "" {
				// Try to unmarshal the JSON string into the struct
				instance := reflect.New(pointedType).Interface()
				if err := json.Unmarshal([]byte(jsonStr), instance); err == nil {
					field.Elem().Set(reflect.ValueOf(instance).Elem())
				} else {
					// If JSON parsing fails, try to bind individual fields
					bindStructFromQuery(c, field.Elem(), fieldName)
				}
			} else {
				// Try to bind struct fields individually based on nested field names
				bindStructFromQuery(c, field.Elem(), fieldName)
			}
		} else {
			// Handle pointer to basic types
			value := c.Query(fieldName)
			if value != "" {
				elemVal := reflect.New(pointedType).Elem()
				setFieldValue(elemVal, value)
				field.Elem().Set(elemVal)
			}
		}
		return
	}

	// Handle struct types
	if field.Kind() == reflect.Struct {
		bindStructFromQuery(c, field, fieldName)
		return
	}

	// Handle basic types
	value := c.Query(fieldName)
	if value != "" {
		setFieldValue(field, value)
	}
}

// bindStructFromQuery binds a struct from query parameters with field prefix
func bindStructFromQuery(c *gin.Context, structVal reflect.Value, prefix string) {
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			// Fallback to field name
			jsonTag = fieldType.Name
		} else {
			// Remove omitempty and other options
			parts := strings.Split(jsonTag, ",")
			jsonTag = parts[0]
		}

		fieldName := prefix + "." + jsonTag
		bindComplexField(c, field, fieldName)
	}
}

// setFieldValue sets the value of a field based on string input, handling different types
func setFieldValue(field reflect.Value, value string) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(value, 10, 64); err == nil {
			field.SetUint(u)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(value); err == nil {
			field.SetBool(b)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			field.SetFloat(f)
		}
	case reflect.Ptr:
		// For nil pointers, create a new instance and set the value
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		setFieldValue(field.Elem(), value)
	}
}

func BindRequest(c *gin.Context, req interface{}) {
	BindUriParams(c, req)

	BindQueryParams(c, req)

	c.ShouldBindJSON(req)
}
