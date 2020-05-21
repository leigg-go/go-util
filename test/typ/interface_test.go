package typ

import (
	"fmt"
	"github.com/bmizerany/assert"
	"github.com/leigg-go/go-util/_typ/_interf"
	"testing"
)

func TestToSliceInterface(t *testing.T) {
	testData := []struct {
		before interface{}
		after  []interface{}
	}{
		{[]int{1, 2, 3}, []interface{}{1, 2, 3}},
		{[]float64{1.0, 2.0, 3.0}, []interface{}{1.0, 2.0, 3.0}},
		{[]error{fmt.Errorf("1"), fmt.Errorf("2")}, []interface{}{fmt.Errorf("1"), fmt.Errorf("2")}},
	}
	for _, d := range testData {
		assert.Equal(t, _interf.ToSliceInterface(d.before), d.after)
	}

}
