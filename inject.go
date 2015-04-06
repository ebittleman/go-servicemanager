package servicemanager

import (
	"fmt"
	"reflect"
)

const INJECT_TAG string = "inject"

// Creates a map where the keys are the field names of the passed struct
// and the values are named services to be injected into those fields
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

// Takes a service instance and injects the map of field names as keys and
// instantiated services as values. The service locator is specifically left
// out of this method call so that users can do injections with their own systems
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
