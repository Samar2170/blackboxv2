package design

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func run3(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				fmt.Println("Break the loop")
				return nil
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in a loop")

			}
		}
	})
	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				fmt.Println("Break the loop")
				return nil
			case <-time.After(1 * time.Second):
				fmt.Println("Hello in a loop")

			}
		}
	})
	err := g.Wait()
	if err != nil {
		fmt.Println("Error group: ", err)
	}
	fmt.Println("Main done")
}
