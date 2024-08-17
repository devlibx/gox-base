package serialization

import (
	"github.com/devlibx/gox-base/v2/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type action struct {
	Seq       int    `json:"seq"`
	Action    string `json:"action"`
	Mandatory bool   `json:"mandatory"`
}

type state struct {
	Name     string    `json:"name"`
	Actions  []*action `json:"actions"`
	Target   string    `json:"target"`
	Terminal bool      `json:"terminal"`
}

type stateMachine struct {
	States []*state `json:"states"`
}

// Test to see if we are able to read a json file
func TestReadJson(t *testing.T) {
	var sm stateMachine
	if err := ReadJson("../testdata/sm_test.json", &sm); err != nil {
		assert.Fail(t, "failed to read json file", err)
	} else {
		assert.Equal(t, len(sm.States), 3)
		assert.Equal(t, "initial", sm.States[0].Name)
		assert.Equal(t, "gox.create", sm.States[0].Actions[0].Action)
	}
}

func TestReadJsonWithBadFile(t *testing.T) {
	var sm stateMachine
	if err := ReadJson("../testdata/sm_test_does_not_exist.json", &sm); err == nil {
		assert.NotNil(t, "We expected a errors", err)
	} else {
		var e errors.Error
		if errors.As(err, &e) {
			assert.Equal(t, errors.FileOpenErrorCode, e.GetCode())
		} else {
			assert.Fail(t, "We expected the error to be type of errors.Error")
		}
	}
}

func TestToBytes(t *testing.T) {
	var sm stateMachine
	if err := ReadJson("../testdata/sm_test.json", &sm); err != nil {
		assert.Fail(t, "failed to read json file", err)
	}

	data, err := ToBytes(sm)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	first := StringifySuppressError(sm, "BAD")
	assert.NotEqual(t, "BAD", first)

	first = StringifyOrEmptyOnError(sm)
	assert.NotEqual(t, "", first)

	first = StringifyOrDefaultOnError(sm, "BAD")
	assert.NotEqual(t, "BAD", first)

	second := string(data)
	assert.Equal(t, first+"\n", second)
}

func TestToByteWithoutTags(t *testing.T) {
	type Test struct {
		User string
		Age  int
	}
	test := Test{User: "testing", Age: 10}
	data, err := ToBytes(test)
	assert.Nil(t, err, "we should see a error because struct with no tags")
	assert.NotNil(t, data)

	fromBytes := Test{}
	err = JsonBytesToObject(data, &fromBytes)
	assert.Nil(t, err, "we should see a error")
	assert.Equal(t, "testing", fromBytes.User)
	assert.Equal(t, 10, fromBytes.Age)
}
