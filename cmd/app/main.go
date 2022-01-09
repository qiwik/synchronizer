package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/qiwik/synchronizer/internal/filestruct"
	"github.com/qiwik/synchronizer/internal/parameters/initial"
	scanning "github.com/qiwik/synchronizer/internal/scan"
	"github.com/qiwik/synchronizer/pkg/appcloser"
	"github.com/qiwik/synchronizer/pkg/logger"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	param                        initial.Parameters
	ticker                       *time.Ticker
	tick, tickerMain, tickerCopy int64
)

func main() {
	runtime.GOMAXPROCS(4)
	//creating a context and a log file with the required levels
	ctx := context.Background()
	logFile, logErr := logger.LogFileInit()
	if logErr != nil {
		log.Fatal(logErr)
	}
	logs := logger.LogInit(logFile)
	logs.Info.Println("log file created successfully")

	//waiting for the program to exit via ctrl+c
	appcloser.CloseApp(logFile, logs)

	//flags for interior parameters
	flag.StringVar(&param.SourcePath, "source", "", "The path of major folder")
	flag.StringVar(&param.CopyPath, "copy", "", "The path of copied folder")
	flag.Parse()

	logs.Info.Printf("initial data was read successfully with:\n"+
		"\tsource path: %s, copy path: %s\n", param.SourcePath, param.CopyPath)

	//file tree traversal at program start
	var folderStructure1, folderStructure2 filestruct.MainFolder
	var wg sync.WaitGroup
	var err error

	wg.Add(2)
	go func() {
		defer wg.Done()
		tickerMain, err = scanning.MainPathScan(param.SourcePath, &folderStructure1, logs)
		if err != nil {
			logs.Fatal.Fatal("error with scanning of the source directory")
		}
		logs.Info.Println("source directory scanned")
	}()

	go func() {
		defer wg.Done()
		tickerCopy, err = scanning.MainPathScan(param.CopyPath, &folderStructure2, logs)
		if err != nil {
			logs.Fatal.Fatal("error with scanning of the copy directory")
		}
		logs.Info.Println("copy directory scanned")
	}()
	wg.Wait()

	//set a ticker based on the total weight of the directory
	if tickerMain >= tickerCopy {
		tick = tickerMain
		ticker = time.NewTicker(time.Duration(tickerMain) * time.Second)
	} else {
		tick = tickerCopy
		ticker = time.NewTicker(time.Duration(tickerCopy) * time.Second)
	}

	//work with files, synchronization
	folderStructure1.FoldersSearch(param, &folderStructure2, logs, ctx)
	logs.Info.Printf("***** first iteration passed successfully *****\n\n")
	fmt.Println("Synchronized!")

	//background operations while not entered ctrl+c
	for {
		tickerPast := tick
		select {
		case <-ticker.C:
			stop := false
			var (
				folderStructure1, folderStructure2 filestruct.MainFolder
			)

			//scan file trees for new iteration
			wg.Add(2)
			go func() {
				defer wg.Done()
				tickerMain, err = scanning.MainPathScan(param.SourcePath, &folderStructure1, logs)
				if err != nil {
					logs.Fatal.Fatal("error with scanning of the source directory")
				}
				logs.Info.Println("source directory scanned")
			}()

			go func() {
				defer wg.Done()
				tickerCopy, err = scanning.MainPathScan(param.CopyPath, &folderStructure2, logs)
				if err != nil {
					logs.Fatal.Fatal("error with scanning of the copy directory")
				}
				logs.Info.Println("copy directory scanned")
			}()
			wg.Wait()

			//new tick for new files' weight
			if tickerMain >= tickerCopy {
				tick = tickerMain
			} else {
				tick = tickerCopy
			}

			if tickerPast < tick {
				stop = true
				ticker.Stop()
			} else {
				tickerPast = tick
				ticker.Reset(time.Duration(tick) * time.Second)
				fmt.Println(tick, "second")
			}

			//new work with files, synchronization
			folderStructure1.FoldersSearch(param, &folderStructure2, logs, ctx)
			logs.Info.Printf("***** new iteration passed successfully *****\n\n")
			fmt.Println("Synchronized!")

			//new current ticker
			if stop == true {
				ticker = time.NewTicker(time.Duration(tick) * time.Second)
			}
		}
	}
}
