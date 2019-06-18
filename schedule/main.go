package main

import (
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("go,%d,%p", i, &i)
			//for 阻塞调度,不会出让cpu
			for {
			}
		}()
	}
	wg.Wait()
}
