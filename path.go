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
func newPath(path string, optionalSlash bool) (*Path, error) {
	params := []string{}

	var regexPath string
	if strings.HasSuffix(path, "/*") {
		regexPath = "^" + strings.ReplaceAll(path[:len(path)-2], "/", "\\/") + "(?:\\/(?P<path>.*))?$"

		params = append(params, "path")
	} else {
		// Replace OpenAPI wildcards with Go RegExp named wildcards
		regexPath = regOpenAPIPath.ReplaceAllStringFunc(path, func(m string) string {
			param := m[1 : len(m)-1]          // Strip outer braces from parameter
			params = append(params, param)    // Add param to list
			return "(?P<" + param + ">[^/]+)" // Replace with Go Regexp named wildcard
		})

		// Make sure the path starts with a "^", ends with a "$", and escape slashes
		regexPath = "^" + strings.ReplaceAll(regexPath, "/", "\\/") + "$"

		if optionalSlash && strings.HasSuffix(regexPath, "\\/$") {
			regexPath = regexPath[:len(regexPath)-3] + "\\/?$"
		}
	}

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

	matches := map[string]string{}

	for i, name := range p.reg.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}

		if name == "path" {
			matches["*"] = ms[i]
		} else {
			matches[name] = ms[i]
		}
	}

	return matches
}

func (p *Path) contains(v string) bool {
	for _, i := range p.names {
		if i == v {
			return true
		}
	}

	return false
}
