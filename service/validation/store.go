package validation

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

const NilID = ID("")

type ID string

func newID() ID {
	return ID(uuid.New().String())
}

func ParseID(src string) (ID, error) {
	if src == "" {
		return NilID, fmt.Errorf("src must be a non-empty string")
	}

	return ID(src), nil
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
