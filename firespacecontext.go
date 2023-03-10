package firespace

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/google/shlex"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
)

type FirespaceContext struct {
	CommonSettings
	CanControllHome
	CanSetHome
	HasOverwrites `yaml:",inline"`
	Executeable   string
	PreFlags      []string `yaml:"pre_flags"`
	Flags         []string
}

func NewFirespaceFromConfig(config *ConfigFile, space string, program string) *FirespaceContext {

	fromGlobal := &FirespaceContext{
		CommonSettings: config.Global.CommonSettings,
	}

	fromSpace := newFirespaceFromSpace(config, space)

	fromProgram := newFirespaceFromProgram(config, program)

	fromProgramSpace := newFirespaceFromProgramSpace(config, program, space)

	merged := Merge(fromGlobal, fromSpace)

	merged = Merge(merged, fromProgram)

	merged = Merge(merged, fromProgramSpace)

	return merged
}

func Merge(base *FirespaceContext, updates *FirespaceContext) *FirespaceContext {

	merged := FirespaceContext{}

	if base == nil && updates == nil {
		return &FirespaceContext{}
	}

	if updates == nil {
		merged = *base
		return &merged
	}

	if base == nil {
		merged = *updates
		return &merged
	}

	m := FirespaceContext{
		CommonSettings: CommonSettings{
			HasEnv: HasEnv{
				Env: mergeOrReplaceMap(base.Env, updates.Env, updates.Overwrites.Env),
			},
			Before:        mergeSliceOrOverwrite(base.Before, updates.Before, updates.Overwrites.Before),
			After:         mergeSliceOrOverwrite(base.After, updates.After, updates.Overwrites.After),
			FirejailFlags: mergeSliceOrOverwrite(base.FirejailFlags, updates.FirejailFlags, updates.Overwrites.FirejailFlags),
		},
		CanControllHome: CanControllHome{
			AllowEmptyHome: base.AllowEmptyHome || updates.AllowEmptyHome,
			NoPrivate:      base.NoPrivate || updates.NoPrivate,
		},
		CanSetHome: CanSetHome{
			Home: updateIfNotEmpty(base.Home, updates.Home),
		},
		Executeable: updateIfNotEmpty(base.Executeable, updates.Executeable),
		PreFlags:    mergeSliceOrOverwrite(base.PreFlags, updates.PreFlags, updates.Overwrites.PreFlags),
		Flags:       mergeSliceOrOverwrite(base.Flags, updates.Flags, updates.Overwrites.Flags),
	}

	return &m
}

func updateIfNotEmpty[V comparable](base, update V) V {
	var zeroV V
	if update == zeroV {
		return base
	}
	return update
}

func mergeSliceOrOverwrite[V any](base, update []V, overwrite bool) []V {
	newSlice := []V{}

	if !overwrite {
		newSlice = append(newSlice, base...)
	}

	newSlice = append(newSlice, update...)
	return newSlice

}

// mergeOrReplaceMap merges the too map one by one and replacing existing key, or if overwrite is true replaces the whole map
func mergeOrReplaceMap[K comparable, V any, M map[K]V](base, update map[K]V, overwrite bool) M {

	newMap := M{}

	if !overwrite {
		for k, v := range base {
			newMap[k] = v
		}
	}

	for k, v := range update {
		newMap[k] = v
	}

	return newMap

}

func newFirespaceFromSpace(config *ConfigFile, space string) *FirespaceContext {

	if config == nil {
		return nil
	}

	spaceSettings, ok := config.Spaces[space]
	if !ok {
		return nil
	}

	return &FirespaceContext{
		CommonSettings:  spaceSettings.CommonSettings,
		CanControllHome: spaceSettings.CanControllHome,
		CanSetHome:      spaceSettings.CanSetHome,
		HasOverwrites:   spaceSettings.HasOverwrites,
	}
}

func newFirespaceFromProgram(config *ConfigFile, program string) *FirespaceContext {

	if config == nil {
		return nil
	}

	programSettings, ok := config.Programms[program]
	if !ok {
		return nil
	}

	return &FirespaceContext{
		CommonSettings: programSettings.CommonSettings,
		HasOverwrites:  programSettings.HasOverwrites,
		Executeable:    programSettings.Executeable,
		PreFlags:       programSettings.PreFlags,
		Flags:          programSettings.Flags,
	}
}

func newFirespaceFromProgramSpace(config *ConfigFile, program string, space string) *FirespaceContext {

	if config == nil {
		return nil
	}

	if program == "" {
		return nil
	}

	if space == "" {
		return nil
	}

	programSettings, ok := config.Programms[program]

	if !ok {
		return nil
	}

	spaceSettings, ok := programSettings.Spaces[space]
	if !ok {
		return nil
	}

	return &FirespaceContext{
		CommonSettings: spaceSettings.CommonSettings,
		HasOverwrites:  spaceSettings.HasOverwrites,
	}

}

type TemplateContext struct {
	OS struct {
		Env map[string]string
	}
	Space struct {
		Env    map[string]string
		Config FirespaceContext
	}
}

func getOsEnvMap() map[string]string {
	envMap := map[string]string{}
	for _, v := range os.Environ() {
		splits := strings.SplitN(v, "=", 2)
		key := splits[0]
		rest := splits[1]
		envMap[key] = rest

	}
	return envMap
}

func envMapToSlice(strMap map[string]string) []string {
	envSlice := []string{}
	for k, v := range strMap {
		envSlice = append(envSlice, fmt.Sprint(k, "=", v))
	}
	return envSlice
}

func (space FirespaceContext) ExecuteTemplates() *FirespaceContext {

	templateContext := TemplateContext{
		OS: struct{ Env map[string]string }{
			Env: getOsEnvMap(),
		},
		Space: struct {
			Env    map[string]string
			Config FirespaceContext
		}{
			Env:    space.Env,
			Config: space,
		},
	}

	newSpace := FirespaceContext{
		CommonSettings: CommonSettings{
			HasEnv: HasEnv{Env: executeTemplateOnMap(space.Env, templateContext)},
			Before: executeTemplateOnExtendetShellCommand(space.Before, templateContext),
			After:  executeTemplateOnStringSlice(space.After, templateContext),
		},
		CanControllHome: space.CanControllHome,
		CanSetHome: CanSetHome{
			Home: space.CanSetHome.Home,
		},
		HasOverwrites: space.HasOverwrites,
		Executeable:   space.Executeable,
		PreFlags:      executeTemplateOnStringSlice(space.PreFlags, templateContext),
		Flags:         executeTemplateOnStringSlice(space.Flags, templateContext),
	}

	return &newSpace
}

func (space FirespaceContext) Start(cliArgs []string, dry bool) (err error) {
	defer func() {
		rec := recover()
		if iErr, ok := rec.(error); ok {
			err = iErr
		}
	}()
	space.runBeforeCommands(dry)

	cmd := space.BuildFirejailCommand(cliArgs)

	defer func() {
		space.runAfterCommands(dry)
	}()

	err = runShell(cmd, dry)
	if err != nil {
		return err
	}

	return nil

}

func (space FirespaceContext) BuildFirejailArgs(cliArgs []string) []string {

	args := []string{}

	if space.Home != "" {
		args = append(args, fmt.Sprint("--private=", space.Home))
	}
	if space.Home == "" && space.NoPrivate == false {
		args = append(args, "--private")
	}

	args = append(args, space.FirejailFlags...)

	args = append(args, space.Executeable)

	args = append(args, space.PreFlags...)

	args = append(args, cliArgs...)

	args = append(args, space.Flags...)

	return args
}

func (space FirespaceContext) BuildFirejailCommand(cliArgs []string) *exec.Cmd {

	args := space.BuildFirejailArgs(cliArgs)

	cmd := exec.Command("firejail", args...)
	cmd.Env = envMapToSlice(space.Env)

	return cmd

}

func (space FirespaceContext) runBeforeCommands(dry bool) {
	for _, beforeCommand := range space.Before {
		fields, err := shlex.Split(beforeCommand.Command)
		if err != nil {
			sugar.Panic("parsing before command", zap.Error(err))
		}

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Env = envMapToSlice(space.Env)

		err = runShell(cmd, dry)
		if err != nil {
			if !beforeCommand.AllowError {
				sugar.Panic("running before command", zap.Error(err), zap.String("command", beforeCommand.Command))
			}
		}
	}
}

func (space FirespaceContext) runAfterCommands(dry bool) {
	for _, command := range space.After {
		fields, err := shlex.Split(command)
		if err != nil {
			sugar.Panic("parsing after command", zap.Error(err))
		}

		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Env = envMapToSlice(space.Env)

		err = runShell(cmd, dry)
		if err != nil {
			sugar.Errorw("running after command", zap.Error(err), zap.String("command", command))
		}
	}
}

func executeTemplateOnMap[Key comparable](stringMap map[Key]string, data interface{}) map[Key]string {

	newMap := map[Key]string{}

	for k, v := range stringMap {

		tmplte := template.New(fmt.Sprintf("%v", k))
		tmplte, err := tmplte.Parse(v)
		if err != nil {
			sugar.Panicw("parsing template", zap.Error(err))
		}

		buf := buffer.Buffer{}

		err = tmplte.Execute(&buf, data)
		if err != nil {
			sugar.Panicw("executing template", zap.Error(err))
		}

		newMap[k] = buf.String()

	}
	return newMap

}

func executeTemplateOnStringSlice(stringSlice []string, data interface{}) []string {

	newSlice := make([]string, len(stringSlice))

	for k, v := range stringSlice {

		tmplte := template.New(fmt.Sprintf("%v", k))
		tmplte, err := tmplte.Parse(v)
		if err != nil {
			sugar.Panicw("parsing template", zap.Error(err))
		}

		buf := bytes.Buffer{}

		err = tmplte.Execute(&buf, data)
		if err != nil {
			sugar.Panicw("executing template", zap.Error(err))
		}

		newSlice[k] = buf.String()

	}
	return newSlice

}

func executeTemplateOnExtendetShellCommand(cmds []ExtendedShellCommand, data interface{}) []ExtendedShellCommand {

	newSlice := make([]ExtendedShellCommand, len(cmds))

	for k, v := range cmds {

		tmplte := template.New(fmt.Sprintf("%v", k))
		tmplte, err := tmplte.Parse(v.Command)
		if err != nil {
			sugar.Panicw("parsing template", zap.Error(err))
		}

		buf := bytes.Buffer{}

		err = tmplte.Execute(&buf, data)
		if err != nil {
			sugar.Panicw("executing template", zap.Error(err))
		}

		newSlice[k].Command = buf.String()

	}
	return newSlice

}

func runShell(cmd *exec.Cmd, dry bool) error {

	if dry {
		return runShellDry(cmd)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func runShellDry(cmd *exec.Cmd) error {

	sugar.Infow("running cmd", zap.String("cmd", cmd.String()))
	return nil
}
