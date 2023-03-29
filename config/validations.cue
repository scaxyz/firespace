package firespace

#AbsPath: =~"^/" & #Path
#Path:    !~"^\\s*$"

#ProgramSettings: _
#SpaceSettings:   _

#ProgramMap: [_name=string]: {
	executeable: #ProgramSettings.executeable | *_name
}

#SpaceMap: [string]: {
	allow_empty_home: #SpaceSettings.allow_empty_home

	home: (string) | (#AbsPath)

	if allow_empty_home == false {
		home: #AbsPath
	}

	if home != "" {
		no_private: false
	}

}

#CommonSettings: {
	_#Flag: =~"^--"
	firejail_flags: [..._#Flag]
}
