package servicemanager_test

import (
	"testing"

	"github.com/ebittleman/go-servicemanager"
)

const SERVICENAME string = "FakeService"

func Test_CreateServiceManager(t *testing.T) {
	shouldCreateManager(t)
}

type FakeService struct {
	Created int
}

type FakeFactory struct{}

func FakeCallback(sl servicemanager.ServiceLocator) (interface{}, error) {
	return &FakeService{1}, nil
}

func (f *FakeFactory) Create(sl servicemanager.ServiceLocator) (interface{}, error) {
	return FakeCallback(sl)
}

func Test_SetFactory(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetFactory(sm, SERVICENAME, t)

	shouldNotSetFactory(sm, SERVICENAME, t)
}

func Test_SetCallback(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetCallback(sm, SERVICENAME, t)

	shouldNotSetCallback(sm, SERVICENAME, t)
}

func Test_HasFactory(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetFactory(sm, SERVICENAME, t)

	shouldHave(sm, SERVICENAME, t)
}

func Test_HasCallback(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetCallback(sm, SERVICENAME, t)

	shouldHave(sm, SERVICENAME, t)
}

func Test_CreateFromFactory(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetFactory(sm, SERVICENAME, t)

	shouldCreateInstance(sm, SERVICENAME, t)
}

func Test_CreateFromCallback(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetCallback(sm, SERVICENAME, t)

	shouldCreateInstance(sm, SERVICENAME, t)
}

func Test_DoubleSetCallback(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetFactory(sm, SERVICENAME, t)

	shouldNotSetCallback(sm, SERVICENAME, t)
}

func Test_DoubleSetFactory(t *testing.T) {
	sm := shouldCreateManager(t)

	shouldSetCallback(sm, SERVICENAME, t)

	shouldNotSetFactory(sm, SERVICENAME, t)
}

func shouldCreateManager(t *testing.T) servicemanager.ServiceManager {
	sm := servicemanager.New()

	if sm == nil {
		t.Fatal("Could not create new service manager")
	}

	return sm
}

func shouldSetFactory(sm servicemanager.ServiceManager, name string, t *testing.T) {
	factory := &FakeFactory{}

	err := sm.SetFactory(name, factory)

	if err != nil {
		t.Fatal(err)
	}
}

func shouldNotSetFactory(sm servicemanager.ServiceManager, name string, t *testing.T) {
	factory := &FakeFactory{}

	err := sm.SetFactory(name, factory)

	if err == nil {
		t.Fatal("Should Not Have Been Able To Set Services")
	}
}

func shouldSetCallback(sm servicemanager.ServiceManager, name string, t *testing.T) {
	err := sm.Set(name, FakeCallback)

	if err != nil {
		t.Fatal(err)
	}
}

func shouldNotSetCallback(sm servicemanager.ServiceManager, name string, t *testing.T) {
	err := sm.Set(name, FakeCallback)

	if err == nil {
		t.Fatal("Should Not Have Been Able To Set Services")
	}
}

func shouldHave(sm servicemanager.ServiceManager, name string, t *testing.T) {
	has := sm.Has(name)

	if !has {
		t.Fatalf("Service Not Found: %s", name)
	}
}

func shouldCreateInstance(sm servicemanager.ServiceManager, name string, t *testing.T) {
	inst, err := sm.Get(name)

	if err != nil {
		t.Fatal(err)
	}

	_, ok := inst.(*FakeService)

	if !ok {
		t.Fatal("Cannot Assert Correct Type of Created Service")
	}
}
