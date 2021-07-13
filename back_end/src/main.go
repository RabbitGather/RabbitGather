package main

import (
	"context"
	//"fmt"
	//"log"
	"os"
	"os/signal"
	"rabbit_gather/src/api_server"
	"rabbit_gather/src/logger"
	"rabbit_gather/src/reverse_proxy_server"
	"rabbit_gather/src/web_server"
	"syscall"
	// database init
	"rabbit_gather/src/neo4j_db"
)

var log = logger.NewLogger("main")

func main() {
	log.DEBUG.Println("Start Main.")
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	reverseProxyServer := reverse_proxy_server.ReverseProxyServer{}
	err := reverseProxyServer.Startup(ctx, shutdownCallback)
	if err != nil {
		cancle()
		panic(err.Error())
	}

	ctx1, _ := context.WithCancel(ctx)
	webserver := web_server.WebServer{}
	err = webserver.Startup(ctx1, shutdownCallback)
	if err != nil {
		cancle()
		panic(err.Error())
	}
	ctx2, _ := context.WithCancel(ctx)
	apiServer := api_server.APIServer{}
	err = apiServer.Startup(ctx2, shutdownCallback)
	if err != nil {
		cancle()
		panic(err.Error())
	}

	//ctx3, _ := context.WithCancel(ctx)
	//websocketServer := websocket_server.WebsocketServer{}
	//err = websocketServer.Startup(ctx3, shutdownCallback)
	//if err != nil {
	//	cancle()
	//	panic(err.Error())
	//}

	waitForShutdown(ctx)
	log.DEBUG.Println("Main process end.")
	finalize()
}

func finalize() {
	err := neo4j_db.Close()
	if err != nil {
		panic(err.Error())
	}
}

func waitForShutdown(ctx context.Context) {
	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.DEBUG.Println("Shutdown with Context done")
	case <-quitSignal:
		log.DEBUG.Println("Shutdown with QuitSignal")
	}
	runShutdownCallbacks()
}

func runShutdownCallbacks() {
	for _, f := range shutdownCallbackQueue {
		f()
	}
}

var shutdownCallbackQueue = []func(){}

func shutdownCallback(f func()) {
	shutdownCallbackQueue = append(shutdownCallbackQueue, f)
}
