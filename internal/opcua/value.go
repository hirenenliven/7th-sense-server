package opcua

import (
	"reflect"
	"time"

	"github.com/gopcua/opcua/ua"
)

// Value is a node value as the collector stores it: a number (booleans are 1 / 0) or text.
type Value struct {
	Number float64
	Text   string
	IsText bool
}

// ValueOf converts an OPC UA variant. Every numeric type becomes a float64, a boolean 1 or 0,
// text and date-times text. A one-element array is its element, because Siemens addresses such
// as DB9004.DBD8:REAL[1] read as a Float array of length 1. ok is false for a null variant or a
// type the collector does not store (longer arrays, structures, byte strings).
func ValueOf(variant *ua.Variant) (Value, bool) {
	if variant == nil {
		return Value{}, false
	}
	value := variant.Value()
	if variant.Has(ua.VariantArrayValues) {
		array := reflect.ValueOf(value)
		if array.Kind() != reflect.Slice || array.Len() != 1 {
			return Value{}, false
		}
		value = array.Index(0).Interface()
	}
	return scalarValueOf(value)
}

// scalarValueOf converts one decoded variant value (not an array).
func scalarValueOf(value any) (Value, bool) {
	switch value := value.(type) {
	case bool:
		if value {
			return Value{Number: 1}, true
		}
		return Value{Number: 0}, true
	case int8:
		return Value{Number: float64(value)}, true
	case uint8:
		return Value{Number: float64(value)}, true
	case int16:
		return Value{Number: float64(value)}, true
	case uint16:
		return Value{Number: float64(value)}, true
	case int32:
		return Value{Number: float64(value)}, true
	case uint32:
		return Value{Number: float64(value)}, true
	case int64:
		return Value{Number: float64(value)}, true
	case uint64:
		return Value{Number: float64(value)}, true
	case float32:
		return Value{Number: float64(value)}, true
	case float64:
		return Value{Number: value}, true
	case string:
		return Value{Text: value, IsText: true}, true
	case *ua.LocalizedText:
		if value == nil {
			return Value{}, false
		}
		return Value{Text: value.Text, IsText: true}, true
	case time.Time:
		return Value{Text: value.UTC().Format(time.RFC3339Nano), IsText: true}, true
	}
	return Value{}, false
}
