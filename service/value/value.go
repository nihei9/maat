package value

import "encoding/json"

type ValueType string

func (t ValueType) String() string {
	return string(t)
}

const (
	ValueTypeText = ValueType("text")
	ValueTypeList = ValueType("list")
)

type Value interface {
	json.Marshaler
	Test(target Value) bool
	vType() ValueType
}

type TextValue struct {
	value string
}

func NewTextValue(text string) Value {
	return &TextValue{
		value: text,
	}
}

func (v *TextValue) MarshalJSON() ([]byte, error) {
	data := struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{
		Type:  v.vType().String(),
		Value: v.value,
	}
	return json.Marshal(&data)
}

func (v *TextValue) vType() ValueType {
	return ValueTypeText
}

func (v *TextValue) Test(target Value) bool {
	if target.vType() != ValueTypeText {
		return false
	}
	t := target.(*TextValue)

	return v.equal(t)
}

func (v *TextValue) equal(target *TextValue) bool {
	return v.value == target.value
}

type ListValue struct {
	value []Value
}

func NewListValue() Value {
	return &ListValue{
		value: []Value{},
	}
}

func (v *ListValue) MarshalJSON() ([]byte, error) {
	data := struct {
		Type  string  `json:"type"`
		Value []Value `json:"value"`
	}{
		Type:  v.vType().String(),
		Value: v.value,
	}
	return json.Marshal(&data)
}

func (v *ListValue) Append(e Value) {
	v.value = append(v.value, e)
}

func (v *ListValue) vType() ValueType {
	return ValueTypeList
}

func (v *ListValue) Test(target Value) bool {
	if target.vType() != ValueTypeList {
		return false
	}
	t := target.(*ListValue)

	return v.equal(t)
}

func (v *ListValue) equal(target *ListValue) bool {
	if len(v.value) != len(target.value) {
		return false
	}
	for i, e := range v.value {
		if !e.Test(target.value[i]) {
			return false
		}
	}
	return true
}
