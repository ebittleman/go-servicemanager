package example

// This is the

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	// sm "github.com/ebittleman/go-servicemanager"
	sm "gopkg.in/ebittleman/go-servicemanager.v0"
)

const LOG_SERVICE = "Log"
const CONFIG_SERVICE = "Config"

// Service interface intended to fine the sum of 2 integers
type Adder interface {
	Add(int, int) int
}

// Simple implementation of adder that acts as expected
type simpleAdder struct{}

func (s *simpleAdder) Add(a, b int) int {
	return a + b
}

// Unique implementation of adder that will divide the sum in half before
// returning it
type halfAdder struct{}

func (s *halfAdder) Add(a, b int) int {
	return int((a + b) / 2)
}

// A struct that does not implement Adder
type notAdder struct{}

// Calculator is a struct that requires both an Adder and an *log.Logger to be
// injected during its construction. In this example I have elected to inject
// the Adder in the factory and the Logger as a named service via the built-in
// dependecy injection in the service manager
type Calculator struct {
	// a simple dependency that will need to be injected from the factory
	Adder Adder

	// notice the inject tag, the service manager will automatically get an
	// instance of the named service "Log" and set this field to
	// that value. That factory will not have to inject this value.
	Log *log.Logger `inject:"Log"`
}

// Calculator itself is an implementation of Adder, but that is just a
// coincedence of this example. Its Add implementation will use the injected
// adder to do a calculation, output that answer to the Logger, then return
// the found sum
func (c *Calculator) Add(a, b int) int {
	ans := c.Adder.Add(a, b)
	c.Log.Printf("%d + %d = %d\n", a, b, ans)
	return ans
}

var factories map[string]sm.ServiceFactoryCallback

func init() {
	// initialized a map of named services to factories
	factories = map[string]sm.ServiceFactoryCallback{
		// Create a SimpleAdder Services
		"SimpleAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &simpleAdder{}, nil
		},

		// Create a HalfAdder Service
		"HalfAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &halfAdder{}, nil
		},

		// Create a service that does not implement Adder
		"NotAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &notAdder{}, nil
		},

		// Create a calculator that consumes SimpleAdder service
		"Calculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("SimpleAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		// Create a calculator that consumes a HalfAdder
		"HalfCalculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("HalfAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		// Create a calculator that trys to consume a service that does not
		// implment Adders. Hint, this should fail
		"NotCalculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("NotAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		// this is a basic example of how to load system configurations
		// typically I would suggest using the factory to parse some config file
		// or use system environment variables to hydrate and actual struct
		// defined to house the options of the system. But for this example a
		// simple map will suffice
		CONFIG_SERVICE: func(sl sm.ServiceLocator) (interface{}, error) {
			return map[string]interface{}{
				"log.writer": os.Stdout,
				"log.prefix": "example: ",
				"log.flags":  log.Lshortfile,
			}, nil
		},

		// Create a new *log.Logger and inject it with the system configuration
		LOG_SERVICE: func(sl sm.ServiceLocator) (interface{}, error) {
			config := GetConfig(sl)
			return log.New(
				config["log.writer"].(io.Writer),
				config["log.prefix"].(string),
				config["log.flags"].(int),
			), nil
		},
	}
}

func Example() {
	var (
		err     error             = nil
		halfy   *Calculator       = nil
		calc    *Calculator       = nil
		manager sm.ServiceManager = sm.New()
	)

	for name, factory := range factories {
		manager.Set(name, factory)
	}

	log := GetLog(manager)

	// Should get an error here
	if _, err := GetCalculator("NotCalculator", manager); err != nil {
		log.Println(err)
	}

	// should succefuly create a normal Calculator
	if calc, err = GetCalculator("Calculator", manager); err != nil {
		log.Fatalln(err)
	}

	two := calc.Add(1, 1)

	if two != 2 {
		log.Fatalln(fmt.Errorf("Expected: 2, Got: %d", two))
	}

	// should create a calculator that divids all sumations by 2
	if halfy, err = GetCalculator("HalfCalculator", manager); err != nil {
		log.Fatalln(err)
	}

	one := halfy.Add(1, 1)

	if one != 1 {
		log.Fatalln(fmt.Errorf("Expected: 1, Got: %d", one))
	}

	//output:
	// example: example_test.go:146: interface conversion: *example.notAdder is not example.Adder: missing method Add
	// example: example_test.go:62: 1 + 1 = 2
	// example: example_test.go:62: 1 + 1 = 1
}

func GetLog(sl sm.ServiceLocator) *log.Logger {
	inst, err := sl.Get(LOG_SERVICE)

	if err != nil {
		panic(err)
	}

	return inst.(*log.Logger)
}

func GetConfig(sl sm.ServiceLocator) map[string]interface{} {
	inst, err := sl.Get(CONFIG_SERVICE)

	if err != nil {
		panic(err)
	}

	return inst.(map[string]interface{})
}

func GetCalculator(name string, sl sm.ServiceLocator) (*Calculator, error) {
	var (
		err        error       = nil
		ok         bool        = false
		inst       interface{} = nil
		calculator *Calculator = nil
	)

	if inst, err = sl.Get(name); err != nil {
		return nil, err
	}

	if calculator, ok = inst.(*Calculator); !ok {
		return nil, fmt.Errorf("Expected: *Calculator, Got: %v", reflect.TypeOf(inst))
	}

	return calculator, nil
}
