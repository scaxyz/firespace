package firespace

type ConfigFile struct {
	Global    GlobalSettings             `json:"global"`
	Spaces    map[string]SpaceSettings   `json:"spaces"`
	Programms map[string]ProgramSettings `json:"programms"`
}
type Env map[string]string

type HasEnv struct {
	Env Env `json:"env"`
}

type CanSetHome struct {
	Home string `json:"home"`
}

type CanControllHome struct {
	AllowEmptyHome bool `json:"allow_empty_home"`
	NoPrivate      bool `json:"no_private"`
}

type HasOverwrites struct {
	Overwrites Overwrites `json:"overwrites"`
}

type Overwrites struct {
	Env           bool
	Before        bool
	After         bool
	FirejailFlags bool
	PreFlags      bool
	Flags         bool
}

type GlobalSettings struct {
	/*
		#CommonSettings
	*/
	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	Env Env `json:"env"`
}

type SpaceSettings struct {
	/// CommonSettings `json:",inline"
	HasEnv        `json:",inline"`
	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	// HasOverwrites   `json:",inline"`
	Overwrites Overwrites `json:"overwrites"`

	// CanControllHome `json:",inline"`
	AllowEmptyHome bool `json:"allow_empty_home"`
	NoPrivate      bool `json:"no_private"`

	// CanSetHome      `json:",inline"`
	Home string `json:"home"`
}

type CommonSpaceSettings struct {
	// CommonSettings `json:",inline"`
	// HasEnv        `json:",inline"`
	Env Env `json:"env"`

	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	// HasOverwrites  `json:",inline"`
	Overwrites Overwrites `json:"overwrites"`
}

type ProgramSettings struct {
	// CommonProgramSettings `json:",inline"`
	// CommonSettings `json:",inline"`
	// HasEnv        `json:",inline"`
	Env Env `json:"env"`

	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	// HasOverwrites  `json:",inline"`
	Overwrites Overwrites `json:"overwrites"`

	PreFlags []string `json:"pre_flags"`
	Flags    []string `json:"flags"`

	Spaces      map[string]AdditionalSpaceSettings `json:"spaces"`
	Executeable string                             `json:"executeable,omitempty"`
}

type CommonProgramSettings struct {
	// CommonSettings `json:",inline"`
	// HasEnv        `json:",inline"`
	Env Env `json:"env"`

	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	// HasOverwrites  `json:",inline"`
	Overwrites Overwrites `json:"overwrites"`

	PreFlags []string `json:"pre_flags"`
	Flags    []string `json:"flags"`
}

type CommonSettings struct {
	// HasEnv        `json:",inline"`
	Env Env `json:"env"`

	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`
}

type ShellCommand string
type ShellCommands []string

type ExtendedShellCommand struct {
	Command    string `json:"command"`
	AllowError bool   `json:"allow_error"`
}

type AdditionalSpaceSettings struct {
	// CommonSpaceSettings   `json:",inline"`
	// CommonSettings `json:",inline"`
	// HasEnv        `json:",inline"`
	Env Env `json:"env"`

	Before        []ExtendedShellCommand `json:"before,omitempty"`
	After         ShellCommands          `json:"after,omitempty"`
	FirejailFlags []string               `json:"firejail_flags"`

	// HasOverwrites `json:",inline"`
	Overwrites Overwrites `json:"overwrites"`

	//CommonProgramSettings `json:",inline"`

	PreFlags []string `json:"pre_flags"`
	Flags    []string `json:"flags"`
}
