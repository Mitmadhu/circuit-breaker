package main

import (
	"fmt"
	"time"

	"github.com/eapache/go-resiliency/breaker"
)

type callWithCB struct {
	cbObj *breaker.Breaker
}

func NewCB() callWithCB {
	return callWithCB{
		cbObj: breaker.New(3, 2, time.Second),
	}
}

func (cb *callWithCB) withCB(testCall func() error) error {

	err := cb.cbObj.Run(testCall)

	if err == nil {
		println("success in api call")
		return nil
	}

	if err == breaker.ErrBreakerOpen {
		println("circuit is open")
	} else {
		println("circuit is close and got error: " + err.Error())
	}

	return err

}

func main() {
	cb := NewCB()

	for i := 0; i < 50; i++ {
		cb.withCB(func() error {
			if i%4 == 0 {
				return fmt.Errorf("[%d] error in api call", i)
			}
			return nil
		})

		time.Sleep(100 * time.Millisecond)
	}

}
