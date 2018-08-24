package main

import (
	"utils"
	"core"
)

func main() {
	general := utils.NewGeneral()
	general.Run(core.NewRouter()) //start blob system
}