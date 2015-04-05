package main

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

type Adder interface {
	Add(int, int) int
}

type simpleAdder struct{}

func (s *simpleAdder) Add(a, b int) int {
	return a + b
}

type halfAdder struct{}

func (s *halfAdder) Add(a, b int) int {
	return int((a + b) / 2)
}

type notAdder struct{}

type Calculator struct {
	Adder Adder
	Log   *log.Logger `inject:"Log"`
}

func (c *Calculator) Add(a, b int) int {
	ans := c.Adder.Add(a, b)
	c.Log.Printf("%d + %d = %d\n", a, b, ans)
	return ans
}

var factories map[string]sm.ServiceFactoryCallback

func init() {
	factories = map[string]sm.ServiceFactoryCallback{
		"SimpleAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &simpleAdder{}, nil
		},

		"HalfAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &halfAdder{}, nil
		},

		"NotAdder": func(sl sm.ServiceLocator) (interface{}, error) {
			return &notAdder{}, nil
		},

		"Calculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("SimpleAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		"HalfCalculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("HalfAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		"NotCalculator": func(sl sm.ServiceLocator) (interface{}, error) {
			adder, _ := sl.Get("NotAdder")
			return &Calculator{Adder: adder.(Adder)}, nil
		},

		CONFIG_SERVICE: func(sl sm.ServiceLocator) (interface{}, error) {
			return map[string]interface{}{
				"log.writer": os.Stderr,
				"log.prefix": "example: ",
				"log.flags":  log.LstdFlags,
			}, nil
		},

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

func main() {
	manager := sm.New()

	for name, factory := range factories {
		manager.Set(name, factory)
	}

	log := GetLog(manager)

	_, err := GetCalculator("NotCalculator", manager)

	if err != nil {
		log.Println(err)
	}

	calc, err := GetCalculator("Calculator", manager)

	if err != nil {
		log.Fatalln(err)
	}

	halfy, err := GetCalculator("HalfCalculator", manager)

	if err != nil {
		log.Fatalln(err)
	}

	one := halfy.Add(1, 1)

	if one != 1 {
		log.Fatalln(fmt.Errorf("Expected: 1, Got: %d", one))
	}

	two := calc.Add(1, 1)

	if two != 2 {
		log.Fatalln(fmt.Errorf("Expected: 2, Got: %d", two))
	}

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
	inst, err := sl.Get(name)

	if err != nil {
		return nil, err
	}

	calculator, ok := inst.(*Calculator)

	if !ok {
		return nil, fmt.Errorf("Expected: *Calculator, Got: %v", reflect.TypeOf(inst))
	}

	return calculator, nil
}
