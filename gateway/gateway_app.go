package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "q,r,t,y,i,o,g,h"
	arr := strings.Split(s, ",")
	for i := range arr {
		arr[i] = "*"
		join := strings.Join(arr, ".")
		fmt.Println(join)
	}
}
