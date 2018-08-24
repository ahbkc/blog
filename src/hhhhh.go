package main

import (
	"core"
	"utils"
)

func main() {
	general := utils.NewGeneral()
	general.Run(core.NewRouter()) //start blob system
}
