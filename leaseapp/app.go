// Package main provides ...
package main

import (
	"github.com/fdgo/leaseapp/routers"
)

func main() {
	route := router.InitRouter()
	route.Run(":9000")
}
