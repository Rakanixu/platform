package utils

import (
	"crypto/aes"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
	JWT_TOKEN    = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJyb2xlcyI6WyJwZXJzb25hbCJdLCJpc3MiOiJodHRwczovL2them91cC5ldS5hdXRoMC5jb20vIiwic3ViIjoiZ29vZ2xlLWFwcHN8cGFibG8uYWd1aXJyZUBrYXpvdXAuY29tIiwiYXVkIjoiNnpJRG04SW5oYlRScDFiTDJDNG0xVEs0TGxyNGFyVHkiLCJleHAiOjE0OTQ1MjY1OTMsImlhdCI6MTQ5NDQ5MDU5MywiYXpwIjoiNU9DSll1VHE1RG9nOTYwYzNsZlZFc0JscXVEWDlLYTIifQ.H_pgMvVg6u4Jern6fetiZfsBAczZ9iwpKZac0FDBuUY"
)

func TestParseUserIdFromContext(t *testing.T) {
	type out struct {
		str string
		err error
	}

	testData := []struct {
		in  context.Context
		out out
	}{
		{
			in: context.WithValue(
				context.TODO(),
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue(TEST_USER_ID),
			),
			out: out{
				str: TEST_USER_ID,
				err: nil,
			},
		},
		{
			in: context.WithValue(
				context.TODO(),
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue(""),
			),
			out: out{
				str: "",
				err: errors.ErrNoUserInCtx,
			},
		},
		{
			in: context.TODO(),
			out: out{
				str: "",
				err: errors.ErrInvalidUserInCtx,
			},
		},
	}

	for _, tt := range testData {
		result, err := ParseUserIdFromContext(tt.in)

		if tt.out.str != result {
			t.Errorf("Expected %v, got %v", tt.out.str, result)
		}

		if tt.out.err != err {
			t.Errorf("Expected %v, got %v", tt.out.err, err)
		}
	}
}

/* Token will expired and test will fail.*/
/*
func TestParseRolesFromContext(t *testing.T) {
	type out struct {
		roles []string
		err   error
	}

	testData := []struct {
		in  context.Context
		out out
	}{
		{
			in: metadata.NewContext(
				context.TODO(),
				metadata.Metadata{"Authorization": JWT_TOKEN},
			),
			out: out{
				roles: []string{"personal"},
				err:   nil,
			},
		},
		{
			in: context.TODO(),
			out: out{
				roles: []string{},
				err:   errors.ErrInvalidMetadata,
			},
		},
		{
			in: metadata.NewContext(
				context.TODO(),
				metadata.Metadata{"Authorization": ""},
			),
			out: out{
				roles: []string{},
				err:   errors.ErrNoAuthHeader,
			},
		},
	}

	for _, tt := range testData {
		result, err := ParseRolesFromContext(tt.in)

		if !reflect.DeepEqual(tt.out.roles, result) {
			t.Errorf("Expected %v, got %v", tt.out.roles, result)
		}

		if tt.out.err != err {
			t.Errorf("Expected %v, got %v", tt.out.err, err)
		}
	}
}

func TestParseJWTToken(t *testing.T) {
	type out struct {
		userId string
		err    error
	}

	testData := []struct {
		in  string
		out out
	}{
		{
			in: JWT_TOKEN,
			out: out{
				userId: "",
				err:    nil,
			},
		},
	}

	for _, tt := range testData {
		result, err := ParseJWTToken(tt.in)

		if tt.out.err != err {
			t.Errorf("Expected %v, got %v", tt.out.err, err)
		}
	}
}
*/

func TestNewUUID(t *testing.T) {
	str, err := NewUUID()

	if err != nil {
		t.Error(err)
	}

	if len(str) <= 0 {
		t.Errorf("NewUUID empty")
	}
}

func TestGetMD5Hash(t *testing.T) {
	testData := []struct {
		in string
	}{
		{
			in: "abcde",
		},
		{
			in: "a213123123123bcd1e",
		},
	}

	for _, tt := range testData {
		if GetMD5Hash(tt.in) != GetMD5Hash(tt.in) {
			t.Error("Different hash for same string")
		}
	}
}

func TestGetMimeType(t *testing.T) {
	if len(GetMimeType("1", "2")) > 0 {
		t.Error("Incorrect MimeType")
	}
}

func TestGoogleDriveExportAs(t *testing.T) {
	if len(GoogleDriveExportAs("1")) > 0 {
		t.Error("Incorrect original MimeType")
	}
}

func TestGetDocumentTemplate(t *testing.T) {
	type in struct {
		fileType string
		fullName bool
	}

	testData := []struct {
		in  in
		out string
	}{
		{
			in: in{
				fileType: globals.DOCUMENT,
				fullName: true,
			},
			out: "docx.docx",
		},
		{
			in: in{
				fileType: globals.PRESENTATION,
				fullName: true,
			},
			out: "pptx.pptx",
		},
		{
			in: in{
				fileType: globals.SPREADSHEET,
				fullName: true,
			},
			out: "xlsx.xlsx",
		},
		{
			in: in{
				fileType: globals.TEXT,
				fullName: true,
			},
			out: "txt.txt",
		},
		{
			in: in{
				fileType: globals.DOCUMENT,
			},
			out: "docx",
		},
		{
			in: in{
				fileType: globals.PRESENTATION,
			},
			out: "pptx",
		},
		{
			in: in{
				fileType: globals.SPREADSHEET,
			},
			out: "xlsx",
		},
		{
			in: in{
				fileType: globals.TEXT,
			},
			out: "txt",
		},
	}

	for _, tt := range testData {
		result := GetDocumentTemplate(tt.in.fileType, tt.in.fullName)
		if tt.out != result {
			t.Errorf("Expected %v, got %v", tt.out, result)
		}
	}
}

func TestEncrypt(t *testing.T) {
	type in struct {
		key  []byte
		text []byte
	}

	testData := []struct {
		in  in
		err error
	}{
		{
			in: in{
				key:  []byte(globals.ENCRYTION_KEY_32),
				text: []byte("info"),
			},
			err: nil,
		},
		{
			in: in{
				key:  []byte(""),
				text: []byte("info"),
			},
			err: aes.KeySizeError(0),
		},
	}

	for _, tt := range testData {
		result, err := Encrypt(tt.in.key, tt.in.text)

		if reflect.DeepEqual(tt.in.text, result) {
			t.Errorf("Expected %v, got %v", tt.in.text, result)
		}

		if tt.err != err {
			t.Errorf("Expected %v, got %v", tt.err, err)
		}
	}
}

func TestDecrypt(t *testing.T) {
	testData := []struct {
		key  []byte
		text []byte
		err  error
	}{
		{
			key:  []byte(globals.ENCRYTION_KEY_32),
			text: []byte("info"),
			err:  nil,
		},
		{
			key: []byte(globals.ENCRYTION_KEY_32),
			text: []byte(`3142asd asdf
			asdf	dsfgsdf\n asfasdf`),
			err: nil,
		},
	}

	for _, tt := range testData {
		result, err := Encrypt(tt.key, tt.text)
		if err != nil {
			t.Error(err)
		}

		decrypted, err := Decrypt(tt.key, result)
		if err != nil {
			t.Error(err)
		}

		if string(tt.text) != string(decrypted) {
			t.Errorf("Expected %v, got %v", string(tt.text), decrypted)
		}
	}
}
