// This library is created as an experiment to see how a ServiceLocator
// type pattern holds up in a strictly typed environment. Dependecy Injection
// and the Factory Pattern are tools that are essential when rapidly building
// a system with changing requirements. Binding Named services to factories
// allows us to easily change the concrete implementation of a service and
// gives us the opportunity to compose and inject any other services registered
// to the system
package servicemanager

import "fmt"

// Read only version of service manager that is to be passed to factories
type ServiceLocator interface {
	Get(name string) (interface{}, error) // Retreives a defined service
	Has(name string) bool                 // Let you know if a service has been defined
}

// Quick and dirty anonymous factories
type ServiceFactoryCallback func(sl ServiceLocator) (interface{}, error)

// sometimes you may need a bit more than just a callback for your factories
type ServiceFactory interface {
	Create(sl ServiceLocator) (interface{}, error)
}

// main interface that allows registering and instantiating services
type ServiceManager interface {
	Set(name string, sl ServiceFactoryCallback) error // register a callback factory
	SetFactory(name string, sl ServiceFactory) error  // register an interface factory
	ServiceLocator
}

// concrete implementation of the service manager
type serviceManager struct {
	factories map[string]interface{} // registry of defined factories
	instances map[string]interface{} // registry of succesfully instantiated services
}

// Create an new ServiceManager
func New() ServiceManager {
	return &serviceManager{
		factories: map[string]interface{}{},
		instances: map[string]interface{}{},
	}
}

// Will create and return an instance of a named service and inject any
// tagged dependencies
func (s *serviceManager) Get(name string) (inst interface{}, err error) {
	inst = nil
	err = nil

	defer func() {
		if r := recover(); r != nil {
			inst = nil
			err = fmt.Errorf("%v", r)
		}
	}()

	if !s.Has(name) {
		return nil, fmt.Errorf("Service Not Found: %s", name)
	}

	inst, ok := s.instances[name]

	if ok {
		if inst == nil {
			return nil, fmt.Errorf("Circular Dependency Detected")
		}
		return inst, nil
	}

	// protect against circular dependencies
	s.instances[name] = nil

	factoryInst := s.factories[name]

	switch factory := factoryInst.(type) {
	case ServiceFactory:
		inst, err = s.getFactory(factory)
	case ServiceFactoryCallback:
		inst, err = s.getCallback(factory)
	default:
		return nil, fmt.Errorf("Invalid Factory Type")
	}

	if err != nil {
		delete(s.instances, name)
		return nil, err
	}

	s.instances[name] = inst

	inst, err = s.injectDependencies(inst)

	if err != nil {
		delete(s.instances, name)
		return nil, err
	}

	return inst, nil
}

func (s *serviceManager) getFactory(factory ServiceFactory) (interface{}, error) {
	return factory.Create(s)
}

func (s *serviceManager) getCallback(cb ServiceFactoryCallback) (interface{}, error) {
	return cb(s)
}

// Has() looks up if the requested service has been regisitered against the ServiceManager
func (s *serviceManager) Has(name string) bool {
	_, ok := s.factories[name]

	return ok
}

// SetFactory registers a ServiceFactory interface to a named service
func (s *serviceManager) SetFactory(name string, factory ServiceFactory) error {
	return s.set(name, factory)
}

// SetFactory registers a Callback method as a factory to a named service
// this provides a much lighter weight approach to defining factorys
func (s *serviceManager) Set(name string, cb ServiceFactoryCallback) error {
	return s.set(name, cb)
}

func (s *serviceManager) set(name string, constructor interface{}) error {
	if s.Has(name) {
		return fmt.Errorf("Service Already Set: %s", name)
	}

	s.factories[name] = constructor

	return nil
}

func (s *serviceManager) injectDependencies(inst interface{}) (interface{}, error) {
	dict := GetDependencies(inst)
	fieldValues := make(map[string]interface{})

	for fieldName, serviceName := range dict {
		service, err := s.Get(serviceName)

		if err != nil {
			return nil, err
		}

		fieldValues[fieldName] = service
	}

	return InjectDependencies(inst, fieldValues)
}
