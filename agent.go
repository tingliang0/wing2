package main

import (
	"encoding/json"
	"fmt"

	"github.com/antonholmquist/jason"
)

type Agent struct {
	Name     string
	Age      int
	handlers map[string]func(o *jason.Object) ([]byte, error)
	conn     *connection
}

func (a *Agent) Init(name string, age int) {
	a.Name = name
	a.Age = age

	a.handlers = make(map[string]func(o *jason.Object) ([]byte, error))
	a.handlers["hello"] = a.Hello
	a.handlers["add_age"] = a.AddAge
}

func (a *Agent) Hello(o *jason.Object) ([]byte, error) {
	return json.Marshal(a)
}

func (a *Agent) AddAge(o *jason.Object) ([]byte, error) {
	a.Age += 10
	return json.Marshal(a)
}

func (a *Agent) dispatcher(msg []byte) {
	o, err := jason.NewObjectFromBytes(msg)
	if err != nil {
		fmt.Printf("json -> object get err %s\n", err)
		return
	}
	name, _ := o.GetString("Name")
	handler := a.handlers[name]
	if handler == nil {
		fmt.Printf("unknown proto: %s\n", name)
		return
	}
	ret, err := handler(o)
	if err != nil {
		fmt.Printf("handle %s proto err %s\n", name, err)
		return
	}
	a.conn.send <- ret
}
