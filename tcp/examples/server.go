package main

import (
	"github.com/jilieryuyi/wing-go/tcp"
	"context"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	_ "net/http/pprof"
	"github.com/sirupsen/logrus"
	"fmt"
)
func main() {
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
		if err:=recover();err!=nil{
			fmt.Println("recover error: ", err) // 这里的err其实就是panic传入的内容，55
		}
	}()
	address := "127.0.0.1:7771"
	server  := tcp.NewServer(context.Background(), address, tcp.SetOnServerMessage(func(node *tcp.ClientNode, msgId int64, data []byte) {
		logrus.Infof("server send, msgid=[%v], data=[%v]", msgId, string(data))
		n,e := node.Send(msgId, data)
		if e != nil || n != len(data) {
			logrus.Errorf("send fail")
		}
	}))
	server.Start()
	defer server.Close()


	go func() {
		//http://localhost:8880/debug/pprof/  内存性能分析工具
		//go tool pprof logDemo.exe --text a.prof
		//go tool pprof your-executable-name profile-filename
		//go tool pprof your-executable-name http://localhost:8880/debug/pprof/heap
		//go tool pprof wing-binlog-go http://localhost:8880/debug/pprof/heap
		//https://lrita.github.io/2017/05/26/golang-memory-pprof/
		//然后执行 text
		//go tool pprof -alloc_space http://127.0.0.1:8880/debug/pprof/heap
		//top20 -cum

		//下载文件 http://localhost:8880/debug/pprof/profile
		//分析 go tool pprof -web /Users/yuyi/Downloads/profile
		http.ListenAndServe("127.0.0.1:7772", nil)
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sc
}