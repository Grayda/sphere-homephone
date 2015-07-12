package main

import (
	"fmt"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	ledmodel "github.com/ninjasphere/sphere-go-led-controller/model"
)

type something struct {
	Derp string
}

var led *ninja.ServiceClient

func showPane() {

	led := driver.conn.GetServiceClient("$node/" + config.Serial() + "/led-controller")
	fmt.Println(led)
	err := led.Call("displayIcon", ledmodel.IconRequest{
		Icon: "orphaned.gif",
	}, nil, 0)

	if err != nil {
		fmt.Println("========")
		fmt.Println(err)
		fmt.Println("========")
	}
}
