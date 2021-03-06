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
	DefConfFileName   = "config"
	DefConfFolderName = "config"
	DefDeployCmdFlag  = "deploy"
	// config confPath by env, parse `$GO_DEPLOY_DIR/$ConfFolderName/$ConfFileName+(.suffix)`
	DefGoDeployDirEnv = "GO_DEPLOY_DIR"
)

// ConfLoader api
type ConfLoader interface {
	// File name without suffix, e.g. [config].json, should not contain `.json`
	SetDeployDir(dir string)
	SetDeployDirCmdFlag(name string)
	SetDeployDirEnvName(name string)

	SetFileName(name string)
	SetConfFolderName(name string)

	MustLoad(conf interface{})
	GetDeployDir() string
	GetConfPath() string
}

type share struct {
	deployDir      string
	fName          string
	confFolderName string
	confPath       string
	cmdFlag        string
	envName        string
}

// prior
func (s *share) SetDeployDir(dir string) {
	s.deployDir = dir
}

func (s *share) SetFileName(name string) {
	s.fName = name
}

func (s *share) SetConfFolderName(name string) {
	s.confFolderName = name
}

func (s *share) SetDeployDirCmdFlag(name string) {
	s.cmdFlag = name
}

func (s *share) SetDeployDirEnvName(name string) {
	s.envName = name
}

func (s *share) GetDeployDir() string {
	return s.deployDir
}

func (s *share) GetConfPath() string {
	return s.confPath
}

// LoadPath from cmd-line args, env in order.
func (s *share) LoadPath(suffix string) {
	defer func() {
		if s.deployDir == "" {
			panic(fmt.Sprintf("_config: no deployDir"))
		}
		// Rule of confPath concat
		s.confPath = filepath.Join(s.deployDir, s.confFolderName, s.fName+suffix)
		absP, err := filepath.Abs(s.confPath)
		_util.PanicIfErr(err, nil, fmt.Sprintf("_config: %v", err))
		s.confPath = absP
	}()

	// set default
	_util.If(s.fName == "", func() { s.fName = DefConfFileName })
	_util.If(s.confFolderName == "", func() { s.confFolderName = DefConfFolderName })
	_util.If(s.cmdFlag == "", func() { s.cmdFlag = DefDeployCmdFlag })
	_util.If(s.envName == "", func() { s.envName = DefGoDeployDirEnv })

	// called SetDeployDir()
	if s.deployDir != "" {
		return
	}

	already := false
	flag.Visit(func(i *flag.Flag) {
		if i.Name == s.cmdFlag {
			already = true
		}
	})
	if !already {
		flag.StringVar(&s.deployDir, s.cmdFlag, "", "[set by _config]: deploy dir")
	}

	flag.Parse()

	if s.deployDir != "" {
		return
	}
	s.deployDir = os.Getenv(s.envName)
}

type JsonLoader struct {
	share
}

func (l *JsonLoader) MustLoad(conf interface{}) {
	l.LoadPath(".json")

	b, err := ioutil.ReadFile(l.confPath)
	_util.PanicIfErr(err, nil, fmt.Sprintf("_config: %v", err))
	_util.PanicIfErr(json.Unmarshal(b, conf), nil)
}

type YamlLoader struct {
	share
}

func (l *YamlLoader) MustLoad(conf interface{}) {
	l.LoadPath(".yaml")

	b, err := ioutil.ReadFile(l.confPath)
	_util.PanicIfErr(err, nil, fmt.Sprintf("_config: %v", err))
	_util.PanicIfErr(yaml.Unmarshal(b, conf), nil)
}

type TomlLoader struct {
	share
}

func (l *TomlLoader) MustLoad(conf interface{}) {
	l.LoadPath(".toml")

	b, err := ioutil.ReadFile(l.confPath)
	_util.PanicIfErr(err, nil, fmt.Sprintf("_config: %v", err))
	_util.PanicIfErr(toml.Unmarshal(b, conf), nil)
}

func NewLoader(fileTyp string) ConfLoader {
	lower := strings.ToLower(fileTyp)
	switch lower {
	case "json":
		return &JsonLoader{}
	case "yaml", "yml":
		return &YamlLoader{}
	case "toml":
		return &TomlLoader{}
	}
	panic(fmt.Sprintf("_config: not support file type <%s>", lower))
}
