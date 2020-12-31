package _error

import (
	"github.com/leigg-go/go-util/_util"
	"github.com/pkg/errors"
)

/*
error定义，除了敏感错误类型，其他错误都应该使用 text 尽可能表达清楚错误类型（text会直接传到前端）；
-- Core error指的是程序错误，需要立即解决
-- 注意Core error的 text 定义格式应该是 [system err NUMBER]
-- Normal error指的是可以不正确的操作导致的错误，一般不需要解决
*/

//Core error
var (
	ErrSys              = errors.New("[system err]")
	ErrMysql            = errors.New("[system err 001]")
	ErrRedis            = errors.New("[system err 002]")
	ErrSySConfig        = errors.New("[system config err]")
	ErrUnmarshal        = errors.New("[unmarshal err]")
	ErrExtractReqParams = errors.New("[extract req params err]")
	ErrParams           = errors.New("[params err]")
)

//Normal error
var (
	ErrNotAllowed       = errors.New("[operation not be allowed]")
	ErrResourceNotFound = errors.New("[resource not found]")
)

func WrapDBErr(err error) error {
	if _util.IsDBErr(err) {
		return errors.Wrap(ErrMysql, err.Error())
	}
	return nil
}

func WrapRedisErr(err error) error {
	if _util.IsRedisErr(err) {
		return errors.Wrap(ErrRedis, err.Error())
	}
	return nil
}
