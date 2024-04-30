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
		"ToUpper": strings.ToUpper,
        "FirstUpper": FirstUpper,
        "FirstUpperForArray": FirstUpperForArray,
		"ToLower": strings.ToLower,
        "FirstLower": FirstLower,
        "FirstLowerForArray": FirstLowerForArray,
        "JoinBy": JoinBy,
        "SplitBy": SplitBy,
        "Slice": Slice,
        "SliceStr": SliceStr,
        "ParseTmpStr": ParseTmpStr,
        "ReplaceAllStr": ReplaceAllStr,
        "ParseRouteStr": ParseRouteStr,
        "GetRouteParams": GetRouteParams,
		"Distinct": Distinct,
		"Contains": Contains,
		"In": In,
		"LastStr": LastStr,
		"Get": Get,
		"GetStr": GetStr,
		"SchemaToTsType": SchemaToTsType,
    }
	tmpl, err := template.New(AppName).Funcs(funcMap).Parse(content)
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
