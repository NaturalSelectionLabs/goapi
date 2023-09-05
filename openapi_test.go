package goapi_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/NaturalSelectionLabs/goapi"
	"github.com/ysmood/got"
	"github.com/ysmood/vary"
)

type Res interface {
	goapi.Response
}

var iRes = vary.New(new(Res))

type Res01 struct {
	goapi.StatusOK
	ID string `description:"response id"`
}

var _ = iRes.Add(Res01{})

type Res02 struct {
	goapi.StatusForbidden
	Error goapi.Error
}

func (Res02) Description() string {
	return "returns 403"
}

var _ = iRes.Add(Res02{})

func TestOpenAPI(t *testing.T) {
	g := got.T(t)

	r := goapi.New()
	tr := g.Serve()
	tr.Mux.Handle("/", r.Server())

	r.GET("/test", func(p struct {
		goapi.InURL
		ID string
	}, h struct {
		goapi.InHeader
		UA string
	}) Res {
		return Res01{}
	})

	doc := r.OpenAPI(nil).JSON()

	// Ensure you have nodejs installed
	{
		g.E(os.WriteFile("tmp/openapi.json", []byte(doc), 0666))
		out, err := exec.Command("npx", strings.Split("rdme openapi:validate tmp/openapi.json", " ")...).CombinedOutput()
		g.Desc("%s", out).E(err)
	}

	g.Eq(g.JSON(doc), map[string]interface{} /* len=4 */ {
		"components": map[string]interface{}{
			"schemas": map[string]interface{} /* len=5 */ {
				"Error": map[string]interface{} /* len=5 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.Error`, /* len=43 */
					"properties": map[string]interface{} /* len=5 */ {
						"code": map[string]interface{}{
							"type": "string",
						},
						"details": map[string]interface{} /* len=2 */ {
							"items": map[string]interface{}{
								"$ref": `#/components/schemas/Error`, /* len=26 */
							},
							"type": "array",
						},
						"innererror": map[string]interface{}{
							"type": "object",
						},
						"message": map[string]interface{}{
							"type": "string",
						},
						"target": map[string]interface{}{
							"type": "string",
						},
					},
					"title": "Error",
					"type":  "object",
				},
				"Res01": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi_test.Res01`, /* len=48 */
					"properties": map[string]interface{}{
						"ID": map[string]interface{} /* len=2 */ {
							"description": "response id",
							"type":        "string",
						},
					},
					"required": []interface{} /* len=1 cap=1 */ {
						"ID",
					},
					"title": "Res01",
					"type":  "object",
				},
				"Res02": map[string]interface{} /* len=6 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi_test.Res02`, /* len=48 */
					"properties": map[string]interface{}{
						"Error": map[string]interface{}{
							"$ref": `#/components/schemas/Error`, /* len=26 */
						},
					},
					"required": []interface{} /* len=1 cap=1 */ {
						"Error",
					},
					"title": "Res02",
					"type":  "object",
				},
				"StatusForbidden": map[string]interface{} /* len=4 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.StatusForbidden`, /* len=53 */
					"title":                              "StatusForbidden",
					"type":                               "object",
				},
				"StatusOK": map[string]interface{} /* len=4 */ {
					`additionalProperties` /* len=20 */ : false,
					"description":                        `github.com/NaturalSelectionLabs/goapi.StatusOK`, /* len=46 */
					"title":                              "StatusOK",
					"type":                               "object",
				},
			},
		},
		"info": map[string]interface{} /* len=2 */ {
			"title":   "",
			"version": "",
		},
		"openapi": "3.1.0",
		"paths": map[string]interface{}{
			"/test": map[string]interface{}{
				"get": map[string]interface{} /* len=2 */ {
					"parameters": []interface{} /* len=2 cap=2 */ {
						map[string]interface{} /* len=4 */ {
							"in":       "query",
							"name":     "id",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
						map[string]interface{} /* len=4 */ {
							"in":       "header",
							"name":     "ua",
							"required": true,
							"schema": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"responses": map[string]interface{} /* len=2 */ {
						"200": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": `#/components/schemas/Res01`, /* len=26 */
									},
								},
							},
							"description": "",
						},
						"403": map[string]interface{} /* len=2 */ {
							"content": map[string]interface{}{
								"application/json" /* len=16 */ : map[string]interface{}{
									"schema": map[string]interface{}{
										"$ref": `#/components/schemas/Res02`, /* len=26 */
									},
								},
							},
							"description": "returns 403",
						},
					},
				},
			},
		},
	})
}
