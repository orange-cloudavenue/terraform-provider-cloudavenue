package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/kr/pretty"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"

	_ "embed"
)

const typeTemplate = `
package {{ .PackageName }}

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

{{ if not (setKeyValue "Name" .Name) }}
Error
{{ end }}

{{ template "typeStruct" . }}


{{ define "typeStruct" }}

{{ $Parent := "" }}
{{ if existKeyValue "ParentSubName" }}
	{{ $Parent = (getKeyValue "ParentSubName") }}
{{ end }}
{{ $Name := (getKeyValue "Name") }}


{{ if existKeyValue "SubName" }}
// * {{ toUpperCamel (getKeyValue "SubName") -}}
{{ end }}
type {{ toUpperCamel $Name }}Model{{ if existKeyValue "SubName" }}{{singular (toUpperCamel (getKeyValue "SubName"))}}{{end}} struct {
	{{- range $aN, $aD := .Attributes }}
		{{ toUpperCamel $aN }} {{ terraformValue . }}	{{ tfsdk $aN -}}
	{{ end }}
}
{{ if existKeyValue "SubName" }}{{ if not (delKeyValue "SubName")}}Error{{end}}{{end}}

{{/* if schemaType is Nested create new Type */}}
{{- range $aN, $aD := .Attributes }}
	{{- if isNestedAttribute . }}
		// * {{ toUpperCamel $aN -}} 
		{{ if isMap . }}
			type {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }} map[string]{{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $aN) }}
		{{ else }}
			type {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }} []{{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $aN) }}
		{{ end }}
		{{ if not (setKeyValue "SubName" $aN) }}
		Error
		{{ end }}
		{{ if not (setKeyValue "ParentSubName" $aN) }}
		Error
		{{ end }}
		{{ template "typeStruct" .NestedObject -}}

	{{ end }}
	{{/* End if isNestedAttribute */}}

	{{- if isSingle . }}

		{{ if not (setKeyValue "SubName" $aN) }}
		Error
		{{ end }}
		{{ if not (setKeyValue "ParentSubName" $aN) }}
		Error
		{{ end }}
		{{ template "typeStruct" . -}}

	{{ end }}
	{{/* End if isSingle */}}

	{{ if isArray . }}
		{{ if isMap . }}
			type {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }} map[string]{{ terraformValue . }}
		{{ else }}
			type {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }} []{{ terraformValue . }}
		{{ end }}
	{{ end }}
	{{/* End if isArray */}}



{{ end }}
{{ end }}


func New{{ toUpperCamel .Name }}(t any) *{{ toUpperCamel .Name }}Model {
	switch x := t.(type) {
	case tfsdk.State: //nolint:dupl
		return {{ template "structNewFunc" . }}
	case tfsdk.Plan: //nolint:dupl
		return {{ template "structNewFunc" . }}
	case tfsdk.Config: //nolint:dupl
		return {{ template "structNewFunc" . }}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *{{ toUpperCamel .Name }}Model) Copy() *{{ toUpperCamel .Name }}Model {
	x := &{{ toUpperCamel .Name }}Model{}
	utils.ModelCopy(rm, x)
	return x
}

{{- range $aN, $aD := .Attributes }}
	{{- $Name := (getKeyValue "Name") -}}

	{{ if isNestedOrArrayAttribute . }}
		// Get{{ toUpperCamel $aN }} returns the value of the {{ toUpperCamel $aN }} field.
		func (rm *{{ toUpperCamel $Name }}Model) Get{{ toUpperCamel $aN }}(ctx context.Context) (values {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN}}, diags diag.Diagnostics) {
			values = make({{ toUpperCamel $Name }}Model{{ toUpperCamel $aN}}, 0)
			d := rm.{{ toUpperCamel $aN }}.Get(ctx, &values, false)
			return values, d
		}
	{{ end }}

	{{ if isSingle . }}
		// Get{{ toUpperCamel $aN }} returns the value of the {{ toUpperCamel $aN }} field.
		func (rm *{{ toUpperCamel $Name }}Model) Get{{ toUpperCamel $aN }}(ctx context.Context) (values {{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $aN)}}, diags diag.Diagnostics) {
			values = {{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $aN)}}{}
			d := rm.{{ toUpperCamel $aN }}.Get(ctx, &values, basetypes.ObjectAsOptions{})
			return values, d
		}
	{{ end }}
{{ end }}



{{ define "structNewFunc" }}
	{{- $Name := (getKeyValue "Name") -}}
	&{{ toUpperCamel $Name }}Model{
	{{ range $aN, $aD := .Attributes }}
		{{- if isNestedOrArrayAttribute . -}}
			{{ toUpperCamel $aN }}: {{ funcNullOrUnkown . }}(x.Schema.GetAttributes()["{{ toSnakeCase $aN }}"].GetType().({{ terraformType . }}).ElementType()),
		{{- else if isSingle . -}}
			{{ toUpperCamel $aN }}: {{ funcNullOrUnkown . }}(x.Schema.GetAttributes()["{{ toSnakeCase $aN }}"].GetType().({{ terraformType . }}).AttributeTypes()),
		{{- else -}}
			{{ toUpperCamel $aN }}: {{ funcNullOrUnkown . }}(),
		{{ end }}
	{{ end -}}
	}

{{ end }}


`

type templateData struct {
	Name        string
	PackageName string
	Attributes  map[string]schema.Attribute
}

var KeyValueStore *map[string]any

func main() {
	resourceName := new(string)
	isResource := new(bool)
	isDataSource := new(bool)
	filePath := new(string)

	flag.StringVar(filePath, "file", "", "file path")
	flag.StringVar(resourceName, "resource", "", "resource name")
	flag.BoolVar(isResource, "is-resource", false, "is resource")
	flag.BoolVar(isDataSource, "is-data-source", false, "is data source")
	flag.Parse()

	if *resourceName == "" || (!*isResource && !*isDataSource) || *filePath == "" {
		flag.PrintDefaults()
		return
	}

	log.Info().Msgf("Looking for resource %s", *resourceName)

	ctx := context.Background()

	cavP := provider.New(provider.VCDVersion)

	for _, res := range cavP().Resources(ctx) {
		metadataResponse := &resource.MetadataResponse{}
		res().Metadata(ctx, resource.MetadataRequest{}, metadataResponse)

		log.Info().Msgf("Find resource %s", metadataResponse.TypeName)
		if "cloudavenue"+metadataResponse.TypeName == *resourceName {
			log.Info().Msgf("Found resource %s", *resourceName)

			if *isResource {
				resp := &resource.SchemaResponse{}
				res().Schema(ctx, resource.SchemaRequest{}, resp)

				// read file
				content, err := ioutil.ReadFile(*filePath)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to read file")
				}
				// Do something with the content
				packageName := ""
				for _, line := range strings.Split(string(content), "\n") {
					if !strings.Contains(line, "//") {
						words := strings.Split(line, " ")
						// Get package name
						packageName = words[len(words)-1]
						break
					}
				}

				// metadataResponse.TypeName = "_demo_superschema_supertypes"
				// if metadataResponse.TypeName contains two or more underscores, remove the two first underscores

				metadataResponse.TypeName = strings.TrimPrefix(metadataResponse.TypeName, "_")

				if strings.Count(metadataResponse.TypeName, "_") > 1 {
					first := strings.Split(metadataResponse.TypeName, "_")[0]
					metadataResponse.TypeName = strings.TrimPrefix(metadataResponse.TypeName, first+"_")
				}

				tD := templateData{
					Name:        metadataResponse.TypeName,
					PackageName: packageName,
					Attributes:  resp.Schema.Attributes,
				}

				KeyValueStore = &map[string]any{}

				templateFuncs := template.FuncMap{
					"toLowerCamel": func(s string) string {
						return strcase.ToLowerCamel(s)
					},
					"toUpperCamel": func(s string) string {
						return strcase.ToCamel(s)
					},
					"toSnakeCase": func(s string) string {
						return strcase.ToSnake(s)
					},
					"isNestedAttribute": func(a schema.Attribute) bool {
						return IsNested(reflect.TypeOf(a).String())
					},
					"isNestedOrArrayAttribute": func(a schema.Attribute) bool {
						return IsNestedOrArray(reflect.TypeOf(a).String())
					},
					"isArray": func(a schema.Attribute) bool {
						return IsArray(reflect.TypeOf(a).String())
					},
					"isList": func(a schema.Attribute) bool {
						return IsList(reflect.TypeOf(a).String())
					},
					"isSet": func(a schema.Attribute) bool {
						return IsSet(reflect.TypeOf(a).String())
					},
					"isMap": func(a schema.Attribute) bool {
						return IsMap(reflect.TypeOf(a).String())
					},
					"isSingle": func(a schema.Attribute) bool {
						return IsSingle(reflect.TypeOf(a).String())
					},
					"terraformType": func(a schema.Attribute) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToTerraformType()
					},
					"terraformValue": func(a schema.Attribute) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToTerraformValue()
					},
					"elementType": func(a any) string {
						return NewElementType(a).ToTerraformType()
					},
					"baseTypeValue": func(a any) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToBaseTypeValue()
					},
					"attrType": func(a any) string {
						return NewAttributeType(a)
					},
					"funcNull": func(a any) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToFuncNull()
					},
					"funcUnkown": func(a any) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToFuncUnkown()
					},
					"funcNullOrUnkown": func(a schema.Attribute) string {
						if a.IsComputed() {
							return NewSchemaType(reflect.TypeOf(a).String()).ToFuncUnkown()
						}
						return NewSchemaType(reflect.TypeOf(a).String()).ToFuncNull()
					},
					"funcFromValue": func(a any) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToValueFrom()
					},
					"tfsdk": func(a string) string {
						return "`tfsdk:" + "\"" + a + "\"" + "`"
					},
					"setKeyValue": func(k string, v any) bool {
						(*KeyValueStore)[k] = v
						return true
					},
					"getKeyValue": func(k string) any {
						if v, ok := (*KeyValueStore)[k]; ok {
							return v
						}
						return nil
					},
					"delKeyValue": func(k string) bool {
						delete(*KeyValueStore, k)
						return true
					},
					"existKeyValue": func(k string) bool {
						_, ok := (*KeyValueStore)[k]
						return ok
					},
					"singular": func(s string) string {
						return strings.TrimSuffix(s, "s")
					},
				}

				tmpl, err := template.New("template").Funcs(templateFuncs).Parse(typeTemplate)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to parse template")
				}

				var tplTypes bytes.Buffer

				if err := tmpl.Execute(&tplTypes, tD); err != nil {
					log.Fatal().Err(err).Msg("Failed to execute template")
					return
				}

				file := strings.TrimSuffix(*filePath, ".go") + "_types.go"

				if err := os.WriteFile(file, tplTypes.Bytes(), 0o600); err != nil {
					log.Fatal().Err(err).Msg("Failed to write file")
					return
				}

				// format go file
				cmd := exec.Command("gofmt", "-s", "-w", file)
				if err := cmd.Run(); err != nil {
					log.Fatal().Err(err).Msg("Failed to format file")
					return
				}

				pretty.Print(resp.Schema)
			}
			return
		}
	}

	log.Error().Msgf("Resource %s not found", *resourceName)
}
