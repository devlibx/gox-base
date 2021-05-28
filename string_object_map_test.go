package gox

import (
	"fmt"
	"github.com/devlibx/gox-base/serialization"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringObjectMap_IntOrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue int
		Output       int
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: 1, Output: 0, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 1, Output: 1, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 11, Output: 1, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": 11}, Name: "int", DefaultValue: 0, Output: 11, err: nil},
		{TestName: "4", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 0, Output: 0, err: nil},
		{TestName: "6", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 11, Output: 11, err: nil},
		{TestName: "7", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "8", Input: StringObjectMap{"int": int32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "9", Input: StringObjectMap{"int": int64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "10", Input: StringObjectMap{"int": uint32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "11", Input: StringObjectMap{"int": uint64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "12", Input: StringObjectMap{"int": float32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "13", Input: StringObjectMap{"int": float64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "14", Input: StringObjectMap{"int": "10"}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "15", Input: StringObjectMap{"int": "10bad"}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.IntOrDefault(tt.Name, tt.DefaultValue)
			if out != tt.Output {
				t.Errorf("got %q, want %q", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_Int32OrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue int
		Output       int
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: 1, Output: 0, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 1, Output: 1, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 11, Output: 1, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": 11}, Name: "int", DefaultValue: 0, Output: 11, err: nil},
		{TestName: "4", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 0, Output: 0, err: nil},
		{TestName: "6", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 11, Output: 11, err: nil},
		{TestName: "7", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "8", Input: StringObjectMap{"int": int32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "9", Input: StringObjectMap{"int": int64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "10", Input: StringObjectMap{"int": uint32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "11", Input: StringObjectMap{"int": uint64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "12", Input: StringObjectMap{"int": float32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "13", Input: StringObjectMap{"int": float64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "14", Input: StringObjectMap{"int": "10"}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "15", Input: StringObjectMap{"int": "10bad"}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.Int32OrDefault(tt.Name, int32(tt.DefaultValue))
			if out != int32(tt.Output) {
				t.Errorf("got %q, want %q", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_Int64OrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue int
		Output       int64
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: 1, Output: 0, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 1, Output: 1, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 11, Output: 1, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": 11}, Name: "int", DefaultValue: 0, Output: 11, err: nil},
		{TestName: "4", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 0, Output: 0, err: nil},
		{TestName: "6", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 11, Output: 11, err: nil},
		{TestName: "7", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "8", Input: StringObjectMap{"int": int32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "9", Input: StringObjectMap{"int": int64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "10", Input: StringObjectMap{"int": uint32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "11", Input: StringObjectMap{"int": uint64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "12", Input: StringObjectMap{"int": float32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "13", Input: StringObjectMap{"int": float64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "14", Input: StringObjectMap{"int": "10"}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "15", Input: StringObjectMap{"int": "10bad"}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "16", Input: StringObjectMap{"int": fmt.Sprintf("%d", 1<<50)}, Name: "int", DefaultValue: 12, Output: 1 << 50, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.Int64OrDefault(tt.Name, int64(tt.DefaultValue))
			if out != int64(tt.Output) {
				t.Errorf("got %q, want %q", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_Float32OrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue int
		Output       int64
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: 1, Output: 0, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 1, Output: 1, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 11, Output: 1, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": 11}, Name: "int", DefaultValue: 0, Output: 11, err: nil},
		{TestName: "4", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 0, Output: 0, err: nil},
		{TestName: "6", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 11, Output: 11, err: nil},
		{TestName: "7", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "8", Input: StringObjectMap{"int": int32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "9", Input: StringObjectMap{"int": int64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "10", Input: StringObjectMap{"int": uint32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "11", Input: StringObjectMap{"int": uint64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "12", Input: StringObjectMap{"int": float32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "13", Input: StringObjectMap{"int": float64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "14", Input: StringObjectMap{"int": "10"}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "15", Input: StringObjectMap{"int": "10bad"}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "16", Input: StringObjectMap{"int": fmt.Sprintf("%d", 1<<50)}, Name: "int", DefaultValue: 12, Output: 1 << 50, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.Float32OrDefault(tt.Name, float32(tt.DefaultValue))
			if out != float32(tt.Output) {
				t.Errorf("got %f, want %d", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_Float64OrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue int
		Output       int64
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: 1, Output: 0, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 1, Output: 1, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: 11, Output: 1, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": 11}, Name: "int", DefaultValue: 0, Output: 11, err: nil},
		{TestName: "4", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 0, Output: 0, err: nil},
		{TestName: "6", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 11, Output: 11, err: nil},
		{TestName: "7", Input: StringObjectMap{"no_int": 11}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "8", Input: StringObjectMap{"int": int32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "9", Input: StringObjectMap{"int": int64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "10", Input: StringObjectMap{"int": uint32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "11", Input: StringObjectMap{"int": uint64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "12", Input: StringObjectMap{"int": float32(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "13", Input: StringObjectMap{"int": float64(10)}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "14", Input: StringObjectMap{"int": "10"}, Name: "int", DefaultValue: 12, Output: 10, err: nil},
		{TestName: "15", Input: StringObjectMap{"int": "10bad"}, Name: "int", DefaultValue: 12, Output: 12, err: nil},
		{TestName: "16", Input: StringObjectMap{"int": fmt.Sprintf("%d", 1<<50)}, Name: "int", DefaultValue: 12, Output: 1 << 50, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.Float64OrDefault(tt.Name, float64(tt.DefaultValue))
			if out != float64(tt.Output) {
				t.Errorf("got %f, want %d", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_BoolOrDefault(t *testing.T) {

	tests := []struct {
		TestName     string
		Input        StringObjectMap
		Name         string
		DefaultValue bool
		Output       bool
		err          error
	}{
		{TestName: "1", Input: StringObjectMap{"int": 0}, Name: "int", DefaultValue: true, Output: false, err: nil},
		{TestName: "2", Input: StringObjectMap{"int": 1}, Name: "int", DefaultValue: false, Output: true, err: nil},
		{TestName: "3", Input: StringObjectMap{"int": "true"}, Name: "int", DefaultValue: false, Output: true, err: nil},
		{TestName: "4", Input: StringObjectMap{"int": "false"}, Name: "int", DefaultValue: true, Output: false, err: nil},
		{TestName: "5", Input: StringObjectMap{"int": true}, Name: "int", DefaultValue: false, Output: true, err: nil},
		{TestName: "6", Input: StringObjectMap{"int": false}, Name: "int", DefaultValue: true, Output: false, err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			out := tt.Input.BoolOrDefault(tt.Name, tt.DefaultValue)
			if out != tt.Output {
				t.Errorf("got %t, want %t", out, tt.Output)
			}
		})
	}
}

func TestStringObjectMap_ObjectOrDefault(t *testing.T) {
	type testMe struct {
		Name string
	}
	type testMeNot struct {
		Name string
	}

	m := map[string]interface{}{"obj": &testMe{Name: "user"}}
	var sm StringObjectMap
	sm = m

	// Test with correct type
	value, ok := sm.Object("obj", &testMe{})
	assert.True(t, ok)
	assert.Equal(t, "user", value.(*testMe).Name)

	// Test with bad type
	value, ok = sm.Object("obj", &testMeNot{})
	assert.False(t, ok)

	sm = StringObjectMap{"obj": &testMe{Name: "user"}}

	// Test with correct type
	value, ok = sm.Object("obj", &testMe{})
	assert.True(t, ok)
	assert.Equal(t, "user", value.(*testMe).Name)

	// Test with bad type
	value, ok = sm.Object("obj", &testMeNot{})
	assert.False(t, ok)

	sm = StringObjectMap{"obj": testMe{Name: "user"}}

	// Test with correct type
	value, ok = sm.Object("obj", testMe{})
	assert.True(t, ok)
	assert.Equal(t, "user", value.(testMe).Name)

	// Test with bad type
	value, ok = sm.Object("obj", testMeNot{})
	assert.False(t, ok)

	sm = StringObjectMap{"no": testMe{Name: "user"}}

	// Test with correct type
	value = sm.ObjectOrDefault("obj", testMe{}, testMeNot{Name: "Not"})
	assert.Equal(t, "Not", value.(testMeNot).Name)

	// Test with bad type
	value = sm.ObjectOrDefault("obj", testMe{}, &testMeNot{Name: "Not"})
	assert.Equal(t, "Not", value.(*testMeNot).Name)
}

func TestStringObjectMap_ObjectOrDefaultWithString(t *testing.T) {
	type testMe struct {
		Name string
		Age  int
		no   int
	}
	test := testMe{
		Name: "user",
		Age:  10,
		no:   11,
	}
	sm := StringObjectMap{"obj": serialization.StringifySuppressError(test, "")}

	// Test - this must fail and give default (because we are not sending pointer)
	value := sm.ObjectOrDefault("obj", testMe{}, testMe{Name: "Not"})
	assert.Equal(t, "Not", value.(testMe).Name)

	// Test - this should pass
	value = sm.ObjectOrDefault("obj", &testMe{}, testMe{Name: "Not"})
	assert.Equal(t, "user", value.(*testMe).Name)
}

func TestStringObjectMap_TestUsage(t *testing.T) {
	m := map[string]interface{}{"int": 10, "bool": "true", "str": "Some test string"}
	var sm StringObjectMap
	sm = m
	assert.Equal(t, 10, sm.IntOrDefault("int", 0))
	assert.True(t, sm.BoolOrDefault("bool", false))
	assert.Equal(t, "Some test string", sm.StringOrDefault("str", "bad"))
}

func TestStringObjectMap_MapOrDefault(t *testing.T) {
	m := map[string]interface{}{"int": 10, "bool": "true", "str": "Some test string"}
	mapSm := StringObjectMap{"a": 10}
	sm := StringObjectMap{}
	sm["map_obj"] = m
	sm["map_obj_sm"] = mapSm

	assert.NotNil(t, sm.MapOrDefault("map_obj", nil))
	assert.NotNil(t, sm.MapOrDefault("map_obj_sm", nil))

	assert.NotNil(t, sm.MapOrDefault("map_obj_not_present", map[string]interface{}{"k": "v"}))
	assert.Equal(t, map[string]interface{}{"k": "v"}, sm.MapOrDefault("map_obj_not_present", map[string]interface{}{"k": "v"}))
	assert.NotNil(t, sm.MapOrEmpty("map_obj_not_present"))

	mapReturned := sm.MapOrEmpty("map_obj_sm")
	assert.Equal(t, 10, mapReturned["a"])
}

func TestStringObjectMap_StringObjectMapFromString(t *testing.T) {
	m := map[string]interface{}{"int": 10, "bool": "true", "str": "Some test string"}
	mapSm := StringObjectMap{"a": 10}
	sm := StringObjectMap{}
	sm["map_obj"] = m
	sm["map_obj_sm"] = mapSm

	str, err := serialization.ToBytes(sm)
	assert.NoError(t, err)

	back, err := StringObjectMapFromString(string(str))
	assert.NoError(t, err)

	str1, err := serialization.ToBytes(back)
	assert.Equal(t, str, str1)
}


func TestStringObjectMap_StringifyOrEmptyJsonOnError(t *testing.T) {
	back := serialization.StringifyOrEmptyJson(nil)
	assert.Equal(t, "{}", back)
}