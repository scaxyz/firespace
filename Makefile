.PHONY: test-specs

test-specs:
	cue eval test-specs.yml config/specifications.cue -d=#ConfigFile

test-validations:
	cue eval test-validations.yml config/specifications.cue config/validations.cue -c -d=#ConfigFile


test-validations-fail:
	cue eval test-validations-fail.yml config/specifications.cue config/validations.cue -c -d=#ConfigFile