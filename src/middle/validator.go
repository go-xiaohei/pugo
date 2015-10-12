package middle

import (
	"fmt"
	"github.com/tango-contrib/binding"
)

var (
	ValidatorRenderName = "ValidateError"

	_ IValidator = (*baseValidator)(nil)
	_ IValidator = (*Validator)(nil)
)

type IValidator interface {
	Validate(v interface{}, controller interface{}) error
}

type baseValidator struct {
	binding.Binder
	onFail func(errors binding.Errors, controller interface{})
}

func (v *baseValidator) setOnFail() {
	if v.onFail != nil {
		return
	}
	v.onFail = func(errors binding.Errors, controller interface{}) {
		if errors.Len() == 0 {
			return
		}

		// assign to render
		if render, ok := controller.(ITheme); ok {
			render.Assign(ValidatorRenderName, fmt.Sprintf("%s %s", errors[0].FieldNames[0], errors[0].Message))
		}
	}
}

func (bv *baseValidator) Validate(value interface{}, controller interface{}) error {
	bv.setOnFail()
	if errors := bv.Binder.Bind(value); errors.Len() > 0 {
		bv.onFail(errors, controller)
		return fmt.Errorf("%s %s", errors[0].FieldNames[0], errors[0].Message)
	}
	return nil
}

type Validator struct {
	baseValidator
}
