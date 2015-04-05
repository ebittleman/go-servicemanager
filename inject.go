package servicemanager

import (
	"fmt"
	"reflect"
)

const INJECT_TAG string = "inject"

func GetDependencies(inst interface{}) map[string]string {
	dict := map[string]string{}

	refVal := reflect.ValueOf(inst)

	if refVal.Kind() != reflect.Ptr {
		return dict
	}

	val := refVal.Elem()
	typeOf := val.Type()

	numFields := val.NumField()

	for i := 0; i < numFields; i++ {
		tag := typeOf.Field(i).Tag
		name := tag.Get(INJECT_TAG)

		if name == "" {
			continue
		}

		dict[typeOf.Field(i).Name] = name
	}

	return dict
}

func InjectDependencies(inst interface{}, fieldValues map[string]interface{}) (instance interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			instance = nil
			err = fmt.Errorf("%v", r)
		}
	}()

	refVal := reflect.ValueOf(inst)

	if refVal.Kind() != reflect.Ptr {
		return inst, nil
	}

	val := refVal.Elem()

	for name, fieldValue := range fieldValues {

		if fieldValue == nil {
			continue
		}

		val.FieldByName(name).Set(reflect.ValueOf(fieldValue))
	}

	return inst, nil
}
