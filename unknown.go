package unknown

import (
	"encoding/json"
	"reflect"
	"strings"
)

func UnmarshalJSON(jsonStr []byte, obj interface{}) (otherFields map[string]json.RawMessage, err error) {
	otherFields = make(map[string]json.RawMessage)
	objValue := reflect.ValueOf(obj).Elem()
	knownFields := map[string]reflect.Value{}
	for i := 0; i != objValue.NumField(); i++ {
		jsonName := strings.Split(objValue.Type().Field(i).Tag.Get("json"), ",")[0]
		if jsonName == "" {
			jsonName = objValue.Type().Field(i).Name
		}
		knownFields[strings.ToLower(jsonName)] = objValue.Field(i)
	}

	err = json.Unmarshal(jsonStr, &otherFields)
	if err != nil {
		return
	}

	for key, chunk := range otherFields {
		if field, found := knownFields[key]; found {
			err = json.Unmarshal(chunk, field.Addr().Interface())
			if err != nil {
				return
			}
			delete(otherFields, key)
		}
	}
	return
}

func MarshalJSON(obj interface{}, otherFields map[string]json.RawMessage) (bts []byte, err error) {
	objValue := reflect.ValueOf(obj)
	knownFields := map[string]reflect.Value{}
	for i := 0; i != objValue.NumField(); i++ {
		jsonName := strings.Split(objValue.Type().Field(i).Tag.Get("json"), ",")[0]
		if jsonName == "" {
			jsonName = objValue.Type().Field(i).Name
		}
		knownFields[jsonName] = objValue.Field(i)
	}

	for k := range knownFields {
		if field, found := knownFields[k]; found {
			b, err := json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			otherFields[k] = b
		}
	}
	return json.Marshal(otherFields)
}
