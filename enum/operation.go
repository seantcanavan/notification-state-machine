package enum

type Operation string

const (
	Create = "Create"
	Delete = "Delete"
	Read   = "Read"
	Update = "Update"
)

var operationToString = map[Operation]string{
	Create: "Create",
	Delete: "Delete",
	Read:   "Read",
	Update: "Update",
}

func (o Operation) String() string {
	return operationToString[o]
}

func (o Operation) Valid() bool {
	return operationToString[o] != ""
}
