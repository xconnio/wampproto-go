package util

import "fmt"

func AsInt64(i any) (int64, bool) {
	switch v := i.(type) {
	case int64:
		return v, true
	case uint64:
		return int64(v), true // #nosec
	case uint8:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int32:
		return int64(v), true
	case uint:
		return int64(v), true // #nosec
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case float64:
		return int64(v), true
	case float32:
		return int64(v), true
	}
	return 0, false
}

func AsFloat64(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int32:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	}
	return 0.0, false
}

func AsBool(i any) (bool, bool) {
	boolean, ok := i.(bool)
	return boolean, ok
}

func ToBool(i any) bool {
	boolean, _ := i.(bool)
	return boolean
}

func AsString(i any) (string, bool) {
	str, ok := i.(string)
	return str, ok
}

func ToString(i any) string {
	str, _ := i.(string)
	return str
}

func AnysToStrings(input []any) ([]string, error) {
	result := make([]string, 0, len(input))
	for _, item := range input {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("element %v is not a string", item)
		}
		result = append(result, str)
	}
	return result, nil
}
