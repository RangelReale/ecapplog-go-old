package main

import (
	"fmt"
	"time"

	"github.com/RangelReale/ecapplog-go"
)

func main() {
	c := ecapplog.NewClient()
	c.Open()
	defer c.Close()

	for i := 0; i < 100; i++ {
		c.Log(time.Now(), ecapplog.Priority_DEBUG, "app", fmt.Sprintf("First log: %d", i))
		c.Log(time.Now(), ecapplog.Priority_INFORMATION, "app", fmt.Sprintf("Second log: %d", i))
		c.Log(time.Now(), ecapplog.Priority_ERROR, "app", fmt.Sprintf("Third log: %d", i))
	}

	fmt.Printf("Sleeping 10 seconds...\n")
	select {
	case <-time.After(time.Second * 10):
		break
	}

	fmt.Printf("Closing and sleeping 10 seconds...\n")
	c.Close()
	select {
	case <-time.After(time.Second * 10):
		break
	}
}
