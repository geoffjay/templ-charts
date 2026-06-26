package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// PropertyAccessor resolves a value of type V from a datum of type D. It is
// either a func(D) V or a string path (e.g. "id", "data.color") that is
// resolved via reflection over struct fields and map keys. Mirrors nivo's
// PropertyAccessor + get function.
type PropertyAccessor[D any, V any] any // func(D) V | string

// GetPropertyAccessor converts a PropertyAccessor into a func(D) V. A nil or
// empty accessor returns the zero V. A string path is resolved via reflection
// (struct field or map key chain, dot-separated). A func(D) V is returned as-is.
func GetPropertyAccessor[D any, V any](accessor PropertyAccessor[D, V]) func(D) V {
	switch a := any(accessor).(type) {
	case nil:
		return func(D) V { var z V; return z }
	case string:
		path := strings.TrimSpace(a)
		if path == "" {
			return func(D) V { var z V; return z }
		}
		parts := strings.Split(path, ".")
		return func(d D) V {
			v, ok := getPath(reflect.ValueOf(d), parts)
			if !ok {
				var z V
				return z
			}
			return coerceValue[V](v)
		}
	case func(D) V:
		return a
	default:
		// Allow func-like values via reflection (e.g. func(D) V wrapped as any).
		rv := reflect.ValueOf(a)
		if rv.Kind() == reflect.Func {
			return func(d D) V {
				out := rv.Call([]reflect.Value{reflect.ValueOf(d)})
				if len(out) == 0 {
					var z V
					return z
				}
				v, ok := out[0].Interface().(V)
				if !ok {
					var z V
					return z
				}
				return v
			}
		}
		return func(D) V { var z V; return z }
	}
}

// getPath walks a dot path over struct fields and map[string]any keys.
// Returns the resolved interface value and ok=false if any step is missing.
func getPath(rv reflect.Value, parts []string) (any, bool) {
	for _, p := range parts {
		for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Struct:
			f := rv.FieldByName(p)
			if !f.IsValid() {
				// case-insensitive fallback
				rt := rv.Type()
				for i := 0; i < rt.NumField(); i++ {
					if strings.EqualFold(rt.Field(i).Name, p) {
						f = rv.Field(i)
						break
					}
				}
				if !f.IsValid() {
					return nil, false
				}
			}
			rv = f
		case reflect.Map:
			mk := reflect.ValueOf(p)
			if !mk.Type().AssignableTo(rv.Type().Key()) {
				return nil, false
			}
			mv := rv.MapIndex(mk)
			if !mv.IsValid() {
				return nil, false
			}
			rv = mv
		default:
			return nil, false
		}
	}
	return rv.Interface(), true
}

// coerceValue coerces an arbitrary any into V. Supports the numeric/string
// conversions nivo relies on (number→float64/int, anything→string).
func coerceValue[V any](v any) V {
	var z V
	rv := reflect.ValueOf(&z).Elem()
	switch rv.Kind() {
	case reflect.String:
		return any(toString(v)).(V)
	case reflect.Float32, reflect.Float64:
		return any(toFloat64(v)).(V)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f := toFloat64(v)
		rv.SetInt(int64(f))
		return z
	case reflect.Bool:
		return any(toBool(v)).(V)
	default:
		if cv, ok := v.(V); ok {
			return cv
		}
		return z
	}
}

func toFloat64(v any) float64 {
	switch x := v.(type) {
	case nil:
		return 0
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case int32:
		return float64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0
		}
		return f
	case time.Time:
		return float64(x.UnixMilli())
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(rv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(rv.Uint())
		case reflect.Float32, reflect.Float64:
			return rv.Float()
		}
		return 0
	}
}

func toString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case time.Time:
		return x.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toBool(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		b, err := strconv.ParseBool(x)
		if err != nil {
			return x != ""
		}
		return b
	default:
		return true
	}
}

// GetLabelGenerator builds a label generator from an accessor. A nil/empty
// accessor returns a generator that yields the zero V (typically "").
// Equivalent to nivo's getLabelGenerator.
func GetLabelGenerator[D any, V any](label any, labelFrom string) func(D) V {
	if label != nil {
		if g, ok := label.(func(D) V); ok {
			return g
		}
		if s, ok := label.(string); ok && strings.TrimSpace(s) != "" {
			return GetPropertyAccessor[D, V](s)
		}
	}
	if labelFrom != "" {
		return GetPropertyAccessor[D, V](labelFrom)
	}
	return func(D) V { var z V; return z }
}
