package servicemanager

import (
	"fmt"
	"reflect"
)

const INJECT_TAG string = "inject"

func GetDependencies(inst interface{}) map[string]string {
	val := reflect.ValueOf(inst).Elem()
	typeOf := val.Type()

	numFields := val.NumField()

	dict := map[string]string{}

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

	val := reflect.ValueOf(inst).Elem()

	for name, fieldValue := range fieldValues {

		if fieldValue == nil {
			continue
		}

		val.FieldByName(name).Set(reflect.ValueOf(fieldValue))
	}

	return inst, nil
}
