package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/valyala/fasttemplate"
)

// FastTemplate is driver by fasttemplate
func FastTemplate(content string, data map[string]interface{}, config *Config) io.Reader{
    t := fasttemplate.New(content, config.FileStartTag, config.FileEndTag)
    return strings.NewReader(t.ExecuteString(data))
}

// TextTemplate is driver by text/template
func TextTemplate(content string, data map[string]interface{}, config *Config) io.Reader {
    tmpl, err := template.New(AppName).Delims(config.FileStartTag, config.FileEndTag).Parse(content)
    if err != nil{
        panic(err)
    }
    reader, writer := io.Pipe()  
    tmpl.Execute(writer, data)
    return reader

}