package goapi

import (
	"regexp"
	"strings"
)

// Path helps to handle openapi path pattern.
type Path struct {
	path  string
	reg   *regexp.Regexp
	names []string
}

var regOpenAPIPath = regexp.MustCompile(`\{([^}]+)\}`)

// Converts OpenAPI style path to Go Regexp and returns path parameters.
func newPath(path string) (*Path, error) {
	params := []string{}

	// Replace OpenAPI wildcards with Go RegExp named wildcards
	regexPath := regOpenAPIPath.ReplaceAllStringFunc(path, func(m string) string {
		param := m[1 : len(m)-1]          // Strip outer braces from parameter
		params = append(params, param)    // Add param to list
		return "(?P<" + param + ">[^/]+)" // Replace with Go Regexp named wildcard
	})

	// Make sure the path starts with a "^", ends with a "$", and escape slashes
	regexPath = "^" + strings.ReplaceAll(regexPath, "/", "\\/") + "$"

	// Compile the regular expression
	r, err := regexp.Compile(regexPath)
	if err != nil {
		return nil, err
	}

	return &Path{path, r, params}, nil
}

func (p *Path) match(path string) map[string]string {
	ms := p.reg.FindStringSubmatch(path)

	if ms == nil {
		return nil
	}

	params := map[string]string{}

	for i, m := range ms[1:] {
		params[p.names[i]] = m
	}

	return params
}

func (p *Path) contains(v string) bool {
	for _, i := range p.names {
		if i == v {
			return true
		}
	}

	return false
}
