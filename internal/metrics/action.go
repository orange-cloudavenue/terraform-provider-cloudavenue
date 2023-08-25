package metrics

type Action string

const (
	Create Action = "Create"
	Read   Action = "Read"
	Update Action = "Update"
	Delete Action = "Delete"
	Import Action = "Import"
)

// String returns the string representation of the action.
func (a Action) String() string {
	return string(a)
}
