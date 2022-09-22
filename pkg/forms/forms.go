package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

func (f *Form) MaxLength(field string, n int) {
	value := f.Get(field)

	if utf8.RuneCountInString(value) > n {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum %d characters).", n))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, `This field is invalid`)
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
