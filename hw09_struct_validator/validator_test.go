package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func jsonBlob(filename string) []byte {
	blob, _ := os.ReadFile(filename)
	return blob
}

func TestValidate(t *testing.T) {
	var correctUser User
	var correctApp App
	var correctResponse Response
	token := Token{}
	var incorrectUser User
	var incorrectApp App
	var incorrectResponse Response

	json.Unmarshal(jsonBlob("testdata/correct/user.json"), &correctUser)
	json.Unmarshal(jsonBlob("testdata/correct/app.json"), &correctApp)
	json.Unmarshal(jsonBlob("testdata/correct/response.json"), &correctResponse)
	json.Unmarshal(jsonBlob("testdata/wrong/user.json"), &incorrectUser)
	json.Unmarshal(jsonBlob("testdata/wrong/app.json"), &incorrectApp)
	json.Unmarshal(jsonBlob("testdata/wrong/response.json"), &incorrectResponse)

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: correctUser, expectedErr: nil},
		{in: correctApp, expectedErr: nil},
		{in: correctResponse, expectedErr: nil},
		{in: token, expectedErr: nil},
		{in: incorrectUser, expectedErr: ValidationErrors{}},
		{in: incorrectApp, expectedErr: ValidationErrors{}},
		{in: incorrectResponse, expectedErr: ValidationErrors{}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			if tt.expectedErr != nil {
				require.ErrorAs(t, err, &tt.expectedErr)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestProgramErrors(t *testing.T) {
	t.Run("input case", func(t *testing.T) {
		require.NotNil(t, Validate(nil))
		require.NotNil(t, Validate(""))
	})

	tests := []struct {
		in interface{}
	}{
		{in: struct {
			Field int `validate:"in200,404,500"`
		}{Field: 1}},
		{in: struct {
			Field int `validate:"min:a"`
		}{Field: 1}},
		{in: struct {
			Field int `validate:"len:5"`
		}{Field: 1}},
		{in: struct {
			Field string `validate:"len:"`
		}{Field: "1"}},
		{in: struct {
			Field string `validate:"length:1"`
		}{Field: "1"}},
	}

	t.Run("parsing case", func(t *testing.T) {
		var validationErrors ValidationErrors
		for _, test := range tests {
			err := Validate(test.in)
			require.NotNil(t, err)
			require.False(t, errors.As(err, &validationErrors))
		}
	})
}
