package main

import (
	"fmt"
	"github.com/general252/cpu_percent/cpu_percent"
	"time"
)

func main() {
	for {
		if data, err := cpu_percent.GetCpuPercent(time.Second); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(data)
		}
	}
}
