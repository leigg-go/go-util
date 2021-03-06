## go-util (Private)
go's common tools.  

For example:
 -  DB operations: Mysql/Redis(ok), MongoDB(todo)
 -  Go type util, e.g. _`[]string to []interface{}`_, `IntEach`, `StrEach`...
 -  Some util method like `InCollection`, `PanicIfErr`, `Must`, `RandInt(min,max)`
 -  Config Loader that load config from flag/env/config-file with a fixed mode.

All pkg name starts with `_`, e.g. `_redis`, `_config`

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

// load config from `/deploy_dir/staticfile/config/config.json`
// staticfile is folder name, change folder name by loader.SetConfFolderName("YOUR_FOLDER_NAME")
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

About More usage cases, see [test dir](https://github.com/leigg-go/go-util/tree/master/test). 