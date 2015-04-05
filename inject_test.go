package servicemanager_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ebittleman/go-servicemanager"
)

type FakeDependency interface {
	Value() int
}

type fakeDep struct {
	value int
}

func (f *fakeDep) Value() int {
	return f.value
}

type FakeInjectService struct {
	MyField  FakeDependency `inject:"FakeDependency"`
	MyField2 FakeDependency `inject:"FakeDependency2"`
}

func Test_GetDependencies(t *testing.T) {

	service := &FakeInjectService{}

	tags := servicemanager.GetDependencies(service)

	assertEquals(2, len(tags), t)

	assertEquals("FakeDependency", tags["MyField"], t)
	assertEquals("FakeDependency2", tags["MyField2"], t)

	// t.Log(tags)
}

func Test_NonStructSerivce(t *testing.T) {
	service := map[string]interface{}{
		"log.writer": os.Stderr,
		"log.prefix": "example: ",
		"log.flags":  log.LstdFlags,
	}

	tags := servicemanager.GetDependencies(service)
	t.Log(tags)
}

func Test_InjectDependencies(t *testing.T) {
	service := &FakeInjectService{}
	fieldValues := map[string]interface{}{
		"MyField":  &fakeDep{1},
		"MyField2": &fakeDep{2},
	}

	inst, err := servicemanager.InjectDependencies(service, fieldValues)

	if err != nil {
		t.Fatal(err)
	}

	service, ok := inst.(*FakeInjectService)

	if !ok {
		t.Fatalf("Failed to return correct type: %s", reflect.TypeOf(inst))
	}

	if service.MyField.Value() != 1 {
		t.Fatalf("Expected: 1, Got: %d", service.MyField.Value())
	}

	if service.MyField2.Value() != 2 {
		t.Fatalf("Expected: 2, Got: %d", service.MyField2.Value())
	}
}

func assertEquals(expected interface{}, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Fatalf("Expected: %v, Got: %v", expected, actual)
	}
}
