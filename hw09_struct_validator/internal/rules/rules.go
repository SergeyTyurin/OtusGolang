package rules

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	errGetRule  = errors.New("incorrect type or rule's name")
	errNilValue = errors.New("value is nil")

	strTypeError       = "type of value must be %w"
	strParsingError    = "parsing rule error \"%w\""
	strValidationError = "rule: \"%w\", value: \"%w\""

	avaliableRules = map[reflect.Kind][]string{
		reflect.Int:    {"min", "max", "in"},
		reflect.String: {"len", "regexp", "in"},
	}
)

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type TypeError struct {
	Err error
}

func (e TypeError) Error() string {
	return e.Err.Error()
}

type ValidationRule interface {
	CheckValue(rule string, value interface{}) error
}

func getStrValue(value reflect.Value) string {
	return fmt.Sprint(value)
}

func getRuleByName(name string) ValidationRule {
	switch name {
	case "min":
		return &MinRule{}
	case "max":
		return &MaxRule{}
	case "len":
		return &LengthRule{}
	case "regexp":
		return &RegexRule{}
	case "in":
		return &InRule{}
	}
	return nil
}

func GetRule(name string, typeKind reflect.Kind) (ValidationRule, error) {
	if typeKind != reflect.Int && typeKind != reflect.String {
		err := TypeError{fmt.Errorf("unsupported type %v", typeKind)}
		return nil, err
	}

	isAvailableRule := false
	for _, v := range avaliableRules[typeKind] {
		isAvailableRule = isAvailableRule || (v == name)
	}

	if !isAvailableRule {
		return nil, errGetRule
	}

	rule := getRuleByName(name)
	if rule != nil {
		return rule, nil
	}

	return nil, errGetRule
}

type MinRule struct{}

func (r MinRule) CheckValue(rule string, value interface{}) error {
	if value == nil {
		return errNilValue
	}
	if len(rule) == 0 {
		return fmt.Errorf(strParsingError, rule)
	}

	ruleInt, err := strconv.Atoi(rule)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(value)

	if val.Type().Kind() != reflect.Int {
		return fmt.Errorf(strTypeError, "int")
	}

	if iVal, _ := strconv.Atoi(getStrValue(val)); ruleInt <= iVal {
		return nil
	}
	return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
}

type MaxRule struct{}

func (r MaxRule) CheckValue(rule string, value interface{}) error {
	if value == nil {
		return errNilValue
	}
	if len(rule) == 0 {
		return fmt.Errorf(strParsingError, rule)
	}

	ruleInt, err := strconv.Atoi(rule)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(value)

	if val.Type().Kind() != reflect.Int {
		return fmt.Errorf(strTypeError, "int")
	}

	if iVal, _ := strconv.Atoi(getStrValue(val)); ruleInt >= iVal {
		return nil
	}

	return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
}

type LengthRule struct{}

func (r LengthRule) CheckValue(rule string, value interface{}) error {
	if value == nil {
		return errNilValue
	}
	if len(rule) == 0 {
		return fmt.Errorf(strParsingError, rule)
	}

	ruleInt, err := strconv.Atoi(rule)
	if err != nil {
		return err
	}
	val := reflect.ValueOf(value)

	if val.Type().Kind() != reflect.String {
		return fmt.Errorf(strTypeError, "string")
	}

	if ruleInt == len(getStrValue(val)) {
		return nil
	}

	return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
}

type RegexRule struct{}

func (r RegexRule) CheckValue(rule string, value interface{}) error {
	if value == nil {
		return errNilValue
	}
	if len(rule) == 0 {
		return fmt.Errorf(strParsingError, rule)
	}

	reg, err := regexp.Compile(rule)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(value)

	if val.Type().Kind() != reflect.String {
		return fmt.Errorf(strTypeError, "string")
	}

	if !reg.MatchString(getStrValue(val)) {
		return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
	}
	return nil
}

type InRule struct{}

func (r InRule) CheckValue(rule string, value interface{}) error {
	if value == nil {
		return errNilValue
	}
	if len(rule) == 0 {
		return fmt.Errorf(strParsingError, rule)
	}
	val := reflect.ValueOf(value)
	if val.Type().Kind() == reflect.Int {
		elements := strings.Split(rule, ",")
		if len(elements) == 0 {
			return fmt.Errorf(strParsingError, rule)
		}

		intValue, _ := strconv.Atoi(getStrValue(val))

		for _, e := range elements {
			ie, err := strconv.Atoi(e)
			if err != nil {
				return err
			}

			if ie == intValue {
				return nil
			}
		}
		return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
	}
	if val.Type().Kind() == reflect.String {
		elements := strings.Split(rule, ",")
		if len(elements) == 0 {
			return fmt.Errorf(strParsingError, rule)
		}
		strValue := getStrValue(val)
		for _, e := range elements {
			if e == strValue {
				return nil
			}
		}
		return ValidationError{Err: fmt.Errorf(strValidationError, rule, value)}
	}
	return fmt.Errorf(strTypeError, "int or string")
}
