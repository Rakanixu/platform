package handler

import (
	proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"strings"
)

type sortByRate []*proto.Quota

func (s sortByRate) Len() int {
	return len(s)
}

func (s sortByRate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortByRate) Less(i, j int) bool {
	return s[i].Rate > s[j].Rate
}

type sortAlphabetically []*proto.Quota

func (s sortAlphabetically) Len() int {
	return len(s)
}

func (s sortAlphabetically) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortAlphabetically) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) <= strings.ToLower(s[j].Name)
}
