package formatter

import (
	"fmt"
	"strings"

	"github.com/yusnelgg/scribe/internal/parser"
)

type Formatter struct{}

func New() *Formatter {
	return &Formatter{}
}

func (f *Formatter) FormatTest(route parser.Route) string {
	var buf strings.Builder

	buf.WriteString("package tests\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"net/http\"\n")
	buf.WriteString("\t\"net/http/httptest\"\n")
	buf.WriteString("\t\"testing\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString(fmt.Sprintf("func Test%s_%s(t *testing.T) {\n",
		strings.ToUpper(route.Method),
		formatTestName(route.Path)))

	buf.WriteString(fmt.Sprintf("\treq := httptest.NewRequest(%q, %q, nil)\n",
		strings.ToUpper(route.Method), route.Path))
	buf.WriteString("\tw := httptest.NewRecorder()\n\n")
	buf.WriteString("\t// TODO: Set up router and handler\n")
	buf.WriteString("\t// router := gin.New()\n")
	buf.WriteString(fmt.Sprintf("\t// router.%s(%q, %s)\n",
		route.Method, route.Path, route.Handler))
	buf.WriteString("\t// router.ServeHTTP(w, req)\n\n")

	buf.WriteString(fmt.Sprintf("\tif w.Code != http.StatusOK {\n"))
	buf.WriteString("\t\tt.Errorf(\"Expected status %d, got %d\", http.StatusOK, w.Code)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	return buf.String()
}

func formatTestName(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.ReplaceAll(path, ":", "_")
	return strings.ToUpper(path)
}
