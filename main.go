package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

func envDef(name, def string) string {
	if v := os.Getenv("SERPENTINISED_" + name); v != "" {
		return v
	}

	return def
}

var (
	redisSentinelAddress = flag.String("sentinel-address", envDef("SENTINEL_ADDRESS", ""), "the address of the Sentinel master")
	redisSentinelMaster  = flag.String("sentinel-master", envDef("SENTINEL_MASTER", "mymaster"), "the name of the Sentinel master")
	bind                 = flag.String("bind", envDef("BIND", "127.0.0.1:26380"), "the address to bind to proxy connections to the active Sentinel")
	timeout              = flag.Int("connect-timeout", 1, "seconds before a connection to a master times out")

	logger *zap.Logger
)

func main() {
	flag.Parse()

	initLogger, err := zap.NewProduction()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to initialize logger: %v", err)
		os.Exit(1)
		return
	}
	logger = initLogger

	if *redisSentinelAddress == "" {
		logger.Error("A Redis Sentinel address is required.")
		os.Exit(1)
		return
	}

	logger.Info("connecting to redis sentinel", zap.String("sentinel-address", *redisSentinelAddress))
	monitor, err := NewRedisSentinelMonitor(*redisSentinelAddress, *redisSentinelMaster)
	if err != nil {
		logger.Info("unable to connect to redis sentinel",
			zap.String("sentinel-address", *redisSentinelAddress),
			zap.Error(err))
		os.Exit(1)
		return
	}

	server, err := NewSerpentinisedServer(*bind, monitor)
	if err != nil {
		logger.Info("unable to bind to address",
			zap.String("address", *bind),
			zap.Error(err))
		os.Exit(1)
		return
	}

	go monitor.monitorSentinelChanges()
	if err := server.Listen(); err != nil {
		logger.Info("error whilst listening for connections", zap.Error(err))
		os.Exit(1)
		return
	}
}
