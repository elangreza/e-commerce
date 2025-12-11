package task

import (
	"context"
	"fmt"
	"time"
)

type (
	orderService interface {
		RemoveExpiryOrder(ctx context.Context, duration time.Duration) (int, error)
	}

	TaskOrder struct {
		closeChan            chan struct{}
		svc                  orderService
		removeExpiryDuration time.Duration
	}
)

func NewTaskOrder(orderService orderService) *TaskOrder {
	to := &TaskOrder{
		closeChan: make(chan struct{}),
		svc:       orderService,
	}

	go to.backgroundJobs()

	return to
}

func (to *TaskOrder) SetRemoveExpiryDuration(d time.Duration) {
	to.removeExpiryDuration = d
}

func (to *TaskOrder) removeExpiryOrder() error {
	expiryOrder, err := to.svc.RemoveExpiryOrder(context.Background(), to.removeExpiryDuration)
	if err != nil {
		return err
	}

	if expiryOrder > 0 {
		fmt.Printf("removing %d order(s)\n", expiryOrder)
	} else {
		// fmt.Println("no expiry order")
	}

	return nil
}

func (to *TaskOrder) Close() {
	to.closeChan <- struct{}{}
}

func (to *TaskOrder) backgroundJobs() {
	fmt.Println("running order backgroundJobs")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			if to.removeExpiryDuration > 0 {
				err := to.removeExpiryOrder()
				if err != nil {
					fmt.Println("getting error from RemoveExpiryOrder", err)
				}
			}

		case <-to.closeChan:
			fmt.Println("order task closed")
			return
		}
	}
}
