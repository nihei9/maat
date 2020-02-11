package service

import (
	"fmt"

	"github.com/nihei9/maat/service/value"
	"github.com/pkg/errors"
)

func unmarshalValue(src interface{}) (value.Value, error) {
	srcMap, ok := src.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value.Value object must be a map[string]interface{}")
	}
	srcValueField, ok := srcMap["value"]
	if !ok {
		return nil, fmt.Errorf("value.Value object must contain a 'value' field")
	}
	switch srcValue := srcValueField.(type) {
	case string:
		return value.NewTextValue(srcValue), nil
	case []interface{}:
		v := value.NewListValue().(*value.ListValue)
		for _, srcElement := range srcValue {
			eValue, err := unmarshalValue(srcElement)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal the object expected to be the value.Value type")
			}
			v.Append(eValue)
		}
		return v, nil
	}

	return nil, fmt.Errorf("'value' field must be a string or a []interface{}")
}
