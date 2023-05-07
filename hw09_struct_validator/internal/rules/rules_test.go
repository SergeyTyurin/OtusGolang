package rules

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckValue(t *testing.T) {
	minRule := MinRule{}
	maxRule := MaxRule{}
	lenRule := LengthRule{}
	regexRule := RegexRule{}
	inRule := InRule{}

	t.Run("correct data case", func(t *testing.T) {
		require.Nil(t, minRule.CheckValue("4", 5))
		require.Nil(t, maxRule.CheckValue("4", 4))
		require.Nil(t, lenRule.CheckValue("4", "abcd"))
		require.Nil(t, regexRule.CheckValue("\\d+", "1234"))
		require.Nil(t, inRule.CheckValue("1,2,3", 1))
	})

	t.Run("incorrect data case", func(t *testing.T) {
		var validationError ValidationError
		require.NotNil(t, minRule.CheckValue("4", "4"))
		require.NotNil(t, minRule.CheckValue("a", 5))
		require.ErrorAs(t, minRule.CheckValue("4", 2), &validationError)

		require.NotNil(t, maxRule.CheckValue("4", "4"))
		require.NotNil(t, maxRule.CheckValue("a", 5))
		require.ErrorAs(t, maxRule.CheckValue("4", 5), &validationError)

		require.NotNil(t, lenRule.CheckValue("_", "4"))
		require.NotNil(t, lenRule.CheckValue("3", 5))
		require.ErrorAs(t, lenRule.CheckValue("3", "abcd"), &validationError)

		require.NotNil(t, regexRule.CheckValue("\\d+", 4))
		require.NotNil(t, regexRule.CheckValue("\\", "4"))
		require.ErrorAs(t, regexRule.CheckValue("\\d+", "abc"), &validationError)

		require.ErrorAs(t, inRule.CheckValue("1,2,3", "4"), &validationError)
		require.ErrorAs(t, inRule.CheckValue("1,2,3", 4), &validationError)
	})

	t.Run("empty data case", func(t *testing.T) {

		require.NotNil(t, minRule.CheckValue("", "4"))
		require.NotNil(t, minRule.CheckValue("a", nil))

		require.NotNil(t, maxRule.CheckValue("", "4"))
		require.NotNil(t, maxRule.CheckValue("a", nil))

		require.NotNil(t, lenRule.CheckValue("", "4"))
		require.NotNil(t, lenRule.CheckValue("3", nil))

		require.NotNil(t, regexRule.CheckValue("\\d+", nil))
		require.NotNil(t, regexRule.CheckValue("", "4"))

		require.NotNil(t, inRule.CheckValue("", "4"))
		require.NotNil(t, inRule.CheckValue("1,2,3", nil))
	})
}

func TestGetRule(t *testing.T) {
	t.Run("correct case", func(t *testing.T) {
		tests := []struct {
			name string
			kind reflect.Kind
		}{
			{"min", reflect.Int},
			{"max", reflect.Int},
			{"len", reflect.String},
			{"regexp", reflect.String},
			{"in", reflect.String},
			{"in", reflect.Int},
		}

		for _, test := range tests {
			_, err := GetRule(test.name, test.kind)
			require.Nil(t, err)
		}
	})

	t.Run("incorrect case", func(t *testing.T) {
		tests := []struct {
			name string
			kind reflect.Kind
		}{
			{"min", reflect.String},
			{"max", reflect.String},
			{"len", reflect.Int},
			{"regexp", reflect.Float32},
			{"in", reflect.Float32},
			{"in", reflect.Array},
			{"incorrect", reflect.String},
		}

		for _, test := range tests {
			_, err := GetRule(test.name, test.kind)
			require.NotNil(t, err)
		}
	})
}
