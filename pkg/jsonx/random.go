package jsonx

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
)

type RandomType int

const (
	RandomString RandomType = iota
	RandomName
	RandomEmail
	RandomPhone
	RandomAddress
	RandomCompanyName
	RandomRegex
)

var types = map[string]RandomType{
	"random.string":      RandomString,
	"random.name":        RandomName,
	"random.email":       RandomEmail,
	"random.phone":       RandomPhone,
	"random.address":     RandomAddress,
	"random.companyName": RandomCompanyName,
	"random.regex":       RandomRegex,
}

const (
	RandomValue Type = "randomValue"
)

type randomValueType struct {
	_type  RandomType
	config string
	name   string
}

func NewRandomValue(name, typeName, config string) *randomValueType {
	_type, found := types[typeName]

	if !found {
		return nil
	}
	return &randomValueType{name: name, _type: _type, config: config}
}

func (n *randomValueType) Equals(other JsonX) bool {
	panic("implement randomValue.Equals")
}

func (n *randomValueType) Diff(_ Context, actual JsonX) Differences {
	panic("implement randomValue.Diff")
}

func (n *randomValueType) Eval(context Context) JsonX {
	evaluated := n.eval()

	if n.name != "" {
		context.Set(n.name, evaluated)
	}

	return evaluated
}

func (n *randomValueType) eval() JsonX {
	switch n._type {
	case RandomString:
		return &stringType{value: gofakeit.LetterN(6)}

	case RandomName:
		return &stringType{value: gofakeit.Name()}

	case RandomEmail:
		return &stringType{value: gofakeit.Email()}

	case RandomPhone:
		return &stringType{value: gofakeit.Phone()}

	case RandomAddress:
		address := gofakeit.Address()
		return &stringType{value: fmt.Sprintf("%s, %s, %s, %s, %s", address.Address, address.City, address.State, address.Zip, address.Country)}

	case RandomCompanyName:
		return &stringType{value: gofakeit.Company()}

	case RandomRegex:
		return &stringType{value: gofakeit.Regex(n.config)}

	default:
		panic("random type unknown or not implemented: " + n.typeName())
	}
}

func (n *randomValueType) typeName() string {
	for k, v := range types {
		if v == n._type {
			return k
		}
	}
	return fmt.Sprintf("unknown %d", n._type)
}
func (n *randomValueType) Type() Type {
	return RandomValue
}

func (n *randomValueType) String() string {
	return fmt.Sprintf("${%s:%s:%s}", n.name, n.typeName(), n.config)
}
