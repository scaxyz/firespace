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

#HasSpaces: {
	spaces?: [string]: #SpaceSettings
}

#ProgramOnlySettings: {
	#Settings
	executeable: string
}

#ProgramSettings: #Overwrites & {
	_src: {
		#ProgramOnlySettings
		#HasSpaces
	}
	_ignore: ["executeable", "spaces"]

}

#HomeSettings: {
	home:             string
	allow_empty_home: bool | *false
}

#SpaceOnlySettings: {
	#Settings
	#HomeSettings
}

#SpaceSettings: #Overwrites & {
	_src: #SpaceOnlySettings
	_ignore: ["home"]
}

#GlobalSettings: {
	#Settings
	#HomeSettings
}

#ConfigFile: {
	global?: #GlobalSettings

	#HasSpaces

	programms?: [string]: #ProgramSettings
}

#FireSpaceContext: {
	#SpaceOnlySettings
	#ProgramOnlySettings
}

// validations
#ConfigFile: programms?: [_program=string]: #ProgramSettings & {
	executeable: _program | =~"/"
}

// in programms only allow space names that are allowed at the top level
#ConfigFile: {
	spaces?: _
	#allowedSpaces: { // also works with closed hidden types; (performance?)
		for k, _ in spaces {
			"\(k)": true
		}
	}
	programms?: [string]: {
		spaces?: [_space=string]: {
			_valid_space_name: #allowedSpaces & {
				"\(_space)": true
			}
		}
	}
}

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

c1: #ConfigFile & {
	global: {}
	spaces:{
		valid:{}
	}
	programms:{
		cat:{
			spaces:{
				valid:{}
			}
		}
	}
}

f0: #FireSpaceContext & {}
