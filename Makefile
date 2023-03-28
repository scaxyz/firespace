.PHONY: test-specs

test-specs:
	cue eval ./testdata/test-specs.yml config/specifications.cue -d=#ConfigFile

test-validations:
	cue eval ./testdata/test-validations.yml config/specifications.cue config/validations.cue -c -d=#ConfigFile


test-validations-fail:
	cue eval ./testdata/test-validations-fail.yml config/specifications.cue config/validations.cue -c -d=#ConfigFile

eval-config:
	cue eval ~/.config/firespace/config.yaml config/specifications.cue config/validations.cue -c -d=#ConfigFile

install:
	GOBIN=~/go/bin go install ./cli/firespace/