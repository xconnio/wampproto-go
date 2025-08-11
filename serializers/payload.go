package serializers

import (
	"encoding/json"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/vmihailenco/msgpack/v5"
)

func decode(arr []any) ([]any, map[string]any, error) {
	if len(arr) == 0 {
		return nil, nil, nil
	}

	if len(arr) > 2 {
		return nil, nil, fmt.Errorf("too many args to decode")
	}

	args, ok := arr[0].([]any)
	if !ok {
		return nil, nil, fmt.Errorf("args element is not []any")
	}

	var kwargs map[string]any
	if len(arr) == 2 {
		kwargs, ok = arr[1].(map[string]any)
		if !ok {
			return nil, nil, fmt.Errorf("kwargs element is not map[string]any")
		}
	}

	return args, kwargs, nil
}

func prepareForEncode(args []any, kwargs map[string]any) []any {
	var data []any
	if len(args) != 0 {
		data = append(data, args)
	}

	if len(kwargs) != 0 {
		if len(args) == 0 {
			data = append(data, []any{})
		}

		data = append(data, kwargs)
	}

	if len(data) == 0 {
		return nil
	}

	return data
}

func CBOREncodePayload(args []any, kwargs map[string]any) ([]byte, error) {
	data := prepareForEncode(args, kwargs)
	if len(data) == 0 {
		return nil, nil
	}

	return cbor.Marshal(data)
}

func CBORDecodePayload(b []byte) ([]any, map[string]any, error) {
	var arr []any
	if err := cbor.Unmarshal(b, &arr); err != nil {
		return nil, nil, err
	}

	return decode(arr)
}

func MsgPackEncodePayload(args []any, kwargs map[string]any) ([]byte, error) {
	data := prepareForEncode(args, kwargs)
	if len(data) == 0 {
		return nil, nil
	}

	return msgpack.Marshal(data)
}

func MsgPackDecodePayload(b []byte) ([]any, map[string]any, error) {
	var arr []any
	if err := msgpack.Unmarshal(b, &arr); err != nil {
		return nil, nil, err
	}

	return decode(arr)
}

func JSONEncodePayload(args []any, kwargs map[string]any) ([]byte, error) {
	data := prepareForEncode(args, kwargs)
	if len(data) == 0 {
		return nil, nil
	}

	return json.Marshal(data)
}

func JSONDecodePayload(b []byte) ([]any, map[string]any, error) {
	var arr []any
	if err := json.Unmarshal(b, &arr); err != nil {
		return nil, nil, err
	}

	return decode(arr)
}

func DeserializePayload(serializerID uint64, payload []byte) ([]any, map[string]any, error) {
	switch serializerID {
	case NoneSerializerID:
		return []any{payload}, make(map[string]any), nil
	case JSONSerializerID:
		return JSONDecodePayload(payload)
	case CBORSerializerID:
		return CBORDecodePayload(payload)
	case MsgPackSerializerID:
		return MsgPackDecodePayload(payload)
	default:
		return nil, nil, fmt.Errorf("serializer %d not recognized", serializerID)
	}
}

func SerializePayload(serializerID uint64, args []any, kwargs map[string]any) ([]byte, error) {
	switch serializerID {
	case NoneSerializerID:
		if len(args) == 0 && len(kwargs) == 0 {
			return nil, nil
		}

		if len(args) != 1 || len(kwargs) != 0 {
			return nil, fmt.Errorf("serializer %d requires exactly 1 arg", serializerID)
		}

		payload, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("serializer %d requires []byte", serializerID)
		}

		return payload, nil
	case JSONSerializerID:
		return JSONEncodePayload(args, kwargs)
	case CBORSerializerID:
		return CBOREncodePayload(args, kwargs)
	case MsgPackSerializerID:
		return MsgPackEncodePayload(args, kwargs)
	default:
		return nil, fmt.Errorf("serializer %d not recognized", serializerID)
	}
}
