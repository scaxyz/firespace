package firespace

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
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

func validateGoStruct(t *testing.T, goInterface interface{}, cuePath string) {
	require.Empty(t, goInterface)

	cCtx := cuecontext.New()

	ccVal, err := loadCueConfigValue(cCtx)
	require.NoError(t, err)
	require.NotEmpty(t, ccVal)

	cueLookupValue := ccVal.LookupPath(cue.ParsePath(cuePath))
	require.NotEmpty(t, cueLookupValue)
	t.Logf("cueLookupValue:\n%#v\n", cueLookupValue)

	cueGoInterface := cCtx.EncodeType(goInterface)

	openCueGoInterface := cueGoInterface
	op, orEntries := cueGoInterface.Expr()

	// if op is | and first entry is nil use the sencond entry
	if op.String() == "|" {
		if orEntries[0].Null() == nil {
			openCueGoInterface = orEntries[1]
		}
	}

	// close the struct
	closedCueGoInterface := closeStructs(cCtx, openCueGoInterface)

	t.Logf("closedCueGoInterface:\n%#v\n", closedCueGoInterface)

	// cross unify to check for not allow or missing entries

	cueGoUnify := cueLookupValue.Unify(closedCueGoInterface)
	require.NoErrorf(t, cueGoUnify.Err(), "cue.unify(go): %s", errors.Details(cueGoUnify.Err(), nil))

	goCueUnify := closedCueGoInterface.Unify(cueLookupValue)
	require.NoErrorf(t, goCueUnify.Err(), "go.unify(cue): %s", errors.Details(cueGoUnify.Err(), nil))
}

func closeStructs(cctx *cue.Context, v cue.Value) cue.Value {
	closedV := cctx.CompileString(fmt.Sprintf("close({%#v})", v))

	return closedV

}
