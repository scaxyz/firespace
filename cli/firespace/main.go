package main

import (
	_ "embed"
	"flag"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/scaxyz/firespace"
	"go.uber.org/zap"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

var configPath = flag.String("config", filepath.Join(xdg.ConfigHome, "firespace", "config.yaml"), "Path to the config file.")
var dry = flag.Bool("dry", false, "Set run-dry mode.")
var debug = flag.Bool("debug", false, "Enable debug logging.")

func main() {
	defer func() {
		logger.Sync()
	}()
	flag.Parse()

	var err error

	if *debug {
		logger, err = zap.NewDevelopment(zap.IncreaseLevel(zap.DebugLevel))
	} else if *dry {
		logger, err = zap.NewDevelopment(zap.IncreaseLevel(zap.InfoLevel))
	} else {
		logger, err = zap.NewDevelopment(zap.IncreaseLevel(zap.WarnLevel))
	}

	if err != nil {
		panic(err)
	}

	sugar = logger.Sugar()

	if len(flag.Args()) < 2 {
		sugar.Error("please specify atleast a space to run")
		os.Exit(1)
	}

	firespace.SetLogger(sugar)

	config, err := firespace.LoadYamlConfig(*configPath)

	if err != nil {
		sugar.Panicw("loading yaml config", zap.Error(err), zap.String("config", *configPath))
	}

	if len(config.Spaces) == 0 {
		sugar.Error("please create at leaset one space in your config")
		os.Exit(1)
	}

	spaceName := flag.Arg(0)
	programName := flag.Arg(1)

	space := firespace.NewFirespaceFromConfig(config, spaceName, programName)

	if space.Executeable == "" {
		space.Executeable = programName
	}

	eSpace := space.ExecuteTemplates()

	err = eSpace.Start(flag.Args()[2:], *dry)

	if err != nil {
		sugar.Panicw("runnig firejail", zap.Error(err))
	}

}
