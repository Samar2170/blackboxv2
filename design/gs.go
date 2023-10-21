package design

import (
	"context"
	"fmt"
	"log"
	"time"
)

func run(ctx context.Context) {
	wait := make(chan struct{}, 1)

	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		for {
			select {
			case <-ctx.Done():
				log.Println("Break the loop")
				break
			case <-time.After(1 * time.Second):
				fmt.Println("hellp in a loop")
			}
		}
	}()

	go func() {
		defer func() {
			wait <- struct{}{}
		}()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("break the loop")
				break
			case <-time.After(1 * time.Second):
				fmt.Println("Ciao in a loop")
			}
		}
	}()

	<-wait
	<-wait

	fmt.Println("Main done")
}
