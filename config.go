package firespace

type ConfigFile struct {
	Global    GlobalSettings             `yaml:"global"`
	Spaces    map[string]SpaceSettings   `yaml:"spaces"`
	Programms map[string]ProgramSettings `yaml:"programms"`
}
type Env map[string]string

type HasEnv struct {
	Env Env
}

type CanSetHome struct {
	Home string
}

type CanControllHome struct {
	AllowEmptyHome bool `yaml:"allow_empty_home"`
	NoPrivate      bool `yaml:"no_private"`
}

type HasOverwrites struct {
	Overwrites Overwrites
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
	CommonSettings `yaml:",inline"`
}

type SpaceSettings struct {
	CommonSpaceSettings `yaml:",inline"`
	CanControllHome     `yaml:",inline"`
	CanSetHome          `yaml:",inline"`
}

type CommonSpaceSettings struct {
	CommonSettings `yaml:",inline"`
	HasOverwrites  `yaml:",inline"`
}

type ProgramSettings struct {
	CommonSettings `yaml:",inline"`
	HasOverwrites  `yaml:",inline"`
	Spaces         map[string]AddionalSpacesSettings
	Executeable    string
	PreFlags       []string `yaml:"pre_flags"`
	Flags          []string
}

type CommonSettings struct {
	HasEnv        `yaml:",inline"`
	Before        []string
	After         []string
	FirejailFlags []string `yaml:"firejail_flags"`
}

type AddionalSpacesSettings struct {
	CommonSpaceSettings `yaml:",inline"`
}
