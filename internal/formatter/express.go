package formatter

import (
	"fmt"
	"strings"

	"github.com/yusnelgg/scribe/internal/parser"
)

type ExpressFormatter struct{}

func NewExpressFormatter() *ExpressFormatter {
	return &ExpressFormatter{}
}

func (f *ExpressFormatter) FormatTest(route parser.Route) string {
	var buf strings.Builder

	buf.WriteString("const request = require('supertest');\n")
	buf.WriteString("const express = require('express');\n\n")

	buf.WriteString("describe('" + route.Method + " " + route.Path + "', () => {\n")
	buf.WriteString("\tlet app;\n\n")

	buf.WriteString("\tbeforeAll(() => {\n")
	buf.WriteString("\t\tapp = express();\n")
	buf.WriteString("\t\tapp.use(express.json());\n")

	if route.Handler != "anonymous" && route.Handler != "" {
		buf.WriteString(fmt.Sprintf("\t\tapp.%s('%s', %s);\n",
			strings.ToLower(route.Method), route.Path, route.Handler))
	} else {
		buf.WriteString(fmt.Sprintf("\t\tapp.%s('%s', (req, res) => res.status(200).json({}));\n",
			strings.ToLower(route.Method), route.Path))
	}

	buf.WriteString("\t});\n\n")

	buf.WriteString("\ttest('should return 200', async () => {\n")
	buf.WriteString(fmt.Sprintf("\t\tconst response = await request(app)\n"))

	method := strings.ToUpper(route.Method)
	switch method {
	case "GET":
		buf.WriteString("\t\t\t.get('" + route.Path + "')\n")
	case "POST":
		buf.WriteString("\t\t\t.post('" + route.Path + "')\n")
	case "PUT":
		buf.WriteString("\t\t\t.put('" + route.Path + "')\n")
	case "DELETE":
		buf.WriteString("\t\t\t.delete('" + route.Path + "')\n")
	case "PATCH":
		buf.WriteString("\t\t\t.patch('" + route.Path + "')\n")
	}

	buf.WriteString("\t\t\t.send({});\n\n")
	buf.WriteString("\t\texpect(response.status).toBe(200);\n")
	buf.WriteString("\t});\n\n")

	if method == "POST" || method == "PUT" || method == "PATCH" {
		buf.WriteString("\ttest('should return 400 for invalid data', async () => {\n")
		buf.WriteString(fmt.Sprintf("\t\tconst response = await request(app)\n"))

		switch method {
		case "POST":
			buf.WriteString("\t\t\t.post('" + route.Path + "')\n")
		case "PUT":
			buf.WriteString("\t\t\t.put('" + route.Path + "')\n")
		case "PATCH":
			buf.WriteString("\t\t\t.patch('" + route.Path + "')\n")
		}

		buf.WriteString("\t\t\t.send({ invalid: 'data' });\n\n")
		buf.WriteString("\t\texpect(response.status).not.toBe(500);\n")
		buf.WriteString("\t});\n")
	}

	buf.WriteString("});\n")

	return buf.String()
}
