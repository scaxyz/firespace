package firespace

#ConfigFile: {
	global?: #GlobalSettings

	spaces?: [string]: #SpaceSettings

	programms?: [string]: #ProgramSettings
}

#GlobalSettings: {
	#CommonSettings
}

#SpaceSettings: {
	#CommonSpaceSettings
	#CanControllHome
	#CanSetHome
}

#ProgramSettings: {
	#CommonProgramSettings
	spaces: [string]: #AdditionalSpaceSettings
	executeable: string
}


#FireSpaceContext:{
	#HasENV
	#CommonSettings
	#CanControllHome
	#CanSetHome
	#ProgramSettings
}

#CommonSettings: {
	#HasENV
	before: [...#ExtendetShellCommand]
	after: [...#ShellCommand]
	allow_debugger: bool |*false
	firejail_flags: [...string]
}

#CommonProgramSettings:{
	#CommonSettings

	#HasOverwrites

	pre_flags: [...string]
	flags: [...string]
}

#CommonSpaceSettings:{
	#CommonSettings
	#HasOverwrites
	
}

#AdditionalSpaceSettings:{
	#CommonSpaceSettings
	#CommonProgramSettings
}

#HasENV: {
	env?: #ENV | null
}

#ENV: [string]: string


#CanSetHome:{
	home: string
}

#CanControllHome:{
	allow_empty_home: bool | *false
	no_private: bool | *false
}


#HasOverwrites:{
	overwrites: [string]: bool
}

#HasAddionalSpacesSettings:{
	spaces: [string]: #AdditionalSpaceSettings
}


#ExtendetShellCommand:{
	command: #ShellCommand
	allow_error: bool | *false
}

#ShellCommand: string





