package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/riomhaire/figura/frameworks/application/figura/bootstrap"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// tracefile, err := os.Create("app.trace")
	// check(err)

	// pprof.StartCPUProfile(tracefile)
	//	trace.Start(tracefile)
	// Shutdown
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	application := bootstrap.Application{}

	go func() {
		<-c
		log.Println("Shutting Down")
		// pprof.StopCPUProfile()
		// //trace.Stop()
		// tracefile.Close()
		application.Stop()
		os.Exit(0)
	}()

	application.Initialize()
	application.Run()

}
