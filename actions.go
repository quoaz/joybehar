package main

import "fmt"

type Action func()

func message(s string) Action {
	return Action(func() {
		fmt.Println(s)
	})
}
