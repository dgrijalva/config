package configs_test

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/configs"
	"reflect"
	"testing"
)

type A struct {
	String string  `json:"string"`
	Float  float64 `json:"float"`
	Int    int64   `json:"int"`
}

var tests = []struct {
	name    string
	config  interface{}
	expect  interface{}
	args    string
	err     error
	options []configs.LoadOption
}{
	{
		name:   "basic",
		config: &A{"foo", 1.23, 123},
	},
	{
		name:    "flags",
		config:  &A{},
		expect:  &A{"bar", 4.56, 456},
		args:    "-string bar -float 4.56 -int 456",
		options: []configs.LoadOption{configs.WithFlags(nil)},
	},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		// Write test config JSON to buffer
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(test.config)

		// Create new item to hold parsed results
		var res interface{} = reflect.New(reflect.Indirect(reflect.ValueOf(test.config)).Type()).Interface()

		// Load config
		if test.options == nil {
			test.options = []configs.LoadOption{}
		}
		test.options = append(test.options, configs.WithReader(buf))
		err := configs.Load(res, test.options...)

		// Handle error cases
		if err != nil {
			if test.err == nil {
				t.Errorf("[%v] Unexpected error: %v", test.name, err)
			}
			if err != test.err {
				t.Errorf("[%v] Error did not meet expectations. Expected %v got %v", test.name, test.err, err)

			}
			continue
		}

		// Handle success cases
		if test.expect == nil {
			test.expect = test.config
		}
		if !reflect.DeepEqual(test.expect, res) {
			t.Errorf("[%v] Parsed config didn't match expectation. Expected %v got %T %v", test.name, test.config, res, res)
		}
	}
}
