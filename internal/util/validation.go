package util

import (
	"fmt"
	"reflect"
	"strings"
)

// ValidateStruct 验证结构体字段
//
//	param s interface{}
//	return error
//	author system
//	update 2025-01-19 12:00:00
func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 检查 required 标签
		if tag := fieldType.Tag.Get("binding"); strings.Contains(tag, "required") {
			if isZeroValue(field) {
				return fmt.Errorf("field %s is required", fieldType.Name)
			}
		}

		// 检查 oneof 标签
		if tag := fieldType.Tag.Get("binding"); strings.Contains(tag, "oneof=") {
			oneofPart := extractOneofValues(tag)
			if oneofPart != "" && !isValidOneofValue(field, oneofPart) {
				return fmt.Errorf("field %s must be one of: %s", fieldType.Name, oneofPart)
			}
		}
	}

	return nil
}

// isZeroValue 检查值是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		return false
	}
}

// extractOneofValues 提取 oneof 标签的值
func extractOneofValues(tag string) string {
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "oneof=") {
			return strings.TrimPrefix(part, "oneof=")
		}
	}
	return ""
}

// isValidOneofValue 检查值是否在 oneof 列表中
func isValidOneofValue(v reflect.Value, oneofValues string) bool {
	if v.Kind() != reflect.String {
		return true // 只验证字符串类型
	}

	value := v.String()
	validValues := strings.Split(oneofValues, " ")
	for _, validValue := range validValues {
		if value == validValue {
			return true
		}
	}
	return false
}