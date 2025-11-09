package grapher

import (
	"fmt"
	"strconv"
)

type IntID int

func (IntID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

func (id *IntID) UnmarshalGraphQL(input any) error {
	var err error
	switch input := input.(type) {
	case string:
		if i, e := strconv.Atoi(input); e != nil {
			err = fmt.Errorf("invalid ID: %s", input)
		} else {
			*id = IntID(i)
		}
	case int32:
		*id = IntID(input)
	default:
		err = fmt.Errorf("wrong type for ID: %T", input)
	}
	return err
}

func (id IntID) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, strconv.Itoa(int(id))), nil
}

type NullIntID struct {
	Value *IntID
	Set   bool
}

func (NullIntID) ImplementsGraphQLType(name string) bool {
	return name == "ID"
}

func (s *NullIntID) UnmarshalGraphQL(input any) error {
	s.Set = true

	if input == nil {
		return nil
	}

	s.Value = new(IntID)
	return s.Value.UnmarshalGraphQL(input)
}

func (s *NullIntID) Nullable() {}
