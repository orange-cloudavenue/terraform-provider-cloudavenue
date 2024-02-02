package testsacc

import (
	"regexp"
)

// * ConfigName
// String returns the string representation of the ConfigName.
func (c ConfigName) String() string {
	return string(c)
}

// IsValid checks if the given ConfigName is valid.
// Parse the string and check if format are prefix_resourceName.configName or data.prefix_resourceName.configName.
func (c ConfigName) IsValid() bool {
	re := regexp.MustCompile(`^(data\.)?cloudavenue_[a-z0-9_]+\.[a-z0-9_]+$`)
	return re.MatchString(c.String())
}

// Parse returns the resource name, the config name and if the config is a data source.
// For example, if the ConfigName is "cloudavenue_catalog.example", the resource name will be "cloudavenue_catalog" and the config name will be "example".
// If the ConfigName is "data.cloudavenue_catalog.example", the resource name will be "cloudavenue_catalog" and the config name will be "example" and the config will be a data source.
func (c ConfigName) Parse() (resourceName ResourceName, configName string, isDataSource bool) {
	x := regexp.MustCompile(`^(data\.)?cloudavenue_([a-z0-9_]+)\.([a-z0-9_]+)$`).FindStringSubmatch(c.String())
	return ResourceName(x[2]), x[3], x[1] == "data."
}
