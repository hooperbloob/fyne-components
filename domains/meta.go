package domains

type FieldDescriptor[T any] struct {
	Label     string
	ValueType string
	Accessor  func(T) string
	Validator func(T) error
}
