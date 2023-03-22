package firespace

import (
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/gocode/gocodec"
	"github.com/stretchr/testify/require"
)

func Test_ConfigFileCueValidion(t *testing.T) {

	goConfigFile := ConfigFile{}

	validateGoStruct(t, &goConfigFile, "#ConfigFile")

}

func Test_GlobalSettingsCueValidion(t *testing.T) {

	goGlobalSettings := GlobalSettings{}

	validateGoStruct(t, &goGlobalSettings, "#GlobalSettings")

}

func Test_SpaceSettingsCueValidion(t *testing.T) {

	goSpaceSettings := SpaceSettings{}

	validateGoStruct(t, &goSpaceSettings, "#SpaceSettings")

}

func Test_ProgramSettingsCueValidion(t *testing.T) {

	goProgramSettings := ProgramSettings{}

	validateGoStruct(t, &goProgramSettings, "#ProgramSettings")

}

func Test_AdditionalSpaceSettingsCueValidion(t *testing.T) {

	goAdditionalSpaceSettings := AdditionalSpaceSettings{}

	validateGoStruct(t, &goAdditionalSpaceSettings, "#AdditionalSpaceSettings")

}

func validateGoStruct(t *testing.T, goValue interface{}, cuePath string) {

	require.Empty(t, goValue)

	cCtx := cuecontext.New()

	ccVal, err := loadCueConfigValue(cCtx)
	require.NoError(t, err)
	require.NotEmpty(t, ccVal)

	cueConfigFile := ccVal.LookupPath(cue.ParsePath(cuePath))
	require.NotEmpty(t, cueConfigFile)

	cCodec := gocodec.New((*cue.Runtime)(cCtx), &gocodec.Config{})

	err = cCodec.Validate(cueConfigFile, goValue)
	require.NoError(t, err)

}
