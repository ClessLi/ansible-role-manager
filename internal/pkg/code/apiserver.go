package code

//go:generate codegen -type=int

// arm-apiserver: inventory errors.
const (
	// ErrGroupNotFound - 404: Group not found.
	ErrGroupNotFound int = iota + 110001

	// ErrGroupAlreadyExist - 400: Group already exist.
	ErrGroupAlreadyExist
)

// arm-apiserver: roles errors.
const ()
