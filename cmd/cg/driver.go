package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/flosch/pongo2/v6"
	"github.com/iancoleman/strcase"
	"github.com/osteele/liquid"
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
		// strcase
		"ToSnake": strcase.ToSnake,
		"ToSnakeWithIgnore": func(ignore, s string) string {
			return strcase.ToSnakeWithIgnore(s, ignore)
		},
		"ToScreamingSnake": strcase.ToScreamingSnake,
		"ToKebab":          strcase.ToKebab,
		"ToScreamingKebab": strcase.ToScreamingKebab,
		"ToDelimited": func(d uint8, s string) string {
			return strcase.ToDelimited(s, d)
		},
		"ToScreamingDelimited": func(delimiter uint8, ignore string, screaming bool, s string) string {
			return strcase.ToScreamingDelimited(s, delimiter, ignore, screaming)
		},
		"ToCamel":      strcase.ToCamel,
		"ToLowerCamel": strcase.ToLowerCamel,
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

// LiquidTemplate is driver by liquid
func LiquidTemplate(content string, data map[string]interface{}, config *Config) io.Reader {
	engine := liquid.NewEngine()
	out, err := engine.ParseAndRenderString(content, data)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(out)
}

// Pongo2Template is driver by pongo2
func Pongo2Template(content string, data map[string]interface{}, config *Config) io.Reader {
	tpl, err := pongo2.FromString(content)
	if err != nil {
		panic(err)
	}
	out, err := tpl.Execute(data)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(out)
}
