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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider"

	_ "embed"
)

const typeTemplate = `
package {{ .PackageName }}

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

{{ if not (setKeyValue "Name" .Name) }}
Error
{{ end }}

{{ template "typeStruct" . }}

{{ define "attributeType" }}
return map[string]attr.Type{
	{{- range $aN, $aD := .Attributes }}
		"{{$aN}}": {{ attrType . }},
	{{- end }}
}
{{end}}

{{ define "typeStruct" }}

{{ $Parent := (getKeyValue "ParentSubName") }}
{{ $Name := (getKeyValue "Name") }}

type {{ toUpperCamel $Name }}Model{{ if existKeyValue "SubName" }}{{singular (toUpperCamel (getKeyValue "SubName"))}}{{end}} struct {
	{{- range $aN, $aD := .Attributes }}
		{{ toUpperCamel $aN }} {{ schemaType . }}	{{ tfsdk $aN -}}
	{{ end }}
}
{{ if existKeyValue "SubName" }}{{ if not (delKeyValue "SubName")}}Error{{end}}{{end}}

{{/* if schemaType is Nested create new Type */}}
{{- range $aN, $aD := .Attributes }}
	{{- if isNestedAttribute . }}

		type {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }} []{{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $aN) }}

		// ObjectType() returns the object type for the nested object.
		func (p *{{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }}) ObjectType(ctx context.Context) types.ObjectType {
			return types.ObjectType{
				AttrTypes: p.AttrTypes(ctx),
			}
		}

		// AttrTypes() returns the attribute types for the nested object.
		func (p *{{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }}) AttrTypes(_ context.Context) map[string]attr.Type {
			{{- template "attributeType" .NestedObject -}}
		}

		// ToPlan() returns the plan representation of the nested object.
		func (p *{{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }}) ToPlan(ctx context.Context) ({{ baseTypeValue . }}, diag.Diagnostics) {
			if p == nil {
				return {{ funcNull . }}(p.ObjectType(ctx)), nil
			}

			return {{ funcFromValue . }}(ctx, p.ObjectType(ctx), p)
		}

		func (rm *{{ toUpperCamel $Name }}Model{{ if existKeyValue "SubName" }}{{ toUpperCamel (getKeyValue "SubName")}}{{end}}) {{ toUpperCamel $aN }}FromPlan(ctx context.Context) ({{ toLowerCamel $aN }} {{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }}, diags diag.Diagnostics) {
			{{ toLowerCamel $aN }} = make({{ toUpperCamel $Name }}Model{{ toUpperCamel $aN }}, 0)
			diags.Append(rm.{{ toUpperCamel $aN }}.ElementsAs(ctx, &{{ toLowerCamel $aN }}, false)...)
			if diags.HasError() {
				return
			}
		
			return {{ toLowerCamel $aN }}, diags
		}

		{{ if not (setKeyValue "SubName" $aN) }}
		Error
		{{ end }}
		{{ if not (setKeyValue "ParentSubName" $aN) }}
		Error
		{{ end }}
		{{ template "typeStruct" .NestedObject -}}

	{{ end }}
	{{/* End if isNestedAttribute */}}

	{{- if isSet . }}
		// * {{toUpperCamel $aN}}
		func (rm *{{ toUpperCamel $Name }}Model{{ singular (toUpperCamel $Parent) }}) {{toUpperCamel $aN}}FromPlan(ctx context.Context) ({{toLowerCamel $aN}} {{elementType .}}, diags diag.Diagnostics) {
			if rm.{{toUpperCamel $aN}}.IsNull() || rm.{{toUpperCamel $aN}}.IsUnknown() {
				return
			}
		
			{{toLowerCamel $aN}} = make({{elementType .}}, 0)
			diags.Append(rm.{{toUpperCamel $aN}}.ElementsAs(ctx, &{{toLowerCamel $aN}}, false)...)
			if diags.HasError() {
				return
			}
		
			return {{toLowerCamel $aN}}, diags
		}					
	{{ end }}
	{{/* End if isSet */}}
{{ end }}
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

				tD := templateData{
					Name:        strings.TrimPrefix(metadataResponse.TypeName, "_"),
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
					"isList": func(a schema.Attribute) bool {
						return IsList(reflect.TypeOf(a).String())
					},
					"isSet": func(a schema.Attribute) bool {
						return IsSet(reflect.TypeOf(a).String())
					},
					"isMap": func(a schema.Attribute) bool {
						return IsMap(reflect.TypeOf(a).String())
					},
					"schemaType": func(a schema.Attribute) string {
						return NewSchemaType(reflect.TypeOf(a).String()).ToTerraformType()
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
						return (*KeyValueStore)[k]
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
