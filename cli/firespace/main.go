package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	cueyaml "cuelang.org/go/encoding/yaml"
	yaml "github.com/goccy/go-yaml"

	"cuelang.org/go/cue/cuecontext"
	"github.com/scaxyz/firespace"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()
var sugar = logger.Sugar()

func main() {
	defer logger.Sync()
	configFilePath := "../../test-validations.yml"

	config, err := loadYamlConfig(configFilePath)

	if err != nil {
		sugar.Fatalw("loading yaml config", zap.Error(err))
	}

	fmt.Println(config)

}

func loadYamlConfig(yamlPath string) (*firespace.ConfigFile, error) {

	yamlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		sugar.Errorw("reading yaml config", zap.String("path", yamlPath), zap.Error(err))
		return nil, err
	}

	cueValue, err := loadCueConfigValue()
	if err != nil {
		sugar.Errorw("loading cue contig value", zap.Error(err))
		return nil, err
	}

	err = validateYamlString(string(yamlData), cueValue)
	if err != nil {
		sugar.Errorw("validation src-yaml vs cue", zap.Error(err))
		return nil, err
	}

	configFile := firespace.ConfigFile{}

	err = yaml.Unmarshal(yamlData, &configFile)

	if err != nil {
		sugar.Errorw("unmarshaling yaml config file", zap.Error(err))
		return nil, err
	}

	remarshaledYAML, err := yaml.Marshal(&configFile)

	if err != nil {
		sugar.Errorw("marshaling struct  to yaml", zap.Error(err))
		return nil, err
	}

	err = validateYamlString(string(remarshaledYAML), cueValue)

	if err != nil {
		sugar.Errorw("validating re-marshaled via cue", zap.Error(err))
		return nil, err
	}

	return &configFile, nil

}

func loadCueConfigValue() (*cue.Value, error) {

	c, err := firespace.CueFiles.ReadDir("config")
	if err != nil {
		sugar.Errorw("reading embeded config directory", zap.Error(err))
		return nil, err
	}

	cctx := cuecontext.New()
	var valueNow cue.Value
	first := true

	for _, v := range c {
		sugar.Debugw("loading cue file", zap.String("file", v.Name()))

		filePath := filepath.Join("config", v.Name())

		schema, err := firespace.CueFiles.ReadFile(filePath)
		if err != nil {
			sugar.Errorw("reading embeded file", zap.Error(err), zap.String("filepath", filePath))
			return nil, err
		}
		opts := []cue.BuildOption{cue.Filename(filePath)}
		if !first {
			opts = append(opts, cue.Scope(valueNow))
		}

		value := cctx.CompileBytes(schema, opts...)
		if !first {
			value = valueNow.Unify(value)
			if err := value.Err(); err != nil {
				sugar.Errorw("loading cue files", zap.Error(err), zap.String("details", errors.Details(err, nil)))
				return nil, err
			}
		}
		valueNow = value

		first = false

	}

	return &valueNow, nil
}

func validateYamlString(yamlData string, cueValue *cue.Value) error {

	err := cueValue.Validate(cue.Optional(true), cue.All())
	if err != nil {
		sugar.Errorw("validating cue", zap.Error(err))
		return err
	}

	cf := cueValue.LookupPath(cue.ParsePath("#ConfigFile"))

	return cueyaml.Validate([]byte(yamlData), cf)

}
