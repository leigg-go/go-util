package os

import (
	"github.com/leigg-go/go-util/_os/_file"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"runtime"
	"testing"
)

func init() {
	log.SetFlags(log.Ltime | log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stderr)
}

func Test_IsExist(t *testing.T) {
	invalidP := "/@^*(*(&%^&"
	exist := _file.IsExist(invalidP)
	//log.Print(err)
	assert.False(t, exist)

	relativeExistedP := "../util"
	ok := _file.IsExist(relativeExistedP)
	assert.True(t, ok)

	relativeNotExistedP := "../util@@@"
	ok = _file.IsExist(relativeNotExistedP)
	assert.False(t, ok)

	absoluteExistedP := "c:"
	if runtime.GOOS != "windows" {
		absoluteExistedP = "/home"
	}
	ok = _file.IsExist(absoluteExistedP)
	assert.True(t, ok)

	absoluteNotExistedP := "c:/@@@"
	if runtime.GOOS != "windows" {
		absoluteNotExistedP = "/home/@@@"
	}
	ok = _file.IsExist(absoluteNotExistedP)
	assert.False(t, ok)
}
