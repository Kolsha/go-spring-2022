//go:build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	field := reflect.ValueOf(obj).Elem().FieldByName(name)
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}
