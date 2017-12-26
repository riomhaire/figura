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
	go func() {
		<-c
		log.Println("Shuting Down")
		// pprof.StopCPUProfile()
		// //trace.Stop()
		// tracefile.Close()
		os.Exit(0)
	}()

	application := bootstrap.Application{}

	application.Initialize()
	application.Run()

}