package main

import (
	"sync"

	"github.com/wigfri/mustore/app"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	go func() {
		defer wg.Done()
		app.NewHttpServer().Start()
	}()

	wg.Wait()
}
