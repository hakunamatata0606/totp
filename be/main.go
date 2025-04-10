package main

import (
	"example/totp/appstate"
	"example/totp/router"
)

func main() {
	state := appstate.GetAppState()
	router := router.GetRouter()

	//todo: get key cert from appstate
	router.Run(state.Config.ServerUrl)
}
