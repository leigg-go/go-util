package config

import (
	"flag"
	"github.com/bmizerany/assert"
	"github.com/leigg-go/go-util/_config"
	"log"
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

func TestJsonLoaderAndReadCmdArgs(t *testing.T) {
	c := new(myConf)
	_config.NewLoader("json").MustLoadFromFile("t.json", c)

	host := flag.String("Host", "", "usage")
	port := flag.Int("port", 0, "usage")

	flagHost := "10.10.10.10"
	flagPort := "1024"
	_ = flag.Set("Host", flagHost)
	_ = flag.Set("port", flagPort)

	_config.ReadCmdArgs(func(i *flag.Flag) {
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

func TestYamlLoader(t *testing.T) {
	c := new(myConf)
	_config.NewLoader("yaml").MustLoadFromFile("t.yaml", c)
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
	_config.NewLoader("toml").MustLoadFromFile("t.toml", c)
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
