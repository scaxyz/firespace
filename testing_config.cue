import (
	"list"
	"strings"
)
// Generators

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

// Schema

#Settings: {
	allow_debugger: bool | *false

	firejail_flags: [...string]
	flags: [...string]
	pre_flags: [...string]

	env?: [string]: string
}

_ProgramSettings: {
	#Settings
	executeable: string
}

#ProgramSettings: #Overwrites & {
	_src: _ProgramSettings
	_ignore: ["executeable"]
}

#HomeSettings: {
	home:             string
	allow_empty_home: bool | *false
}

_SpaceSettings: {
	#Settings
	#HomeSettings
}

#SpaceSettings: #Overwrites & {
	_src: _SpaceSettings
	_ignore: ["home"]
}

#GlobalSettings: {
	#Settings
	#HomeSettings
}

#ConfigFile: {
	global?: #GlobalSettings

	spaces?: #SpaceSettings

	programms?: #ProgramSettings
}

// testing objects


p0: #ProgramSettings & {
	overwrites: {
	}
	executeable: "p0"

	spaces?: [string]:#SpaceSettings

}

p1: #ProgramSettings & {

	executeable: "p1"
}

s0: #SpaceSettings & {
	overwrites: {}
}

s1: #SpaceSettings & {
}

c0: #ConfigFile & {
	global: {}
}
