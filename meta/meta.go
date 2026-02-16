package meta

type FieldDescriptor[T any] struct {
	Label     string
	Accessor  func(T) string
	Validator func(T) error     // field-specific validation
	LessThan  func(a, b T) bool // optional, use if the string values aren't reliable for sorting.. i.e  number strings
}

func NewFieldDescriptor[T any](label string, accessor func(T) string, validator func(T) error, comparator func(a, b T) bool) *FieldDescriptor[T] {

	return &FieldDescriptor[T]{
		Label:     label,
		Accessor:  accessor,
		Validator: validator,
		LessThan:  comparator,
	}
}

func (fd *FieldDescriptor[T]) StringValueFor(item T) string {
	return fd.Accessor(item)
}

type Validator[T any] interface {
	Validate(item T) error // for validating the struct as a whole (fields in conflict?)
}
