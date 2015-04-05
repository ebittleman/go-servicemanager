package main

import (
	"fmt"
	"reflect"

	sm "gopkg.in/ebittleman/go-servicemanager.v0"
)

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
	Adder Adder `inject:"HalfAdder"`
}

func (c *Calculator) Add(a, b int) int {
	return c.Adder.Add(a, b)
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
			return &Calculator{}, nil
		},
	}
}

func main() {
	manager := sm.New()

	for name, factory := range factories {
		manager.Set(name, factory)
	}

	inst, err := manager.Get("Calculator")

	if err != nil {
		panic(err)
	}

	calculator, ok := inst.(*Calculator)

	if !ok {
		panic(fmt.Errorf("Expected: *Calculator, Got: %v", reflect.TypeOf(inst)))
	}

	one := calculator.Add(1, 1)

	if one != 1 {
		panic(fmt.Errorf("Expected: 1, Got: %d", one))
	}

	fmt.Println("IT WORKED!! ", one)
}
