import (
	"list"
	"strings"
)

#Settings: {
	allow_empty_home: bool | *false
	allow_debugger:   bool | *false

	firejail_flags: [...string]
    flags: [...string]
    pre_flags: [...string]
}

_ProgramSettings: {
	#Settings
	executeable: string
}

#OverwritesInline: {
	_src: _
	_ignore: [...string]

	_lower: {
		for k, v in _src {
			"\(k)"?: strings.ToLower(k)
		}
	}

	for k, v in _lower {
		if (!list.Contains(_ignore, k)) {

			"overwrite_\(v)"?: bool
		}
	}
	_src
}

#Overwrites: {
	_src: _
	_ignore: [...string]

	overwrites?: {
		for k, v in _src {
			if (!list.Contains(_ignore, k)) {

				"\(k)"?: bool | *false
			}
		}
	}
	_src
}

#ProgramSettings: #Overwrites & {
	_src: _ProgramSettings
	_ignore: ["executeable"]
}

c0: #ProgramSettings & {
	overwrites: {}

	executeable: "c0"
}

c1: #ProgramSettings & {
	overwrites: {}

	executeable: "c1"
}
