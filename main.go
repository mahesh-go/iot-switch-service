package main

import (
	"fmt"
	"net/http"

	"go.uber.org/dig"
)

type Switch struct {
	Id      string      `json:"id"`
	State   string      `json:"state"`
	Voltage int         `json:"voltage"`
	Wattage int         `json:"wattage"`
	Type    interface{} `json:"type,omitempty"`
}

type SwitchService struct {
	SwitchType *SwitchType
}

type SwitchType struct {
	Type string
}

func main() {
	//Dependency Injection Framework from uber/dig
	container := dig.New()
	container.Provide(func() *SwitchType {
		return &SwitchType{Type: "philips-hue-gen1"}
	})
	container.Provide(func(dep *SwitchType) *SwitchService {
		return &SwitchService{SwitchType: dep}
	})

	var iotSwitchService *SwitchService
	container.Invoke(func(svc *SwitchService) {
		iotSwitchService = svc
	})

	found := InitSwitchConfig(iotSwitchService.SwitchType.Type)

	if !found {
		fmt.Println("Switch Model not found in the registry")
		return
	}

	fmt.Println("IoT Switch Service API Server Ready !!")
	r := NewRouter()
	http.ListenAndServe(":8000", r)
}
