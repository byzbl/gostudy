package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	go sendData(ch1)

	// 读取数据
	for {
		time.Sleep(1 * time.Second)
		v, ok := <-ch1
		if !ok {
			fmt.Println("已经读取了所有数据。。", ok, v)
			break
		}
		fmt.Println("读取的数据是：", v, ok)
	}

	fmt.Println("main over..")
}

func sendData(ch1 chan int) {
	for i := 0; i < 10; i++ {
		ch1 <- i
	}
	close(ch1)
}
