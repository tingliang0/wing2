package main

import (
	"encoding/json"
	"fmt"
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

func (a *Agent) dispatcher(msg []byte) {
	o, err := jason.NewObjectFromBytes(msg)
	if err != nil {
		fmt.Printf("json -> object get err %s\n", err)
		return
	}
	name, _ := o.GetString("Name")
	handler := a.handlers[name]
	if !handler.IsValid() {
		fmt.Printf("unknown proto: %s\n", name)
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
		fmt.Printf("handle %s proto err %s\n", name, err)
		return
	}

	obj, err := json.Marshal(rets[0].Interface())
	if err != nil {
		fmt.Printf("marshal proto err %s\n", err)
		return
	}
	a.conn.send <- obj
}

// Command
func (a Agent) HandleHello(o *jason.Object) (RespMsg, error) {
	return RespMsg{0, a}, nil
}

func (a Agent) HandleAddAge(o *jason.Object) (RespMsg, error) {
	a.Age += 10
	return RespMsg{Errcode: 0, Payload: a}, nil
}
