package firespace2

#ENV: [string]: string

#HasENV: {
	env?: #ENV | null
}

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

#GlobalSettings: {
	#CommonSettings
}

#CommonSettings: {
	#HasENV
	before: [...string]
	after: [...string]
	allow_debugger: bool |*false
	firejail_flags: [...string]
}


#ProgramSettings: {
	#CommonSettings

	#HasOverwrites

	spaces: [string]: #AdditionalSpaceSettings

	executeable: string

	pre_flags: [...string]
	flags: [...string]

}

#SpaceSettings: {
	#CommonSpaceSettings
	#CanControllHome
	#CanSetHome
}

#CommonSpaceSettings:{
	#CommonSettings
	#HasOverwrites
	
}

#AdditionalSpaceSettings:{
	#CommonSpaceSettings
}


#HasAddionalSpacesSettings:{
	spaces: [string]: #AdditionalSpaceSettings
}

#ConfigFile: {
	global?: #GlobalSettings

	spaces: [string]: #SpaceSettings

	programms?: [string]: #ProgramSettings
}


#FireSpaceContext:{
	#HasENV
	#CommonSettings
	#CanControllHome
	#CanSetHome
	#ProgramSettings
}