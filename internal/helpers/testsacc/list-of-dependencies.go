package testsacc

import "log"

// * ListOfDependencies
// Append appends the given dependency to the list of dependencies.
func (l *ListOfDependencies) Append(configName ConfigName, tfData TFData) {
	if configName.IsValid() && !l.Exists(configName) {
		log.Default().Printf("Adding dependency %s", configName)
		(*l)[configName] = tfData
	} else {
		log.Default().Printf("Dependency %s already exists", configName)
	}
}

// Exists checks if the given dependency exists in the list of dependencies.
func (l *ListOfDependencies) Exists(configName ConfigName) bool {
	if l == nil {
		return false
	}

	if len(*l) == 0 {
		return false
	}

	_, ok := (*l)[configName]
	return ok
}

// Get returns the list of dependencies as a slice of string.
func (l *ListOfDependencies) Get() []string {
	x := make([]string, 0)
	for _, v := range *l {
		x = append(x, v.String())
	}
	return x
}

// ToTFData returns the list of dependencies as a TFData (string).
func (l *ListOfDependencies) ToTFData() TFData {
	t := TFData("")
	for _, v := range *l {
		t.Append(v)
	}
	return t
}
