package main

import (
	"testing"
)

func slicesEqual(s1, s2 [][]interface{}) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		s1_1 := s1[i]
		s2_1 := s2[i]
		if len(s1_1) != len(s2_1) {
			return false
		}

		for ii, v := range s1_1 {
			if v != s2_1[ii] {
				return false
			}
		}

	}
	return true
}

type getIgnoreTestpair struct {
	args     []string
	expected [][]interface{}
}

var getIgnoreTests = []getIgnoreTestpair{
	{
		[]string{`{}`, `{}`},
		[][]interface{}{},
	},
	{
		[]string{`{}`, `{"k1": "v1"}`},
		[][]interface{}{{"k1"}},
	},
	{
		[]string{`{"k1" : "v1"}`, `{"k1" : "v1"}`},
		[][]interface{}{},
	},
	{
		[]string{`{"k1" : "v1", "k2": "v2"}`, `{"k1" : "v1", "k2": "diff"}`},
		[][]interface{}{{"k2"}},
	},
	{
		[]string{`{"k1" : "v1"}`, `{}`},
		[][]interface{}{{"k1"}},
	},
	{
		[]string{`{}`, `{"k1" : "v1"}`},
		[][]interface{}{{"k1"}},
	},
	{
		[]string{`{"k1": 1}`, `{"k1" : "1"}`},
		[][]interface{}{{"k1"}},
	},
	{
		[]string{`{"k1": 1}`, `{"k2" : 1}`},
		[][]interface{}{{"k1"}, {"k2"}},
	},
	{
		[]string{`{"k1": {"k1": "v1"}}`, `{"k1": {"k1": "v1"}}`},
		[][]interface{}{},
	},
	{
		[]string{`{"k1": {"k1": "v1"}}`, `{"k1": {"k1": "diff"}}`},
		[][]interface{}{{"k1", "k1"}},
	},
	{
		[]string{`{"k1": {"k1": "v1"}}`, `{"k1": {"k2": "v2"}}`},
		[][]interface{}{{"k1", "k1"}, {"k1", "k2"}},
	},
	{
		[]string{`{"k1": {"k1": "v1"}}`, `{"k1": {"k1": "v1", "k2": "v2"}}`},
		[][]interface{}{{"k1", "k2"}},
	},
	{
		[]string{`{"k1": [1,3]}`, `{"k1": [1,2]}`},
		[][]interface{}{{"k1", 1}},
	},
	{
		[]string{`{"k1": [2,2]}`, `{"k1": [1,2]}`},
		[][]interface{}{{"k1", 0}},
	},
	{
		[]string{`{"k1": [1,2]}`, `{"k1": [1]}`},
		[][]interface{}{{"k1", 1}},
	},
	{
		[]string{`{"k1": [2,2]}`, `{"k1": [1]}`},
		[][]interface{}{{"k1", 0}, {"k1", 1}},
	},
	{
		[]string{`{"k1": [{"k1": "v1"},1]}`, `{"k1": [{"k1": "v1"},2]}`},
		[][]interface{}{{"k1", 1}},
	},
}

func TestGetIgnores(t *testing.T) {
	for _, pair := range getIgnoreTests {

		result, err := getIgnores(pair.args[0], pair.args[1])
		if err != nil {
			t.Error("Got unexpected error ", err)
		}
		if !slicesEqual(result, pair.expected) {
			t.Error("Failed with args: ", pair.args[0], " and ", pair.args[1], ". Expected ", pair.expected, " but got ", result)
		}
	}
}

type getJSONPathTestpair struct {
	arg      []interface{}
	expected string
}

var getJSONPathTests = []getJSONPathTestpair{
	{
		[]interface{}{},
		"",
	},
	{
		[]interface{}{"a", "b", "c"},
		"[\"a\"][\"b\"][\"c\"]",
	},
	{
		[]interface{}{"a", 0, "c"},
		"[\"a\"][0][\"c\"]",
	},
}

func TestGetJSONPath(t *testing.T) {
	for _, pair := range getJSONPathTests {

		result := getJSONPath(pair.arg)
		if result != pair.expected {
			t.Error("Failed with arg: ", pair.arg)
		}
	}
}

func TestGetMessage(t *testing.T) {
	for _, pair := range getIgnoreTests {
		result := jsonDiffMessage(pair.args[0], pair.args[1])
		expected := getMessage(pair.expected)
		if result != expected {
			t.Error("Failed with args: ", pair.args[0], " and ", pair.args[1], ". Expected ", expected, " but got ", result)
		}
	}

}

func TestJSONDiffMessage(t *testing.T) {
	// Test non-JSON
	if jsonDiffMessage("not json", "not json") != "" {
		t.Error("Expected empty message")
	}

	if jsonDiffMessage(`{"k1": "v1"}`, "not json") != "" {
		t.Error("Expected empty message")
	}

	if jsonDiffMessage("not json", `{"k1": "v1"}`) != "" {
		t.Error("Expected empty message")
	}

	if jsonDiffMessage(`{"k1": "v1"}`, `{"k1": "v1"}`) != "" {
		t.Error("Expected empty message")
	}
}
