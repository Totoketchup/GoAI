package main

import (
	"fmt"
	"time"
	)

//testing some functions
func main() {

 fmt.Println(time.Now().Unix())

 c := time.Tick(1 * time.Second)
    for now := range c {
            fmt.Println("Hello" ,now.Unix())
    }

}