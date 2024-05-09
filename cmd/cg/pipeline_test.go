package main

import (
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestParseTmpStr(t *testing.T) {
	assert.Equal(t, ParseTmpStr("{name}-{id}"), "${name}-${id}")
	assert.Equal(t, ParseTmpStr(""), "")
	assert.Equal(t, ParseTmpStr("{name}-{id}-{email}"), "${name}-${id}-${email}")
}

func TestGetRouteParams(t *testing.T) {
	tests := []struct {
		routeStr string
		expected []string
	}{
		{"/api/users/{userID}/repos/{repoName}", []string{"userID", "repoName"}},
		{"/api/users/{userID}", []string{"userID"}},
		{"/api/users", []string{""}},
		{"/api/", []string{""}},
		{"/api", []string{""}},
	}

	for _, tt := range tests {
		result := GetRouteParams(tt.routeStr)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("GetRouteParams(%s) = %v, expected %v", tt.routeStr, result, tt.expected)
		}
	}
}

func TestSchemaToTsType(t *testing.T) {
	type (
		TestCase struct {
			Input  interface{}
			Expect string
		}
	)

	testCases := []TestCase{
		{map[string]interface{}{"type": "object"}, "any"},
		{map[string]interface{}{"type": "array", "items": "string"}, "string[]"},
		{map[string]interface{}{"type": "integer"}, "number"},
		{map[string]interface{}{"type": "object", "anyOf": []interface{}{map[string]interface{}{"type": "string"}}}, "string | number | boolean"},
		{map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "integer"}}}, "ID"},
		{nil, ""},
	}

	for _, testCase := range testCases {
		result := SchemaToTsType(testCase.Input)
		if result != testCase.Expect {
			t.Errorf("Input: %v, Expect: %v, Result: %v", testCase.Input, testCase.Expect, result)
		}
	}
}
