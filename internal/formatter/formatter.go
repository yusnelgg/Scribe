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
	buf.WriteString("\t\"context\"\n")
	buf.WriteString("\t\"encoding/json\"\n")
	buf.WriteString("\t\"net/http\"\n")
	buf.WriteString("\t\"net/http/httptest\"\n")
	buf.WriteString("\t\"testing\"\n")
	buf.WriteString("\n")
	buf.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	buf.WriteString("\t\"github.com/stretchr/testify/assert\"\n")
	buf.WriteString("\t\"github.com/stretchr/testify/require\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString(fmt.Sprintf("func Test%s_%s(t *testing.T) {\n",
		strings.ToUpper(route.Method),
		formatTestName(route.Path)))

	buf.WriteString("\tt.Run(\"success case\", func(t *testing.T) {\n")
	buf.WriteString(fmt.Sprintf("\t\treq := httptest.NewRequest(%q, %q, nil)\n",
		strings.ToUpper(route.Method), route.Path))
	buf.WriteString("\t\treq.Header.Set(\"Content-Type\", \"application/json\")\n")
	buf.WriteString("\t\tw := httptest.NewRecorder()\n\n")
	buf.WriteString("\t\t// Setup router\n")
	buf.WriteString("\t\tgin.SetMode(gin.TestMode)\n")
	buf.WriteString("\t\trouter := gin.New()\n")
	buf.WriteString(fmt.Sprintf("\t\trouter.%s(%q, %s)\n",
		route.Method, route.Path, route.Handler))
	buf.WriteString("\n")
	buf.WriteString("\t\trouter.ServeHTTP(w, req)\n\n")
	buf.WriteString("\t\trequire.Equal(t, http.StatusOK, w.Code, \"expected status 200\")\n")
	buf.WriteString("\t})\n\n")

	buf.WriteString("\tt.Run(\"invalid request\", func(t *testing.T) {\n")
	buf.WriteString(fmt.Sprintf("\t\treq := httptest.NewRequest(%q, %q, nil)\n",
		strings.ToUpper(route.Method), route.Path))
	buf.WriteString("\t\treq.Header.Set(\"Content-Type\", \"application/json\")\n")
	buf.WriteString("\t\tw := httptest.NewRecorder()\n\n")
	buf.WriteString("\t\tgin.SetMode(gin.TestMode)\n")
	buf.WriteString("\t\trouter := gin.New()\n")
	buf.WriteString(fmt.Sprintf("\t\trouter.%s(%q, %s)\n",
		route.Method, route.Path, route.Handler))
	buf.WriteString("\n")
	buf.WriteString("\t\trouter.ServeHTTP(w, req)\n\n")
	buf.WriteString("\t\trequire.NotEqual(t, http.StatusInternalServerError, w.Code, \"should not return 500\")\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n\n")

	f.generateFixtures(&buf, route)

	return buf.String()
}

func (f *Formatter) generateFixtures(buf *strings.Builder, route parser.Route) {
	buf.WriteString("// Helper functions\n\n")
	buf.WriteString("func setupRouter() *gin.Engine {\n")
	buf.WriteString("\tgin.SetMode(gin.TestMode)\n")
	buf.WriteString("\treturn gin.New()\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func invalidRequest() *strings.Reader {\n")
	buf.WriteString("\treturn strings.NewReader(\"{invalid json}\")\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func assertResponseBody(t *testing.T, body string) {\n")
	buf.WriteString("\tassert.NotEmpty(t, body, \"response body should not be empty\")\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func assertJSONContentType(t *testing.T, w *httptest.ResponseRecorder) {\n")
	buf.WriteString("\tassert.Equal(t, \"application/json\", w.Header().Get(\"Content-Type\"))\n")
	buf.WriteString("}\n")
}

func formatTestName(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "-", "_")
	path = strings.ReplaceAll(path, ":", "_")
	return strings.ToUpper(path)
}
