package main


import (
	"github.com/jilieryuyi/wing-go/tcp"
	"context"
	"time"
	"fmt"
	"github.com/sirupsen/logrus"
)
func main() {
	address := "127.0.0.1:7770"
	client  := tcp.NewClient(context.Background())
	err     := client.Connect(address, time.Second * 3)

	if err != nil {
		logrus.Errorf("connect to %v error: %v", address, err)
	}
	defer client.Disconnect()

	w1, _   := client.Send([]byte("hello"))
	w2, _   := client.Send([]byte("word"))
	w3, _   := client.Send([]byte("hahahahahahahahahahah"))
	res1, _ := w1.Wait(time.Second * 3)
	res2, _ := w2.Wait(time.Second * 3)
	res3, _ := w3.Wait(time.Second * 3)

	fmt.Println("w1 return: ", string(res1))
	fmt.Println("w2 return: ", string(res2))
	fmt.Println("w3 return: ", string(res3))
}
