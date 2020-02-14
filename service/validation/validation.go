package validation

import (
	"fmt"

	"github.com/nihei9/maat/service/value"
)

type validationItem struct {
	expected value.Value
	actual   []value.Value
}

type Validation struct {
	items map[string]*validationItem
}

func NewValidation() *Validation {
	return &Validation{
		items: map[string]*validationItem{},
	}
}

func (v *Validation) Expect(name string, expected value.Value) {
	v.items[name] = &validationItem{
		expected: expected,
		actual:   []value.Value{},
	}
}

func (v *Validation) Do(name string, target value.Value) (bool, error) {
	item, ok := v.items[name]
	if !ok {
		return false, fmt.Errorf("unknown validation item '%v'", name)
	}
	passed := item.expected.Test(target)

	item.actual = append(item.actual, target)

	return passed, nil
}
