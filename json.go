package av

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func ToInterface(av types.AttributeValue) interface{} {
	if _, ok := reflect.ValueOf(av).Elem().Type().FieldByName("Value"); ok {
		v := reflect.ValueOf(av).Elem().FieldByName("Value")
		if v.Kind() == reflect.Map {
			m := make(map[string]interface{})
			iter := v.MapRange()
			for iter.Next() {
				mVal := iter.Value().Interface()
				if mAv, ok2 := mVal.(types.AttributeValue); ok2 {
					m[iter.Key().Interface().(string)] = ToInterface(mAv)
				} else {
					m[iter.Key().Interface().(string)] = mVal
				}
			}
			return m
		} else {
			return v.Interface()
		}
	}
	return nil
}

func ToJSON(av types.AttributeValue) ([]byte, error) {
	val := ToInterface(av)

	j, err := json.MarshalIndent(&val, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("av to json conversion err: %w", err)
	}

	return j, nil
}

func FromJSON(js []byte, formatGuidePattern interface{}) (types.AttributeValue, error) {
	var t interface{}

	if formatGuidePattern != nil {
		if reflect.TypeOf(formatGuidePattern).Kind() != reflect.Ptr {
			return nil, errors.New("formatGuidePattern supplied must be a pointer")
		}

		err := json.Unmarshal(js, formatGuidePattern)
		if err != nil {
			return nil, fmt.Errorf("json to av prepare error: %s", err)
		}

		av, err := attributevalue.Marshal(formatGuidePattern)
		if err != nil {
			return nil, fmt.Errorf("json to av conversion error: %w", err)
		}

		return av, nil
	}

	err := json.Unmarshal(js, &t)
	if err != nil {
		return nil, fmt.Errorf("json to av unmarshal error: %w", err)
	}

	av, err := attributevalue.Marshal(t)

	return av, err
}
