package main

import (
	"fmt"
	"time"

	circuitbreaker "github.com/mitmadhu/circuit-breaker/circuit_breaker"
)

var (
	cb = circuitbreaker.NewCB(time.Second, 3, 2)
)

func main() {

	for i := 0; i < 30; i++ {
		err := cb.Run(func() error {
			if i%4 == 0 {
				return fmt.Errorf("[%d] error in call", i)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("[%d] error in call, error : %s \n", i, err.Error())
		} else {
			fmt.Printf("[%d] call was successful\n", i)
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func abcd() {
}
