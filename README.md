# go-util
go utils. e.g. orm/redis...

All first-class dir name starts with `_`, e.g. `_redis`, `_config`

### Usage of `_config`

- Load config from cmd-line args/env in order.
```go
// defines config struct
type myConf struct {
	Host string
	Port int
}
var c = new(myConf)
// use default settings
loader := _config.NewLoader("json")

os.Setenv(_config.DefGoDeployDirEnv, "/deploy_dir")

// load config from `/deploy_dir/staticfile/config.json`
// staticfile is folder name, change folder name by loader.SetConfFolderName("config") = "config"
loader.MustLoad(c)

fmt.Printf("%+v", myConf)
```


### Usage of `_redis`

- use pkg `github.com/go-redis/redis`
```go
var opts = &redis.Options{
	Addr:        "192.168.40.131:63790",
	Password:    "111111",
	DB:          0,
	PoolSize:    3,
	MaxRetries:  3,
	IdleTimeout: 60 * time.Second,
}
// panic on err
_redis.MustInitDefClient(opts)

c := _redis.DefClient

_, err := c.Set("k", "v", 1*time.Second).Result()
assert.NilErr(err) // ok 
s, _ := c.Get("k").Result() // return string,err
assert.Equal(s, "v") // ok
```

### Usage of `_typ/_int`

- int utils
```go
intSlice := []int{1,2,3}
newSlice := []int{}
eachFn := func (_ int, elem int){
	newSlice = append(newSlice, elem)
}
_int.Each(intSlice, eachFn)

// `Each` supports int/int16/int32/int64/string 
```


### Usage of `_typ/_interf`

- interface utils
```go
intSlice := []int{1, 2, 3}
interfSlice := _interf.ToSliceInterface(intSlice) // get []interface{}{1,2,3}

// ToSliceInterface supports int/int16/int32/int64/string/... slice
```