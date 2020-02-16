package validation

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const NilID = ID("")

type ID string

func newID() ID {
	return ID(uuid.New().String())
}

func (id *ID) UnmarshalJSON(src []byte) error {
	var str string
	err := json.Unmarshal(src, &str)
	if err != nil {
		return err
	}
	v, err := ParseID(str)
	if err != nil {
		return err
	}

	*id = v

	return nil
}

func ParseID(src string) (ID, error) {
	if src == "" {
		return NilID, fmt.Errorf("src must be a non-empty string")
	}

	return ID(src), nil
}

func (id ID) IsNil() bool {
	return id == NilID
}

var Store ValidationStore

func init() {
	Store = &simpleStore{}
}

type ValidationStore interface {
	Store(*Validation) (ID, error)
	Load(ID) (*Validation, error)
}

type simpleStore struct {
	m sync.Map
}

func (s *simpleStore) Store(v *Validation) (ID, error) {
	id := newID()
	s.m.Store(id, v)

	return id, nil
}

func (s *simpleStore) Load(id ID) (*Validation, error) {
	maybeValidation, ok := s.m.Load(id)
	if !ok {
		return nil, nil
	}
	validation := maybeValidation.(*Validation)

	return validation, nil
}
