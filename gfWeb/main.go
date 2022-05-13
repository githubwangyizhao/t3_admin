package main

import (
	_ "gfWeb/boot"
	_ "gfWeb/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	s := g.Server()
	s.Run()
}
