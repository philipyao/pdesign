package pconfclient

import (
	"strings"
)

// GonicMapper will consider initialisms when mapping names.
// E.g. id -> ID, user_name -> UserName
type gonicMapper map[string]bool

func (mapper gonicMapper) Tag2GoName(tag string) string {
	if tag == "" {
		return ""
	}
	newstr := make([]rune, 0)

	tag = strings.ToLower(tag)
	parts := strings.Split(tag, "_")

	for _, p := range parts {
		_, isInitialism := mapper[strings.ToUpper(p)]
		for i, r := range p {
			if i == 0 || isInitialism {
				r = toASCIIUpper(r)
			}
			newstr = append(newstr, r)
		}
	}

	return string(newstr)
}

func toASCIIUpper(r rune) rune {
	if 'a' <= r && r <= 'z' {
		r -= ('a' - 'A')
	}
	return r
}

// A GonicMapper that contains a list of common initialisms taken from golang/lint
var lintGonicMapper = gonicMapper{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

func tag2GoName(tag string) string {
	return lintGonicMapper.Tag2GoName(tag)
}
