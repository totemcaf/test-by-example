package jsonx

import (
	"fmt"
	"reflect"
)

const placeholderStart = '$'
const optionsSeparator = ':'
const varExpansionStart = '{'
const varExpansionEnd = '}'
const varExtractorStart = '('
const varExtractorEnd = ')'
const regexStartEnd = '/'
const regexEscapeChar = '\\'

func NewParser() Parser {
	return &parser{}
}

type parser struct {
}

func (p *parser) Parse(jsonObject interface{}) JsonX {
	value := reflect.ValueOf(jsonObject)

	x := p.parse(value)
	return x
}

func (p *parser) parse(value reflect.Value) JsonX {
	switch value.Kind() {
	case reflect.Invalid:
		return NullX
	case reflect.String:
		return p.parseString(value)
	case reflect.Int:
		return &intType{value: value.Int()}
	case reflect.Slice:
		return p.parseSlice(value)
	case reflect.Array:
		return p.parseSlice(value)
	case reflect.Map:
		return p.parseMap(value)
	case reflect.Struct:
		return p.parseStruct(value)
	case reflect.Bool:
		return &boolType{value: value.Bool()}
	case reflect.Interface:
		return p.parse(value.Elem())
	case reflect.Ptr:
		// Improve me
		if value.CanInterface() {
			if v, ok := value.Interface().(JsonX); ok {
				return v
			}
		}
		return p.parse(value.Elem())
	}
	return &nullType{} // error
}

func (p *parser) parseString(value reflect.Value) JsonX {
	return p.parseExpression(value.String())
}

func (p *parser) parseSlice(value reflect.Value) JsonX {
	result := make([]JsonX, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = p.parse(value.Index(i))
	}
	return &arrayType{values: result}
}

func (p *parser) parseStruct(value reflect.Value) JsonX {
	result := make(map[string]JsonX, value.NumField())
	typeOfS := value.Type()
	for i := 0; i < value.NumField(); i++ {
		result[typeOfS.Field(i).Name] = p.parse(value.Field(i))
	}

	return &mapType{result}
}

func (p *parser) parseMap(value reflect.Value) JsonX {
	result := make(map[string]JsonX, value.Len())
	for _, key := range value.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		result[keyStr] = p.parse(value.MapIndex(key))
	}

	return &mapType{result}
}

func (p *parser) parseExpression(s string) JsonX {
	l := len(s)
	values := make([]JsonX, 0, 1)

	pos := 0

	for pos < l {
		var element JsonX
		if s[pos] == placeholderStart {
			var err error
			element, pos, err = parseVarOrDollar(s, pos+1, l)
			if err != nil {
				panic(err)
			}
		} else {
			element, pos = parseLiteral(s, pos, l)

		}
		values = append(values, element)
	}

	if len(values) == 1 {
		return values[0]
	}
	return &concatenationType{values: values}
}

func parseLiteral(s string, pos int, l int) (JsonX, int) {
	end := pos + 1
	for end < l && s[end] != placeholderStart {
		end++
	}
	return &stringType{value: s[pos:end]}, end
}

func NewInvalidExpression(s string, pos int) error {
	return fmt.Errorf("invalid expression: %s at %d", s, pos)
}

func parseVarOrDollar(s string, pos, l int) (JsonX, int, error) {
	if pos >= l {
		return nil, pos, NewInvalidExpression(s, pos)
	}

	switch s[pos] {
	case varExpansionStart:
		return parseVarExtended(s, pos, l)
	case varExtractorStart:
		return parseExtractor(s, pos, l)
	case placeholderStart:
		return &stringType{value: "$"}, pos + 1, nil
	default:
		return parseSimpleVar(s, pos, l)

	}
}

func parseSimpleVar(s string, pos int, l int) (JsonX, int, error) {
	end := pos + 1
	for end < l && isNameChar(s[end]) {
		end++
	}
	return &varExpansionType{varName: s[pos:end]}, end, nil
}

func parseVarExtended(s string, pos int, l int) (JsonX, int, error) {
	end := pos + 1
	for end < l && s[end] != varExpansionEnd && s[end] != optionsSeparator {
		end++
	}
	if end >= l {
		return nil, end, NewInvalidExpression(s, pos)
	}

	name := s[pos+1 : end]

	if s[end] == optionsSeparator {
		return parseVarWithGenerator(name, s, end+1, l)
	}

	return &varExpansionType{varName: name}, end + 1, nil
}

func parseExtractor(s string, pos int, l int) (JsonX, int, error) {
	end := pos + 1
	for end < l && s[end] != varExtractorEnd {
		end++
	}
	if end >= l {
		return nil, end, NewInvalidExpression(s, pos)
	}

	name := s[pos+1 : end]

	return &extractorType{varName: name}, end + 1, nil // TODO
}

func parseVarWithGenerator(name string, s string, pos int, l int) (JsonX, int, error) {
	end := pos
	for end < l && s[end] != varExpansionEnd && s[end] != optionsSeparator {
		end++
	}

	if end >= l {
		return nil, end, NewInvalidExpression(s, pos)
	}

	typeName := s[pos:end]

	var config string

	if s[end] == optionsSeparator {
		var err error
		config, end, err = parseGeneratorConfig(s, end+1, l)
		if err != nil {
			return nil, end, err
		}
	}

	if jsonX := NewRandomValue(name, typeName, config); jsonX != nil {
		return jsonX, end + 1, nil
	} else {
		return nil, end, NewInvalidExpression(s, pos)
	}
}

func parseGeneratorConfig(s string, pos int, l int) (string, int, error) {
	if pos >= l {
		return "", pos, nil
	}

	if s[pos] == regexStartEnd {
		return parseRegex(s, pos, l)
	}

	end := pos
	for end < l && s[end] != varExpansionEnd {
		end++
	}
	return s[pos:end], end, nil
}

// skipTo

func parseRegex(s string, pos int, l int) (string, int, error) {
	end := pos + 1
	for end < l && s[end] != regexStartEnd {
		if s[end] == regexEscapeChar && end+1 < l {
			end++
		}
		end++
	}

	if end >= l {
		return "", end, NewInvalidExpression(s, pos)
	}

	return s[pos+1 : end], end + 1, nil
}

func isNameChar(u uint8) bool {
	return u >= 'a' && u <= 'z' || u >= 'A' && u <= 'Z' || u >= '0' && u <= '9' || u == '_'
}
