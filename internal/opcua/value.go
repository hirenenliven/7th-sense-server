package opcua

import (
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
// text and date-times text; ok is false for a null variant or a type the collector does not
// store (arrays, structures, byte strings).
func ValueOf(variant *ua.Variant) (Value, bool) {
	if variant == nil {
		return Value{}, false
	}
	switch value := variant.Value().(type) {
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
