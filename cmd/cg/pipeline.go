package main

import (
	"net/url"
	"strings"
)

// FirstUpper upper first char
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstUpperForArray upper
func FirstUpperForArray(s []string) []string {
	return Map(s, FirstUpper)
}

// FirstLower lower first char
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// FirstLowerForArray lower
func FirstLowerForArray(s []string) []string {
	return Map(s, FirstLower)
}

// SplitBy split str
func SplitBy(sep string, s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

// Slice slice array
func Slice(start, end int, s []string) []string {
	return s[start:end]
}

// SliceStr slice string
func SliceStr(start, end int, s string) string {
	return s[start:(len(s) - end)]
}

// ParseTmpStr parse template to template str
func ParseTmpStr(s string) string {
	return ReplaceAllStr("{", "${", s)
}

// ReplaceAllStr replace all str
func ReplaceAllStr(old, new, s string) string {
	return strings.ReplaceAll(s, old, new)
}

// Filter filter str array
func Filter[T any](slice []T, f func(T) bool) []T {
	var n []T
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

// Map map
func Map[T any](slice []T, f func(T) T) []T {
	var n []T
	for _, e := range slice {
		n = append(n, f(e))
	}
	return n
}

// ParseRouteStr parse route str
func ParseRouteStr(start int, s string) []string {
	sa := strings.Split(s, "/")
	sa = Filter(sa, func(item string) bool {
		if item == "" {
			return false
		}
		return item[:1] != "{"
	})
	return Map(sa, func(s string) string {
		ss := SplitBy("-", s)
		return JoinBy("", Map(ss, FirstUpper))
	})
}

// GetRouteParams get params from route str
func GetRouteParams(s string) []string {
	sa := strings.Split(s, "/")
	sa = Filter(sa, func(item string) bool {
		if item == "" {
			return false
		}
		return item[:1] == "{"
	})
	return Map(sa, func(s string) string {
		return s[1 : len(s)-1]
	})
}

// JoinBy join str
func JoinBy(sep string, s []string) string {
	return strings.Join(s, sep)
}

// Distinct get distinct
func Distinct(slice []string) []string {
	distinctMap := make(map[string]bool)
	distinctSlice := []string{}
	for _, str := range slice {
		if _, exists := distinctMap[str]; !exists {
			distinctMap[str] = true
			distinctSlice = append(distinctSlice, str)
		}
	}
	return distinctSlice
}

// Contains checks if a character exists in a string array.
func Contains(arr []interface{}, char string) bool {
	for _, str := range arr {
		if str == char {
			return true
		}
	}
	return false
}

// At get element by offset
func At(slice []interface{}, offset int) bool {
	length := len(slice)
    if length == 0 {
        return nil
    }
    offset = (offset % length + length) % length
    return slice[offset]
}

// In check key in map
func In(m map[string]interface{}, char string) bool {
	_, ok := m[char]
	return ok
}

// Get get value
func Get(char string, m map[string]interface{}) interface{} {
	return m[char]
}

// GetStr get str
func GetStr(char string, m map[string]interface{}) string {
	if m[char] == nil {
		return ""
	}
	return m[char].(string)
}

// LastStr return last str
func LastStr(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return s[len(s)-1]
}

// SchemaToTsType return Typescript type
func SchemaToTsType(ss interface{}) string {
	if ss == nil {
		return ""
	}
	s := ss.(map[string]interface{})
	switch s["type"] {
	case "object":
		return "any"
	case "array":
		return SchemaToTsType(s["items"]) + "[]"
	case "integer":
		return "number"
	}
	if s["anyOf"] != nil && len(s["anyOf"].([]interface{})) > 0 {
		var types []string
		for _, v := range s["anyOf"].([]interface{}) {
			types = append(types, SchemaToTsType(v))
		}
		return strings.Join(types, " | ")
	}
	if In(s, "$ref") {
		return LastStr(SplitBy("/", s["$ref"].(string)))
	}
	if s["type"] == nil {
		return ""
	}
	return s["type"].(string)
}

// Parse query string
func QueryParse(queryStr string) map[string]interface{} {
	data := make(map[string]interface{})
	if queryStr != "" {
		values, err := url.ParseQuery(queryStr)
		if err != nil {
			return data
		}
		for name, val := range values {
			if len(val) == 1 {
				data[name] = val[0]
			} else {
				data[name] = val
			}
		}
	}
	return data
}
