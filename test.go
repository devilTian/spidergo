package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan int, 1000)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				fmt.Println(v)
				time.Sleep(time.Second)
			}
		}()
	}
	for i := 0; i < 1000; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	fmt.Println("All Finished")
}
