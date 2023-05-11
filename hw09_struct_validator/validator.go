package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/SergeyTyurin/OtusGolang/hw09_struct_validator/internal/rules"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	builder.WriteString("Validation Error:\n")
	for _, e := range v {
		builder.WriteString("\tField: " + e.Field + ": " + e.Err.Error() + "\n")
	}
	return builder.String()
}

func Validate(v interface{}) error {
	if v == nil {
		return fmt.Errorf("interface is nil")
	}

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("input isn not struct")
	}

	return validateStruct(v)
}

func validateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	validationErrors := make(ValidationErrors, 0)

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		var ruleValidationError rules.ValidationError
		err := validateField(t.Field(i), v.Field(i).Interface())
		if err != nil {
			if errors.As(err, &ruleValidationError) {
				vError := ValidationError{Field: t.Field(i).Name, Err: err}
				validationErrors = append(validationErrors, vError)
			} else {
				return err
			}
		}
	}
	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func validateField(field reflect.StructField, value interface{}) error {
	fieldType := field.Type.Kind()
	if fieldType == reflect.Slice {
		fieldType = field.Type.Elem().Kind()
	}
	if fieldType != reflect.Int && fieldType != reflect.String {
		return nil
	}

	strTag, ok := field.Tag.Lookup("validate")
	if !ok {
		return nil
	}

	if field.Type.Kind() == reflect.Slice {
		reflectedSlice := reflect.ValueOf(value)
		for i := 0; i < reflectedSlice.Len(); i++ {
			err := checkValue(strTag, reflectedSlice.Index(i).Interface())
			if err != nil {
				return err
			}
		}
		return nil
	}
	return checkValue(strTag, value)
}

func checkValue(tag string, value interface{}) error {
	fieldType := reflect.ValueOf(value).Type().Kind()
	rulesStr := strings.Split(tag, "|")
	for _, ruleStr := range rulesStr {
		pair := strings.Split(ruleStr, ":")
		if len(pair) != 2 {
			return fmt.Errorf("incorrect format for tag %s", tag)
		}
		rule, err := rules.GetRule(pair[0], fieldType)
		if err != nil {
			typeError := rules.TypeError{}
			if errors.As(err, &typeError) {
				continue
			} else {
				return err
			}
		}

		err = rule.CheckValue(pair[1], value)
		if err != nil {
			return err
		}
	}
	return nil
}
