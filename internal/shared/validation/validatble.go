package validation

type Validatable interface {
	Validate() error
}
