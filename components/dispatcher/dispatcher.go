package main

import (
	"os"
	"syscall"

	"flag"

	_ "net/http/pprof"

	"os/signal"

	"runtime/debug"

	"github.com/xiaonanln/goworld/engine/binutil"
	"github.com/xiaonanln/goworld/engine/config"
	"github.com/xiaonanln/goworld/engine/consts"
	"github.com/xiaonanln/goworld/engine/gwlog"
)

var (
	configFile = ""
	sigChan    = make(chan os.Signal, 1)
)

func parseArgs() {
	flag.StringVar(&configFile, "configfile", "", "set config file path")
	flag.Parse()
}

func setupGCPercent() {
	debug.SetGCPercent(consts.DISPATCHER_GC_PERCENT)
}

func main() {
	setupGCPercent()
	parseArgs()

	if configFile != "" {
		config.SetConfigFile(configFile)
	}

	dispatcherConfig := config.GetDispatcher()
	binutil.SetupGWLog("dispatcher", dispatcherConfig.LogLevel, dispatcherConfig.LogFile, dispatcherConfig.LogStderr)
	setupSignals()
	binutil.SetupHTTPServer(dispatcherConfig.HTTPIp, dispatcherConfig.HTTPPort, nil)

	dispatcher := newDispatcherService()
	dispatcher.run()
}

func setupSignals() {
	signal.Ignore(syscall.Signal(10), syscall.Signal(12))
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			sig := <-sigChan

			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				// interrupting, quit dispatcher
				gwlog.Infof("Dispatcher quited.")
				os.Exit(0)
			} else {
				gwlog.Infof("unexcepted signal: %s", sig)
			}
		}
	}()
}
