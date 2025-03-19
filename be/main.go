package main

import (
	"example/totp/appstate"
	"example/totp/router"
)

func main() {
	state := appstate.GetAppState()
	router := router.GetRouter()

	router.Run(state.Config.ServerUrl)
}
