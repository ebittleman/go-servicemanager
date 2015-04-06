package servicemanager_test

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/ebittleman/go-servicemanager.v0"
)

const SERVICENAME string = "FakeService"

// First define yourself some interfaces
type MyDependency interface {
	Value() int
}

type MyInjection interface {
	OtherValue() int
}

type myDependency struct {
	Val int
}

func (m *myDependency) Value() int { return m.Val }

type myInjection struct {
	OtherVal int
}

func (m *myInjection) OtherValue() int { return m.OtherVal }

func Example() {

	// Start by taking a peek at some trivial external services our
	// new service will be consuming
	//
	// 	type MyDependency interface {
	// 		Value() int
	//	}
	//
	// 	type MyInjection interface {
	// 		OtherValue() int
	// 	}
	//
	// 	type myDependency struct {
	// 		Val int
	// 	}
	// 	func (m *myDependency) Value() int { return m.Val }
	//
	// 	type myInjection struct {
	// 		OtherVal int
	// 	}
	// 	func (m *myInjection) OtherValue() int { return m.OtherVal }

	// Define the struct of your new service
	type MyService struct {
		// a simple dependency that will need to be injected from the factory
		Dep MyDependency

		// notice the inject tag, the service manager will automatically get an
		// instance of the named service "NamedInjection" and set this field to
		// that value. That factory will not have to inject this value.
		Inject MyInjection `inject:"NamedInjection"`
	}

	// implement functionality (left out for now)

	// get a new instance of a ServiceManager
	sm := servicemanager.New()

	// Register a dependency, Note, there are no naming conventions or rules
	// name services anything you want
	sm.Set("NamedDependency", func(sl servicemanager.ServiceLocator) (interface{}, error) {
		return &myDependency{Val: 1}, nil
	})

	// Register the service that will be injected
	sm.Set("NamedInjection", func(sl servicemanager.ServiceLocator) (interface{}, error) {
		return &myInjection{OtherVal: 5}, nil
	})

	// register the new service just created, see how it leverages the
	// ServiceLocator to get the system's concrete implementation of
	// "NamedDependency" since we are using a static interface bound to some
	// factory, we can change what the "NamedDependency" factory returns with
	// out having to update any other code in our system
	sm.Set("MyNamedService", func(sl servicemanager.ServiceLocator) (interface{}, error) {
		inst, err := sl.Get("NamedDependency")
		return &MyService{Dep: inst.(MyDependency)}, err
	})

	// declare some variables for error checking and type assertions
	var (
		inst    interface{}
		err     error
		service *MyService
		ok      bool
	)

	// Check to see if there were any errors instantiating the service
	if inst, err = sm.Get("MyNamedService"); err != nil {
		fmt.Println(err)
		return
	}

	// assert the instance is of the expected type and assign it to a variable
	if service, ok = inst.(*MyService); !ok {
		fmt.Println(fmt.Errorf("Expected: *MyService, Got: %v", reflect.TypeOf(inst)))
		return
	}

	fmt.Println(service.Inject.OtherValue())
	fmt.Println(service.Dep.Value())
	// output:
	// 5
	// 1
}

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
