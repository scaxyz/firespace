package main

import (
	_ "embed"

	"github.com/scaxyz/firespace"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()
var sugar = logger.Sugar()

func main() {
	defer logger.Sync()
	firespace.SetLogger(sugar)

	configFilePath := "../../test-validations.yml"

	config, err := firespace.LoadYamlConfig(configFilePath)

	if err != nil {
		sugar.Fatalw("loading yaml config", zap.Error(err))
	}

	space := firespace.NewFirespaceFromConfig(config, "user", "bat")

	eSpace := space.ExecuteTemplates()

	err = eSpace.Start([]string{"some-args-1", "some-args-2"}, true)

	if err != nil {
		sugar.Errorw("runnig firejail", zap.Error(err))
	}

}
