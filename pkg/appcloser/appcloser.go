package appcloser

import (
	"fmt"
	"github.com/qiwik/synchronizer/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// CloseApp waits for interrupt signal
func CloseApp(logs *os.File, logger *logger.Logger) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info.Println("utility is closing by ctrl+c")
		logger.Info.Printf("log file is closing\n\n")
		err := logs.Close()
		if err != nil {
			log.Fatalf("log file was not closed\n\n")
		}
		fmt.Println(" - Done")
		os.Exit(0)
	}()
}
