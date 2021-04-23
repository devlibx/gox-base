package gox

import (
	"github.com/divlibx/gox-base/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestContextBuilder(t *testing.T) {
	b := util.NewRequestContextBuilder().
		Tenant("t1").
		City("india").
		Udf1("udf1").
		Build()
	assert.Equal(t, "t1", b.GetTenant())
	assert.Equal(t, "india", b.GetCity())
	assert.Equal(t, "udf1", b.GetUdf1())
	assert.Equal(t, "*", b.GetUdf2())
}
