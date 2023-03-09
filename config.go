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
	CommonSettings `json:",inline"`
}

type SpaceSettings struct {
	CommonSpaceSettings `json:",inline"`
	CanControllHome     `json:",inline"`
	CanSetHome          `json:",inline"`
}

type CommonSpaceSettings struct {
	CommonSettings `json:",inline"`
	HasOverwrites  `json:",inline"`
}

type ProgramSettings struct {
	CommonSettings `json:",inline"`
	HasOverwrites  `json:",inline"`
	Spaces         map[string]AddionalSpacesSettings `json:"spaces"`
	Executeable    string                            `json:"executeable,omitempty"`
	PreFlags       []string                          `json:"pre_flags"`
	Flags          []string                          `json:"flags"`
}

type CommonSettings struct {
	HasEnv        `json:",inline"`
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

type AddionalSpacesSettings struct {
	CommonSpaceSettings `json:",inline"`
}
