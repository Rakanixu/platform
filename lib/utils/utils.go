package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/globals"
	micro_errors "github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"io"
)

// ParseUserIdFromContext returns user_id from context
func ParseUserIdFromContext(ctx context.Context) (string, error) {
	if ctx.Value(kazoup_context.UserIdCtxKey{}) == nil {
		return "", micro_errors.Unauthorized("ParseUserIdFromContext", "Unable to retrieve user from context")
	}

	id := string(ctx.Value(kazoup_context.UserIdCtxKey{}).(kazoup_context.UserIdCtxValue))
	if len(id) == 0 {
		return "", micro_errors.Unauthorized("ParseUserIdFromContext", "No user for given context")
	}

	return id, nil
}

func ParseRolesFromContext(ctx context.Context) ([]string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return []string{}, errors.New("Unable to retrieve metadata")
	}

	if len(md["Authorization"]) == 0 {
		return []string{}, errors.New("No Auth header")
	}

	// We will read claim to know if public user, or paying or whatever
	token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
		if err != nil {
			return nil, err
		}

		return decoded, nil
	})
	if err != nil {
		return []string{}, err
	}

	if token.Claims.(jwt.MapClaims)["roles"] == nil {
		return []string{}, errors.New("Roles not found.")
	}

	var roles []string
	for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
		roles = append(roles, v.(string))
	}

	return roles, nil
}

// ParseJWTToken validates JWT and returns user_id claim
func ParseJWTToken(str string) (string, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, micro_errors.InternalServerError("Unexpected signing method", token.Header["alg"].(string))
		}

		decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
		if err != nil {
			return nil, err
		}

		return decoded, nil
	})

	if err != nil {
		return "", micro_errors.Unauthorized("Token", err.Error())
	}

	if !token.Valid {
		return "", micro_errors.Unauthorized("", "Invalid token")
	}

	return token.Claims.(jwt.MapClaims)["sub"].(string), nil
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetMimeType(fileSystemType, fileType string) string {
	// Be sure to not panic if input not in map
	if globals.FileTypeDict.M[fileSystemType] != nil {
		if len(globals.FileTypeDict.M[fileSystemType][fileType]) > 0 {
			return globals.FileTypeDict.M[fileSystemType][fileType]
		}
	}

	return ""
}

func GoogleDriveExportAs(originalMimeType string) string {
	//https://developers.google.com/drive/v3/web/integrate-open#open_and_convert_google_docs_in_your_app
	switch originalMimeType {
	case globals.GOOGLE_DRIVE_DOCUMENT:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case globals.GOOGLE_DRIVE_PRESETATION:
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case globals.GOOGLE_DRIVE_SPREADSHEET:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	return ""
}

func GetDocumentTemplate(fileType string, fullName bool) string {
	var tmp string

	switch fileType {
	case globals.DOCUMENT:
		tmp = "docx"
	case globals.PRESENTATION:
		tmp = "pptx"
	case globals.SPREADSHEET:
		tmp = "xlsx"
	case globals.TEXT:
		tmp = "txt"
	}

	if fullName {
		tmp = fmt.Sprintf("%s.%s", tmp, tmp)
	}

	return tmp
}

// Encrypt slice of bytes
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt slice of bytes
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
