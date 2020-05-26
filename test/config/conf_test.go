package config

import (
	"flag"
	"github.com/bmizerany/assert"
	"github.com/leigg-go/go-util/_config"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type myConf struct {
	Host   string
	Port   int
	Labels []string
	Nested struct {
		Host   string
		Port   int
		Labels []string
	}
}

var deployDir string

func init() {
	p, _ := filepath.Abs(".")
	deployDir = filepath.Join(filepath.Dir(filepath.Dir(p)), "deploy")
	log.Print(deployDir)
}

// Test jsonLoader and cmd-line parsing
func TestJsonLoaderAndReadCmdArgs(t *testing.T) {

	c := new(myConf)
	// 1. cmd-line args are first, set cmd flag to test
	// but this step is not a required, because pkg did it for help us,
	// just for testing.
	path := flag.String(_config.DefDeployCmdFlag, "", "")
	_ = flag.Set(_config.DefDeployCmdFlag, deployDir)

	// 2. get loader
	loader := _config.NewLoader("json")

	// 3. set filename
	loader.SetFileName("t") // no suffix, default `config`
	loader.SetDeployDir(*path)

	// so it will load config file from `$deployDir/config/t.json` now.
	loader.MustLoad(c)

	// do other args replace from cmd-line args
	host := flag.String("Host", "", "usage")
	port := flag.Int("port", 0, "usage")

	flagHost := "10.10.10.10"
	flagPort := "1024"
	_ = flag.Set("Host", flagHost)
	_ = flag.Set("port", flagPort)

	flag.Parse()
	flag.Visit(func(i *flag.Flag) {
		switch i.Name {
		case "Host":
			c.Host = *host
		case "port":
			c.Port = *port
		}
	})

	should := myConf{
		Host: *host, Port: *port, Labels: []string{"a", "b", "c"}, Nested: struct {
			Host   string
			Port   int
			Labels []string
		}{Host: "2.2.2.2", Port: 6668, Labels: []string{"c", "b", "a"}},
	}
	log.Printf("%+v", *c)
	assert.Equal(t, should, *c)
}

// Test YamlLoader and env parsing
func TestYamlLoader(t *testing.T) {
	c := new(myConf)

	// set env manually.
	_ = os.Setenv(_config.DefGoDeployDirEnv, deployDir)

	loader := _config.NewLoader("yaml")
	loader.SetFileName("t")

	// so it will load config file from `$deployDir/config/t.yaml` now.
	loader.MustLoad(c)

	should := myConf{
		Host: "1.1.1.1", Port: 666, Labels: []string{"a", "b", "c"}, Nested: struct {
			Host   string
			Port   int
			Labels []string
		}{Host: "2.2.2.2", Port: 6668, Labels: []string{"c", "b", "a"}},
	}
	log.Printf("%+v", *c)
	assert.Equal(t, should, *c)
}

type tomlConf struct {
	Title string
	Mysql struct {
		Host   string
		Port   int
		Labels []string
		Nested struct {
			Host   string
			Port   int
			Labels []string
		}
	}
	RedisSvrs map[string]struct {
		Host string
		Port int
	}
}

func TestTomlLoader(t *testing.T) {
	c := new(tomlConf)
	loader := _config.NewLoader("toml")
	loader.SetFileName("t")
	loader.SetDeployDir(deployDir)

	// so it will load config file from `$deployDir/config/t.toml` now.
	loader.MustLoad(c)

	should := tomlConf{
		Title: "toml example",
		Mysql: struct {
			Host   string
			Port   int
			Labels []string
			Nested struct {
				Host   string
				Port   int
				Labels []string
			}
		}{Host: "1.1.1.1", Port: 666, Labels: []string{"a", "b", "c"},
			// nested struct
			Nested: struct {
				Host   string
				Port   int
				Labels []string
			}{Host: "1.1.1.1", Port: 777, Labels: []string{"a", "b", "c"}},
		},
		// map
		RedisSvrs: map[string]struct {
			Host string
			Port int
		}{"svr1": {"3.3.3.3", 888}, "svr2": {"3.3.3.4", 888}},
	}
	log.Printf("%+v", *c)
	assert.Equal(t, should, *c)
}
