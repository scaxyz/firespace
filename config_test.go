package firespace

import (
	"testing"
	"unsafe"

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

	closedCueGoInterface := cueGoInterface
	op, orEntries := cueGoInterface.Expr()

	// if op is | and first entry is nil use the sencond entry
	if op.String() == "|" {
		if orEntries[0].Null() == nil {
			closedCueGoInterface = orEntries[1]
		}
	}

	// close the struct
	closedCueGoInterface = closeStructs(cCtx, closedCueGoInterface)

	t.Logf("closedCueGoInterface:\n%#v\n", closedCueGoInterface)

	// cross unify to check of not allow or missing entries

	cueGoUnify := cueLookupValue.Unify(closedCueGoInterface)
	require.NoErrorf(t, cueGoUnify.Err(), "cue.unify(go): %s", errors.Details(cueGoUnify.Err(), nil))

	goCueUnify := closedCueGoInterface.Unify(cueLookupValue)
	require.NoErrorf(t, goCueUnify.Err(), "go.unify(cue): %s", errors.Details(cueGoUnify.Err(), nil))
}

func closeStructs(cctx *cue.Context, v cue.Value) cue.Value {

	type UnsafeVertex struct {
		// Parent links to a parent Vertex. This parent should only be used to
		// access the parent's Label field to find the relative location within a
		// tree.
		Parent uintptr

		// Label is the feature leading to this vertex.
		Label uint32

		// State:
		//   eval: nil, BaseValue: nil -- unevaluated
		//   eval: *,   BaseValue: nil -- evaluating
		//   eval: *,   BaseValue: *   -- finalized
		//
		state uintptr
		// TODO: move the following status fields to nodeContext.

		// status indicates the evaluation progress of this vertex.
		status int8

		// isData indicates that this Vertex is to be interepreted as data: pattern
		// and additional constraints, as well as optional fields, should be
		// ignored.
		isData                bool
		Closed                bool
		nonMonotonicReject    bool
		nonMonotonicInsertGen int32
		nonMonotonicLookupGen int32
		// ...
	}

	// Value holds any value, which may be a Boolean, Error, List, Null, Number,
	// Struct, or String.
	type UnsafeValue struct {
		idx uintptr
		v   *UnsafeVertex
		// Parent keeps track of the parent if the value corresponding to v.Parent
		// differs, recursively.
		parent_ uintptr
	}

	pV := &v
	pUV := (*UnsafeValue)(unsafe.Pointer(pV))
	pUV.v.Closed = true

	return v
}
