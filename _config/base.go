package _config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/leigg-go/go-util/_util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ConfLoader api
type ConfLoader interface {
	MustLoadFromFile(path string, conf interface{})
}

// ReadCmdArgs use fn to iterate command-line flags that have been set.
func ReadCmdArgs(fn func(*flag.Flag)) {
	flag.Parse()
	flag.Visit(fn)
}

type JsonLoader struct{}

func (j *JsonLoader) MustLoadFromFile(path string, conf interface{}) {
	absPath, err := filepath.Abs(path)
	_util.PanicIfErr(err, nil, "_config: %v")

	b, err := ioutil.ReadFile(absPath)
	_util.PanicIfErr(err, nil, "_config: %v")
	_util.PanicIfErr(json.Unmarshal(b, conf), nil)
}

type YamlLoader struct{}

func (y *YamlLoader) MustLoadFromFile(path string, conf interface{}) {
	absPath, err := filepath.Abs(path)
	_util.PanicIfErr(err, nil, "_config: %v")

	b, err := ioutil.ReadFile(absPath)
	_util.PanicIfErr(err, nil, "_config: %v")
	_util.PanicIfErr(yaml.Unmarshal(b, conf), nil)
}

type TomlLoader struct{}

func (t *TomlLoader) MustLoadFromFile(path string, conf interface{}) {
	absPath, err := filepath.Abs(path)
	_util.PanicIfErr(err, nil, "_config: %v")

	b, err := ioutil.ReadFile(absPath)
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
