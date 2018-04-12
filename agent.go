package main

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"

	"github.com/antonholmquist/jason"
)

type Agent struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	conn     *connection
	handlers map[string]reflect.Value
}

type RespMsg struct {
	Errcode int         `json:"errcode"`
	Payload interface{} `json:"payload"`
}

func (a *Agent) Init(name string, age int) {
	a.Name = name
	a.Age = age

	a.handlers = make(map[string]reflect.Value)
	s := reflect.TypeOf(a)
	fm := reflect.ValueOf(a)
	for i := 0; i < s.NumMethod(); i++ {
		m := s.Method(i)
		if matched, _ := regexp.MatchString("Handle.*", m.Name); matched {
			name = strings.ToLower(m.Name[6:])
			a.handlers[name] = fm.MethodByName(m.Name)
		}
	}
}

func (a *Agent) Dispatcher(msg []byte) {
	o, err := jason.NewObjectFromBytes(msg)
	if err != nil {
		Error.Printf("json -> object err %s\n", err)
		return
	}
	name, err := o.GetString("Name")
	if err != nil {
		Error.Println("invail proto: no name field")
		return
	}
	handler := a.handlers[name]
	if !handler.IsValid() {
		Error.Printf("unknown proto: %s\n", name)
		return
	}
	inputs := make([]reflect.Value, 1)
	inputs[0] = reflect.ValueOf(o)

	// call
	rets := handler.Call(inputs)
	if len(rets) == 0 {
		return
	}

	// error
	if rets[1].Interface() != nil {
		Error.Printf("handle %s proto err %s\n", name, err)
		return
	}

	obj, err := json.Marshal(rets[0].Interface())
	if err != nil {
		Error.Printf("object -> json err %s\n", err)
		return
	}
	a.conn.send <- obj
}

// Command
func (a *Agent) HandleHello(o *jason.Object) (RespMsg, error) {
	return RespMsg{0, a}, nil
}

func (a *Agent) HandleAddAge(o *jason.Object) (RespMsg, error) {
	a.Age += 10
	return RespMsg{Errcode: 0, Payload: a}, nil
}
