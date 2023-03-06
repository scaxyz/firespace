package firespace2

#AbsPath: =~"^/" & #Path
#Path:    !~"^\\s*$"

#ConfigFile: {
	programms: [_name=string]: {
		executeable: #ProgramSettings.executeable | *_name
	}
}

#ConfigFile: {

	spaces: [string]: {
		allow_empty_home: #SpaceSettings.allow_empty_home

		home: (string) | (#AbsPath)

		if allow_empty_home == false {
			home: #AbsPath
		}

		if home != "" {
			no_private: false
		}

	}

}

