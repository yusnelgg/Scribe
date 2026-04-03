package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ExpressParser struct{}

func NewExpressParser() *ExpressParser {
	return &ExpressParser{}
}

func (p *ExpressParser) Parse(file string) ([]Route, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var routes []Route
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	appVar := detectAppVariable(contentStr)

	patterns := []struct {
		regex  string
		method string
	}{
		{appVar + `.get\s*\(\s*"([^"]+)"`, "GET"},
		{appVar + `.get\s*\(\s*'([^']+)'`, "GET"},
		{appVar + `.post\s*\(\s*"([^"]+)"`, "POST"},
		{appVar + `.post\s*\(\s*'([^']+)'`, "POST"},
		{appVar + `.put\s*\(\s*"([^"]+)"`, "PUT"},
		{appVar + `.put\s*\(\s*'([^']+)'`, "PUT"},
		{appVar + `.delete\s*\(\s*"([^"]+)"`, "DELETE"},
		{appVar + `.delete\s*\(\s*'([^']+)'`, "DELETE"},
		{appVar + `.patch\s*\(\s*"([^"]+)"`, "PATCH"},
		{appVar + `.patch\s*\(\s*'([^']+)'`, "PATCH"},
	}

	for _, rp := range patterns {
		re := regexp.MustCompile(rp.regex)
		matches := re.FindAllStringSubmatch(contentStr, -1)

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			path := match[1]
			handler := "anonymous"

			if path == "" || path == "*" {
				continue
			}

			for i, line := range lines {
				if strings.Contains(line, path) {
					routes = append(routes, Route{
						Method:   rp.method,
						Path:     path,
						Handler:  handler,
						Package:  filepath.Base(file),
						Line:     i + 1,
						FilePath: file,
					})
					break
				}
			}
		}
	}

	routerPatterns := []struct {
		regex  string
		method string
	}{
		{`router\.get\s*\(\s*"([^"]+)"`, "GET"},
		{`router\.get\s*\(\s*'([^']+)'`, "GET"},
		{`router\.post\s*\(\s*"([^"]+)"`, "POST"},
		{`router\.post\s*\(\s*'([^']+)'`, "POST"},
		{`router\.put\s*\(\s*"([^"]+)"`, "PUT"},
		{`router\.put\s*\(\s*'([^']+)'`, "PUT"},
		{`router\.delete\s*\(\s*"([^"]+)"`, "DELETE"},
		{`router\.delete\s*\(\s*'([^']+)'`, "DELETE"},
		{`router\.patch\s*\(\s*"([^"]+)"`, "PATCH"},
		{`router\.patch\s*\(\s*'([^']+)'`, "PATCH"},
	}

	for _, rp := range routerPatterns {
		re := regexp.MustCompile(rp.regex)
		matches := re.FindAllStringSubmatch(contentStr, -1)

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			path := match[1]
			handler := "anonymous"

			if path == "" {
				continue
			}

			for i, line := range lines {
				if strings.Contains(line, path) {
					routes = append(routes, Route{
						Method:   rp.method,
						Path:     path,
						Handler:  handler,
						Package:  filepath.Base(file),
						Line:     i + 1,
						FilePath: file,
					})
					break
				}
			}
		}
	}

	return routes, nil
}

func detectAppVariable(content string) string {
	vars := []string{"app", "router", "express"}

	for _, v := range vars {
		if regexp.MustCompile(fmt.Sprintf(`const\s+%s\s*=\s*(?:express|Router)`, v)).MatchString(content) {
			return v
		}
	}

	for _, v := range vars {
		if strings.Contains(content, v+".get") || strings.Contains(content, v+".post") {
			return v
		}
	}

	return "app"
}
