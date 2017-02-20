package googlecloudstorage

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

/*
FIXME:
Following code not in generated library
This is a copy of https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/storage/storage.go
This may be broken when library is updated by google
*/

// SignedURLOptions allows you to restrict the access to the signed URL.
type SignedURLOptions struct {
	// GoogleAccessID represents the authorizer of the signed URL generation.
	// It is typically the Google service account client email address from
	// the Google Developers Console in the form of "xxx@developer.gserviceaccount.com".
	// Required.
	GoogleAccessID string

	// PrivateKey is the Google service account private key. It is obtainable
	// from the Google Developers Console.
	// At https://console.developers.google.com/project/<your-project-id>/apiui/credential,
	// create a service account client ID or reuse one of your existing service account
	// credentials. Click on the "Generate new P12 key" to generate and download
	// a new private key. Once you download the P12 file, use the following command
	// to convert it into a PEM file.
	//
	//    $ openssl pkcs12 -in key.p12 -passin pass:notasecret -out key.pem -nodes
	//
	// Provide the contents of the PEM file as a byte slice.
	// Exactly one of PrivateKey or SignBytes must be non-nil.
	PrivateKey []byte

	// SignBytes is a function for implementing custom signing.
	// If your application is running on Google App Engine, you can use appengine's internal signing function:
	//     ctx := appengine.NewContext(request)
	//     acc, _ := appengine.ServiceAccount(ctx)
	//     url, err := SignedURL("bucket", "object", &SignedURLOptions{
	//     	GoogleAccessID: acc,
	//     	SignBytes: func(b []byte) ([]byte, error) {
	//     		_, signedBytes, err := appengine.SignBytes(ctx, b)
	//     		return signedBytes, err
	//     	},
	//     	// etc.
	//     })
	//
	// Exactly one of PrivateKey or SignBytes must be non-nil.
	SignBytes func([]byte) ([]byte, error)

	// Method is the HTTP method to be used with the signed URL.
	// Signed URLs can be used with GET, HEAD, PUT, and DELETE requests.
	// Required.
	Method string

	// Expires is the expiration time on the signed URL. It must be
	// a datetime in the future.
	// Required.
	Expires time.Time

	// ContentType is the content type header the client must provide
	// to use the generated signed URL.
	// Optional.
	ContentType string

	// Headers is a list of extention headers the client must provide
	// in order to use the generated signed URL.
	// Optional.
	Headers []string

	// MD5 is the base64 encoded MD5 checksum of the file.
	// If provided, the client should provide the exact value on the request
	// header in order to use the signed URL.
	// Optional.
	MD5 []byte
}

// SignedURL returns a URL for the specified object. Signed URLs allow
// the users access to a restricted resource for a limited time without having a
// Google account or signing in. For more information about the signed
// URLs, see https://cloud.google.com/storage/docs/accesscontrol#Signed-URLs.
func SignedURL(bucket, name string, opts *SignedURLOptions) (string, error) {
	if opts == nil {
		return "", errors.New("storage: missing required SignedURLOptions")
	}
	if opts.GoogleAccessID == "" {
		return "", errors.New("storage: missing required GoogleAccessID")
	}
	if (opts.PrivateKey == nil) == (opts.SignBytes == nil) {
		return "", errors.New("storage: exactly one of PrivateKey or SignedBytes must be set")
	}
	if opts.Method == "" {
		return "", errors.New("storage: missing required method option")
	}
	if opts.Expires.IsZero() {
		return "", errors.New("storage: missing required expires option")
	}

	signBytes := opts.SignBytes
	if opts.PrivateKey != nil {
		key, err := parseKey(opts.PrivateKey)
		if err != nil {
			return "", err
		}
		signBytes = func(b []byte) ([]byte, error) {
			sum := sha256.Sum256(b)
			return rsa.SignPKCS1v15(
				rand.Reader,
				key,
				crypto.SHA256,
				sum[:],
			)
		}
	} else {
		signBytes = opts.SignBytes
	}

	u := &url.URL{
		Path: fmt.Sprintf("/%s/%s", bucket, name),
	}

	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s\n", opts.Method)
	fmt.Fprintf(buf, "%s\n", opts.MD5)
	fmt.Fprintf(buf, "%s\n", opts.ContentType)
	fmt.Fprintf(buf, "%d\n", opts.Expires.Unix())
	fmt.Fprintf(buf, "%s", strings.Join(opts.Headers, "\n"))
	fmt.Fprintf(buf, "%s", u.String())

	b, err := signBytes(buf.Bytes())
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	u.Scheme = "https"
	u.Host = "storage.googleapis.com"
	q := u.Query()
	q.Set("GoogleAccessId", opts.GoogleAccessID)
	q.Set("Expires", fmt.Sprintf("%d", opts.Expires.Unix()))
	q.Set("Signature", string(encoded))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// parseKey converts the binary contents of a private key file
// to an *rsa.PrivateKey. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func parseKey(key []byte) (*rsa.PrivateKey, error) {
	if block, _ := pem.Decode(key); block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, err
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("oauth2: private key is invalid")
	}
	return parsed, nil
}
