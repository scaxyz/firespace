package firespace

import (
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/encoding/gocode/gocodec"
	cueyaml "cuelang.org/go/encoding/yaml"

	"go.uber.org/zap"
)

func LoadYamlConfig(yamlPath string) (*ConfigFile, error) {

	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		sugar.Errorw("reading yaml config", zap.String("path", yamlPath), zap.Error(err))
		return nil, err
	}

	cctx := cuecontext.New()

	cueValue, err := loadCueConfigValue(cctx)
	if err != nil {
		sugar.Errorw("loading cue config value", zap.Error(err))
		return nil, err
	}

	i, err := cueyaml.Decode((*cue.Runtime)(cctx), "config.yaml", yamlData)
	if err != nil {
		sugar.Errorw("decoding yaml to cue value", zap.Error(err))
		return nil, err
	}
	ymlValue := i.Value()

	configFileValue := cueValue.LookupPath(cue.ParsePath("#ConfigFile"))

	err = configFileValue.Err()
	if err != nil {
		sugar.Errorw("looging up cue value", zap.Error(err), zap.String("details", errors.Details(err, nil)))
		return nil, err
	}

	unifiedConfig := configFileValue.Unify(ymlValue)
	err = unifiedConfig.Err()
	if err != nil {
		sugar.Errorw("unifying cue schema and yaml config", zap.Error(err), zap.String("details", errors.Details(err, nil)))
		return nil, err
	}

	cueCodec := gocodec.New((*cue.Runtime)(cctx), nil)

	configFile := ConfigFile{}

	err = cueCodec.Validate(unifiedConfig, &configFile)

	if err != nil {
		sugar.Errorw("validating go struct", zap.Error(err), zap.String("details", errors.Details(err, nil)))
		return nil, err
	}

	err = cueCodec.Complete(unifiedConfig, &configFile)
	if err != nil {
		sugar.Errorw("completing go struct", zap.Error(err), zap.String("details", errors.Details(err, nil)))
		return nil, err
	}

	return &configFile, nil

}

func loadCueConfigValue(cctx *cue.Context) (*cue.Value, error) {

	c, err := CueFiles.ReadDir("config")
	if err != nil {
		sugar.Errorw("reading embeded config directory", zap.Error(err))
		return nil, err
	}

	var valueNow cue.Value
	first := true

	for _, v := range c {
		sugar.Debugw("loading cue file", zap.String("file", v.Name()))

		filePath := filepath.Join("config", v.Name())

		schema, err := CueFiles.ReadFile(filePath)
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
