package task

import (
	"context"
	"fmt"
	"time"
)

type (
	paymentService interface {
		RemoveExpiryPayment(ctx context.Context, duration time.Duration) (int, error)
	}

	TaskPayment struct {
		closeChan            chan struct{}
		svc                  paymentService
		removeExpiryDuration time.Duration
	}
)

func NewTaskPayment(paymentService paymentService, duration time.Duration) *TaskPayment {
	to := &TaskPayment{
		closeChan:            make(chan struct{}),
		svc:                  paymentService,
		removeExpiryDuration: duration,
	}

	go to.backgroundJobs()

	return to
}

func (to *TaskPayment) removeExpiryPayment() error {
	expiryPayment, err := to.svc.RemoveExpiryPayment(context.Background(), to.removeExpiryDuration)
	if err != nil {
		return err
	}

	if expiryPayment > 0 {
		fmt.Printf("removing %d payment(s)\n", expiryPayment)
	}

	return nil
}

func (to *TaskPayment) Close() {
	to.closeChan <- struct{}{}
}

func (to *TaskPayment) backgroundJobs() {
	fmt.Println("running payment backgroundJobs")
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			if to.removeExpiryDuration > 0 {
				err := to.removeExpiryPayment()
				if err != nil {
					fmt.Println("getting error from RemoveExpiryPayment", err)
				}
			}

		case <-to.closeChan:
			fmt.Println("payment task closed")
			return
		}
	}
}
