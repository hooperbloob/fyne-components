package meta

type FieldDescriptor[T any] struct {
	Label     string
	Accessor  func(T) string
	Validator func(T) error     // field-specific validation
	lessThan  func(a, b T) bool // optional, use if the string values aren't reliable for sorting.. i.e  numbers, dates, etc
}

func NewFieldDescriptor[T any](label string, accessor func(T) string, validator func(T) error, lessThan func(a, b T) bool) *FieldDescriptor[T] {

	return &FieldDescriptor[T]{
		Label:     label,
		Accessor:  accessor,
		Validator: validator,
		lessThan:  lessThan,
	}
}

func (fd *FieldDescriptor[T]) StringValueFor(item T) string {
	return fd.Accessor(item)
}

func (fd *FieldDescriptor[T]) LessThan() func(a, b T) bool {

	if fd.lessThan != nil {
		return fd.lessThan
	}

	// return one that uses string values
	return func(a, b T) bool {
		return fd.Accessor(a) < fd.Accessor(b)
	}
}

type Validator[T any] interface {
	Validate(item T) error // for validating the struct as a whole (fields in conflict?)
}
