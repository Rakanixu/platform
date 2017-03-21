package text

import (
	"bytes"
	"golang.org/x/text/transform"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

func NormalizeBytes(b []byte) (string, error) {
	skipped := transform.NewReader(bytes.NewReader(b), newSkipper(5))

	var buf bytes.Buffer
	_, err := io.Copy(&buf, skipped)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Normalize(s string) (string, error) {
	skipped := transform.NewReader(strings.NewReader(s), newSkipper(5))

	var buf bytes.Buffer
	_, err := io.Copy(&buf, skipped)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ReplaceTabs(s string) (string, error) {
	t, err := regexp.Compile("\t")
	if err != nil {
		return "", err
	}

	return t.ReplaceAllString(s, " "), nil
}

func ReplaceNewLines(s string) (string, error) {
	nl, err := regexp.Compile("\n")
	if err != nil {
		return "", err
	}

	return nl.ReplaceAllString(s, " "), nil
}

func ReplaceDoubleQuotes(s string) (string, error) {
	q, err := regexp.Compile("\"")
	if err != nil {
		return "", err
	}

	return q.ReplaceAllString(s, "'"), nil
}

type skipper struct {
	pos int
	cnt int
}

// NewSkipper creates a text transformer which will remove the rune at pos
func newSkipper(pos int) transform.Transformer {
	return &skipper{pos: pos}
}

func (s *skipper) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	for utf8.FullRune(src) {
		_, sz := utf8.DecodeRune(src)
		// not enough space in the dst
		if len(dst) < sz {
			return nDst, nSrc, transform.ErrShortDst
		}
		if s.pos != s.cnt {
			copy(dst[:sz], src[:sz])
			// track that we stored in dst
			dst = dst[sz:]
			nDst += sz
		}
		// track that we read from src
		src = src[sz:]
		nSrc += sz
		// on to the next rune
		s.cnt++
	}
	if len(src) > 0 && !atEOF {
		return nDst, nSrc, transform.ErrShortSrc
	}
	return nDst, nSrc, nil
}

func (s *skipper) Reset() {
	s.cnt = 0
}
