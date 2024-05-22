package testsacc

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	lorem "github.com/drhodes/golorem"
	"github.com/thanhpk/randstr"
)

const (
	// prefix of generate.
	generatePrefix = "tftest-"
)

var (
	KeyValueStore = &map[string]any{}

	templateFuncs = template.FuncMap{
		"generate": func(resourceName, key string, extraOpts ...string) string {
			if len(extraOpts) == 0 {
				extraOpts = append(extraOpts, "")
			}

			randomString := generatePrefix + generateRandomString(extraOpts[0])
			(*KeyValueStore)[buildKeyValueStore(resourceName, key)] = randomString
			return returnWithQuotes(randomString)
		},
		"get": func(resourceName, key string) string {
			if v, ok := (*KeyValueStore)[buildKeyValueStore(resourceName, key)]; ok {
				if s, ok := v.(string); ok {
					return returnWithQuotes(s)
				}
			}
			return ""
		},
	}
)

// GenerateFromTemplate generates the Terraform configuration from the given template.
// The template can contain placeholders that will be replaced by the given values.
//
// Who to use:
//
//	resource "cloudavenue_catalog" "example" {
//		name             = {{ get . "name" }}
//		description      = {{ generate . "description" "longString"}}
//		delete_recursive = true
//		delete_force     = true
//	}
//
// Available functions in the template:
//   - generate: generates a random string and stores it in the key-value store. Generate accepts an optional argument that specifies the format of the random string (available formats: "shortString", "longString"). Default format is "shortString".
//   - get: returns the value of the given key from the key-value store.
func GenerateFromTemplate(resourceName, templateData string) TFData {
	// if prefix of resourceName is "data." then remove it
	resourceName = strings.TrimPrefix(resourceName, "data.")

	t, _ := template.New(resourceName).Funcs(templateFuncs).Parse(templateData)
	var tplTypes bytes.Buffer
	_ = t.Execute(&tplTypes, resourceName)

	tmpl := tplTypes.String()
	return TFData(tmpl)
}

// GetValueFromTemplate returns the value of the given key from the key-value store.
func GetValueFromTemplate(resourceName, key string) string {
	// if prefix of resourceName is "data." then remove it
	resourceName = strings.TrimPrefix(resourceName, "data.")

	if v, ok := (*KeyValueStore)[buildKeyValueStore(resourceName, key)]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// buildKeyValueStore builds the key-value store.
func buildKeyValueStore(resourceName, key string) string {
	return resourceName + "." + key
}

// generateRandomString generates a random string.
func generateRandomString(format string) string {
	// generate random string
	switch format {
	case "longString":
		return lorem.Sentence(1, 5)
	default:
		return randstr.String(16, "abcdefghijklmnopqrstuvwxyz")
	}
}

// returnWithQuotes returns the given string with quotes.
func returnWithQuotes(s string) string {
	s = strings.Trim(s, "\n")
	return fmt.Sprintf(`"%s"`, s)
}
