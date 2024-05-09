package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/valyala/fasttemplate"
)

// FastTemplate is driver by fasttemplate
func FastTemplate(content string, data map[string]interface{}, config *Config) io.Reader {
	t := fasttemplate.New(content, config.FileStartTag, config.FileEndTag)
	return strings.NewReader(t.ExecuteString(data))
}

// TextTemplate is driver by text/template
func TextTemplate(content string, data map[string]interface{}, config *Config) io.Reader {
	funcMap := template.FuncMap{
		// String operator pipeline
		"ToUpper":       strings.ToUpper,
		"SplitBy":       SplitBy,
		"FirstUpper":    FirstUpper,
		"ToLower":       strings.ToLower,
		"FirstLower":    FirstLower,
		"SliceStr":      SliceStr,
		"ParseTmpStr":   ParseTmpStr,
		"ReplaceAllStr": ReplaceAllStr,
		"ParseRouteStr": ParseRouteStr,
		// Array operator pipeline
		"FirstUpperForArray": FirstUpperForArray,
		"FirstLowerForArray": FirstLowerForArray,
		"JoinBy":             JoinBy,
		"Slice":              Slice,
		"LastStr":            LastStr,
		"Contains":           Contains,
		"Distinct":           Distinct,
		// Collection operator pipeline
		"In":     In,
		"Get":    Get,
		"GetStr": GetStr,
		// Instrumental operator pipeline
		"GetRouteParams": GetRouteParams,
		"QueryParse":     QueryParse,
		"SchemaToTsType": SchemaToTsType,
	}
	tmpl, err := template.New(AppName).Funcs(funcMap).Delims(config.FileStartTag, config.FileEndTag).Parse(content)
	if err != nil {
		panic(err)
	}
	reader, writer := io.Pipe()
	go func() {
		defer writer.Close()
		err = tmpl.Execute(writer, data)
		if err != nil {
			panic(err)
		}
	}()
	return reader
}
