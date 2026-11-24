package adafruitio

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAccessors checks every generated Get method: it must be safe to call on
// a nil struct, return the zero value for a nil field, and return the field's
// value when it is set.
func TestAccessors(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeOf(&Chart{}),
		reflect.TypeOf(&ChartParameters{}),
		reflect.TypeOf(&Data{}),
		reflect.TypeOf(&Feed{}),
		reflect.TypeOf(&Group{}),
		reflect.TypeOf(&Owner{}),
		reflect.TypeOf(&Throttle{}),
		reflect.TypeOf(&User{}),
	} {
		structType := typ.Elem()

		for i := 0; i < typ.NumMethod(); i++ {
			method := typ.Method(i)
			if !strings.HasPrefix(method.Name, "Get") {
				continue
			}
			name := structType.Name() + "." + method.Name
			field, ok := structType.FieldByName(strings.TrimPrefix(method.Name, "Get"))
			if !assert.True(t, ok, "%s has no matching field", name) {
				continue
			}

			// nil struct and nil field both return the zero value
			nilStruct := reflect.Zero(typ)
			emptyStruct := reflect.New(structType)
			for _, recv := range []reflect.Value{nilStruct, emptyStruct} {
				got := method.Func.Call([]reflect.Value{recv})[0]
				assert.True(t, got.IsZero(), "%s should return the zero value", name)
			}

			// a set field is returned: struct pointer getters return the
			// pointer itself, all others return the value it points to
			returnsPointer := method.Type.Out(0).Kind() == reflect.Ptr
			set := reflect.New(structType)
			want := reflect.New(field.Type.Elem())
			if !returnsPointer {
				want.Elem().Set(nonZero(field.Type.Elem()))
			}
			set.Elem().FieldByIndex(field.Index).Set(want)

			got := method.Func.Call([]reflect.Value{set})[0]
			if returnsPointer {
				assert.Equal(t, want.Interface(), got.Interface(), name)
			} else {
				assert.Equal(t, want.Elem().Interface(), got.Interface(), name)
			}
		}
	}
}

// nonZero returns a non-zero value of a type that a getter returns by value.
func nonZero(typ reflect.Type) reflect.Value {
	if typ == reflect.TypeOf(Timestamp{}) {
		return reflect.ValueOf(Timestamp{time.Date(2026, 11, 24, 16, 5, 19, 0, time.UTC)})
	}

	v := reflect.New(typ).Elem()
	switch typ.Kind() {
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int64:
		v.SetInt(42)
	case reflect.Float64:
		v.SetFloat(4.2)
	case reflect.String:
		v.SetString("value")
	default:
		panic("nonZero: unsupported type " + typ.String())
	}
	return v
}
