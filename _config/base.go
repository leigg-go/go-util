package _config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/leigg-go/go-util/_util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// concat with loader suffix, like config.json/config.yaml
	DefFileName       = "config"
	DefConfFolderName = "staticfile"
	DefConfCmdFlag    = "conf"
	// config path by env, parse `$GO_DEPLOY_DIR/config/`
	DefGoDeployDirEnv = "GO_DEPLOY_DIR"
)

// ConfLoader api
type ConfLoader interface {
	// File name without suffix, e.g. [config].json, should not contain `.json`
	SetFileName(name string)
	SetConfFolderName(name string)
	SetDeployDir(path string)
	SetDeployDirCmdFlag(name string)
	SetDeployDirEnvName(name string)
	MustLoad(conf interface{})
}

// ReadCmdArgs use fn to iterate command-line flags that have been set.
func ReadCmdArgs(fn func(*flag.Flag)) {
	flag.Parse()
	flag.Visit(fn)
}

type share struct {
	fName      string
	folderName string
	path       string
	cmdFlag    string
	envName    string
}

// prior
func (s *share) SetDeployDir(path string) {
	s.path = path
}

func (s *share) SetFileName(name string) {
	s.fName = name
}

func (s *share) SetConfFolderName(name string) {
	s.folderName = name
}

func (s *share) SetDeployDirCmdFlag(name string) {
	s.cmdFlag = name
}

func (s *share) SetDeployDirEnvName(name string) {
	s.envName = name
}

// LoadPath from cmd-line args, env in order.
func (s *share) LoadPath(suffix string) {
	defer func() {
		if s.path == "" {
			panic(fmt.Sprintf("_config: no path"))
		}
		s.path = filepath.Join(s.path, s.folderName, s.fName+suffix)
		absP, err := filepath.Abs(s.path)
		_util.PanicIfErr(err, nil, "_config: %v")
		s.path = absP
	}()
	// called SetDeployDir()
	if s.path != "" {
		return
	}

	// set default
	_util.If(s.fName == "", func() { s.fName = DefFileName })
	_util.If(s.folderName == "", func() { s.folderName = DefConfFolderName })
	_util.If(s.cmdFlag == "", func() { s.cmdFlag = DefConfCmdFlag })
	_util.If(s.envName == "", func() { s.envName = DefGoDeployDirEnv })

	already := false
	flag.Visit(func(i *flag.Flag) {
		if i.Name == s.cmdFlag {
			already = true
		}
	})
	if !already {
		flag.StringVar(&s.path, s.cmdFlag, "", "[set by _config]: config path")
	}

	flag.Parse()

	if s.path != "" {
		return
	}
	s.path = os.Getenv(s.envName)
}

type JsonLoader struct {
	share
}

func (l *JsonLoader) MustLoad(conf interface{}) {
	l.LoadPath(".json")

	b, err := ioutil.ReadFile(l.path)
	_util.PanicIfErr(err, nil, "_config: %v")
	_util.PanicIfErr(json.Unmarshal(b, conf), nil)
}

type YamlLoader struct {
	share
}

func (l *YamlLoader) MustLoad(conf interface{}) {
	l.LoadPath(".yaml")

	b, err := ioutil.ReadFile(l.path)
	_util.PanicIfErr(err, nil, "_config: %v")
	_util.PanicIfErr(yaml.Unmarshal(b, conf), nil)
}

type TomlLoader struct {
	share
}

func (l *TomlLoader) MustLoad(conf interface{}) {
	l.LoadPath(".toml")

	b, err := ioutil.ReadFile(l.path)
	_util.PanicIfErr(err, nil, "_config: %v")
	_util.PanicIfErr(toml.Unmarshal(b, conf), nil)
}

func NewLoader(fileTyp string) ConfLoader {
	lower := strings.ToLower(fileTyp)
	switch lower {
	case "json":
		return &JsonLoader{}
	case "yaml":
		return &YamlLoader{}
	case "toml":
		return &TomlLoader{}
	}
	panic(fmt.Sprintf("_config: not support file type <%s>", lower))
}
