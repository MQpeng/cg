/**
 * {{.swagger.openapi}}
 * {{.swagger.info.title}}
 * {{.swagger.info.version}}
 * auto generate by cg
 */

// =======================================
{{ range $model, $value := .swagger.components.schemas -}}
{{$required := $value.required}}
export interface {{$model}} {
    {{ range $property, $propBody := $value.properties -}}
    {{$property}}
    {{- if Contains $required $property}}{{else}}?{{end -}}
    :
    {{- $propBody | SchemaToTsType }}
    {{ end }}
}
{{ end -}}

{{ range $path, $value := .swagger.paths -}}
{{$params := $path | GetRouteParams}}
// {{$path}}
{{ range $method, $value2 := $value -}}
{{ $operationId := $value2.operationId | SplitBy "_" | Distinct | FirstUpperForArray | JoinBy "" }}
export const {{$operationId}} = ({{range $i, $param := $params}}{{$param}}:string, {{end}}
    {{- if In $value2 "requestBody" -}}
    data:
    {{- if In $value2.requestBody.content `multipart/form-data`  -}}
    {{ $body := Get `multipart/form-data` $value2.requestBody.content -}}
    {{-  GetStr "$ref" $body.schema | SplitBy "/" | LastStr -}}
    {{- else if In $value2.requestBody.content `application/json` -}}
    {{ $body := Get `application/json` $value2.requestBody.content}}
    {{-  GetStr "$ref" $body.schema | SplitBy "/" | LastStr -}}
    {{- else -}}
    {{ In $value2.requestBody.content `application/json` }}
    {{- end -}}
    {{- end -}}
    ) => ({
    url: `{{$path | ParseTmpStr}}`,
    method: '{{$method}}',
    {{ if In $value2 "requestBody" -}}
    data,
    {{- end}}
})
{{ end -}}
// end {{$path}}
{{ end -}}