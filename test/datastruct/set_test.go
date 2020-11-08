package datastruct

import (
	"fmt"
	"github.com/leigg-go/go-util/_datastruct"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {

	set := _datastruct.NewConcurrencyMap()
	for i := 0; i < 10; i++ {
		assert.True(t, set.AddElem(i))
	}

	for i := 0; i < 10; i++ {
		assert.True(t, set.Contains(i))
		assert.True(t, set.RemoveElem(i))
		assert.False(t, set.Contains(i))
	}

	for i := 0; i < 10; i++ {
		assert.True(t, set.AddElem(i))
	}
	iterate := func(seq int, elem interface{}) bool {
		fmt.Println(seq, elem)
		return true
	}
	set.Range(iterate)

	st := _datastruct.NewSet()
	st.Add(1, 23, 4, 1)
	st.Remove(23)
	assert.Equal(t, st.Size(), 2)
	assert.True(t, st.Contains(1) && st.Contains(4))
	assert.False(t, st.Contains(23))
}
