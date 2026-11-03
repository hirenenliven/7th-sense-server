// Command apidocs finishes the API docs that `go generate ./docs/api` produces: it converts the
// Swagger 2.0 spec swag writes (docs/api/.build/swagger.json) to OpenAPI 3.0, sets the session
// cookie auth and the servers, writes docs/api/openapi.yaml, and builds the Postman collection and
// environment in docs/postman from it (folders = tags, requests named like the operations).
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	project       = "7th-sense"
	sessionCookie = "ts_session"
	localBaseURL  = "http://localhost:8080/api"
)

type object = map[string]any

func main() {
	root := moduleRoot()
	raw, err := os.ReadFile(filepath.Join(root, "docs/api/.build/swagger.json"))
	if err != nil {
		log.Fatalf("read swag output (run `go generate ./docs/api`): %v", err)
	}
	var swagger object
	if err := json.Unmarshal(raw, &swagger); err != nil {
		log.Fatal(err)
	}
	spec := dropNulls(toOpenAPI3(swagger)).(object)
	patchSpec(spec)
	addExamples(spec)
	writeYAML(filepath.Join(root, "docs/api/openapi.yaml"), spec)
	collection, environment := postman(spec)
	writeJSON(filepath.Join(root, "docs/postman", project+".postman_collection.json"), collection)
	writeJSON(filepath.Join(root, "docs/postman", project+".postman_environment.json"), environment)
	fmt.Println("wrote docs/api/openapi.yaml and docs/postman/*")
}

func moduleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("go.mod not found")
		}
		dir = parent
	}
}

// toOpenAPI3 converts the parts of Swagger 2.0 swag emits: definitions, body / formData / query /
// path parameters and responses.
func toOpenAPI3(swagger object) object {
	data, _ := json.Marshal(swagger)
	data = []byte(strings.ReplaceAll(string(data), `"#/definitions/`, `"#/components/schemas/`))
	_ = json.Unmarshal(data, &swagger)

	paths := object{}
	for path, value := range swagger["paths"].(object) {
		item := object{}
		for method, raw := range value.(object) {
			item[method] = convertOperation(raw.(object))
		}
		paths[path] = item
	}
	return object{
		"openapi":    "3.0.3",
		"info":       swagger["info"],
		"paths":      paths,
		"components": object{"schemas": swagger["definitions"]},
	}
}

func convertOperation(operation object) object {
	consumes := strings_(operation["consumes"])
	produces := strings_(operation["produces"])
	if len(produces) == 0 {
		produces = []string{"application/json"}
	}
	var parameters []any
	form := object{}
	var formRequired []any
	for _, raw := range list(operation["parameters"]) {
		parameter := raw.(object)
		switch parameter["in"] {
		case "body":
			operation["requestBody"] = object{
				"required":    parameter["required"],
				"description": parameter["description"],
				"content":     object{"application/json": object{"schema": parameter["schema"]}},
			}
		case "formData":
			schema := object{"type": parameter["type"], "description": parameter["description"]}
			if parameter["type"] == "file" {
				schema = object{"type": "string", "format": "binary", "description": parameter["description"]}
			}
			form[parameter["name"].(string)] = schema
			if required, _ := parameter["required"].(bool); required {
				formRequired = append(formRequired, parameter["name"])
			}
		default:
			schema := object{"type": parameter["type"]}
			if format, ok := parameter["format"]; ok {
				schema["format"] = format
			}
			parameters = append(parameters, object{
				"name": parameter["name"], "in": parameter["in"], "required": parameter["required"],
				"description": parameter["description"], "schema": schema,
			})
		}
	}
	if len(form) > 0 {
		media := "multipart/form-data"
		if len(consumes) > 0 {
			media = consumes[0]
		}
		schema := object{"type": "object", "properties": form}
		if len(formRequired) > 0 {
			schema["required"] = formRequired
		}
		operation["requestBody"] = object{"required": true, "content": object{media: object{"schema": schema}}}
	}
	if parameters != nil {
		operation["parameters"] = parameters
	} else {
		delete(operation, "parameters")
	}
	responses := object{}
	for code, raw := range operation["responses"].(object) {
		response := raw.(object)
		converted := object{"description": response["description"]}
		if converted["description"] == nil || converted["description"] == "" {
			converted["description"] = http.StatusText(atoi(code))
		}
		if schema, ok := response["schema"].(object); ok {
			if schema["type"] == "file" {
				converted["content"] = object{"application/octet-stream": object{"schema": object{"type": "string", "format": "binary"}}}
			} else {
				converted["content"] = object{produces[0]: object{"schema": schema}}
			}
		}
		responses[code] = converted
	}
	operation["responses"] = responses
	delete(operation, "consumes")
	delete(operation, "produces")
	return operation
}

// dropNulls removes null values (Swagger 2.0 fields with no OpenAPI 3 counterpart).
func dropNulls(value any) any {
	switch typed := value.(type) {
	case object:
		for key, item := range typed {
			if item == nil {
				delete(typed, key)
				continue
			}
			typed[key] = dropNulls(item)
		}
	case []any:
		for index, item := range typed {
			typed[index] = dropNulls(item)
		}
	}
	return value
}

// addExamples gives every JSON request body and response an example built from its schema: the
// envelope with sample data on success, {success: false, message} on errors.
func addExamples(spec object) {
	forEachOperation(spec, func(_, _ string, operation object) {
		if content := requestContent(spec, operation); content != nil {
			if media, ok := content["application/json"].(object); ok && media["example"] == nil {
				media["example"] = exampleOf(spec, media["schema"], 0)
			}
		}
		code, data := successSchema(spec, operation)
		for status, raw := range operation["responses"].(object) {
			response := raw.(object)
			content, _ := response["content"].(object)
			media, ok := content["application/json"].(object)
			if !ok {
				continue
			}
			number := atoi(status)
			example := object{"statusCode": number, "success": number < 400, "message": response["description"]}
			if status == code {
				example["message"] = map[bool]string{true: "created", false: "success"}[number == http.StatusCreated]
				if sample := exampleOf(spec, data, 0); sample != nil {
					example["data"] = sample
				}
			}
			media["example"] = example
		}
	})
}

func atoi(text string) int {
	value := 0
	for _, char := range text {
		value = value*10 + int(char-'0')
	}
	return value
}

// patchSpec fixes what Swagger 2.0 (swag's output) cannot say: the cookie auth and the servers.
func patchSpec(spec object) {
	spec["servers"] = []any{
		object{"url": "/api", "description": "This server"},
		object{"url": localBaseURL, "description": "Local development"},
	}
	components, _ := spec["components"].(object)
	if components == nil {
		components = object{}
		spec["components"] = components
	}
	components["securitySchemes"] = object{"cookieAuth": object{
		"type": "apiKey", "in": "cookie", "name": sessionCookie,
		"description": "HttpOnly session cookie set by POST /v1/auth/login and cleared by POST /v1/auth/logout.",
	}}
	// Tags in the order the paths use them, each described once
	var tags []string
	forEachOperation(spec, func(_, _ string, operation object) {
		for _, tag := range strings_(operation["tags"]) {
			if !slices.Contains(tags, tag) {
				tags = append(tags, tag)
			}
		}
	})
	sort.Strings(tags)
	tagList := []any{}
	for _, tag := range tags {
		tagList = append(tagList, object{"name": tag})
	}
	spec["tags"] = tagList
}

func forEachOperation(spec object, visit func(path, method string, operation object)) {
	paths, _ := spec["paths"].(object)
	keys := make([]string, 0, len(paths))
	for path := range paths {
		keys = append(keys, path)
	}
	sort.Strings(keys)
	for _, path := range keys {
		item, _ := paths[path].(object)
		for _, method := range []string{"get", "post", "put", "delete"} {
			if operation, ok := item[method].(object); ok {
				visit(path, method, operation)
			}
		}
	}
}

func strings_(value any) []string {
	list, _ := value.([]any)
	out := make([]string, 0, len(list))
	for _, item := range list {
		if text, ok := item.(string); ok {
			out = append(out, text)
		}
	}
	return out
}

// ---------------------------------------------------------------- Postman

type request struct {
	path, method string
	operation    object
	order        int
}

// Runner order inside a folder: list → create → get → update → delete
var methodOrder = map[string]int{"get": 0, "post": 1, "put": 3, "delete": 4}

func postman(spec object) (object, object) {
	folders := map[string][]request{}
	forEachOperation(spec, func(path, method string, operation object) {
		tag := strings_(operation["tags"])[0]
		order := methodOrder[method]
		if method == "get" && strings.HasSuffix(path, "}") {
			order = 2
		}
		folders[tag] = append(folders[tag], request{path: path, method: method, operation: operation, order: order})
	})
	idVariables := map[string]bool{}
	var items []any
	items = append(items, authFolder(spec))
	tags := make([]string, 0, len(folders))
	for tag := range folders {
		if tag != "Auth" {
			tags = append(tags, tag)
		}
	}
	sort.Strings(tags)
	for _, tag := range tags {
		list := folders[tag]
		sort.SliceStable(list, func(i, j int) bool {
			if resource(list[i].path) != resource(list[j].path) {
				return resource(list[i].path) < resource(list[j].path)
			}
			return list[i].order < list[j].order
		})
		var folderItems []any
		for _, req := range list {
			folderItems = append(folderItems, postmanItem(spec, req, idVariables))
		}
		items = append(items, object{"name": tag, "item": folderItems})
	}
	items = append(items, negativeFolder())
	items = append(items, logoutFolder())

	collection := object{
		"info": object{
			"name":        "7th Sense API",
			"description": "Generated from docs/api/openapi.yaml by `go generate ./docs/api` — do not edit by hand.\n\nAuth: run **Auth › Login** first; the server sets the HttpOnly `ts_session` cookie and Postman's cookie jar sends it on every later request (there is no bearer token). Requests that create rows save their ids (e.g. `shiftId`) for the get / update / delete requests after them. Running the whole collection writes to the database: use a development database.",
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"auth": object{"type": "noauth"},
		"item": items,
	}

	values := []any{
		object{"key": "baseUrl", "value": localBaseURL, "type": "default", "enabled": true},
		object{"key": "username", "value": "superadmin", "type": "default", "enabled": true},
		object{"key": "password", "value": "change-me", "type": "secret", "enabled": true},
	}
	names := make([]string, 0, len(idVariables))
	for name := range idVariables {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		values = append(values, object{"key": name, "value": "1", "type": "default", "enabled": true})
	}
	environment := object{
		"name":                    "7th Sense — local",
		"values":                  values,
		"_postman_variable_scope": "environment",
	}
	return collection, environment
}

var paramPattern = regexp.MustCompile(`\{(\w+)\}`)

// resource is the first path segment after the version ("/v1/user-cell/{id}" → "user-cell").
func resource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 1 {
		return parts[1]
	}
	return path
}

// idVariable names the environment variable of a path parameter: {id} of /shifts → shiftId.
func idVariable(path, param string) string {
	if param != "id" {
		return param
	}
	name := resource(path)
	switch {
	case strings.HasSuffix(name, "ches"), strings.HasSuffix(name, "shes"):
		name = strings.TrimSuffix(name, "es")
	case strings.HasSuffix(name, "s"):
		name = strings.TrimSuffix(name, "s")
	}
	words := strings.Split(name, "-")
	for index := 1; index < len(words); index++ {
		words[index] = strings.ToUpper(words[index][:1]) + words[index][1:]
	}
	return strings.Join(words, "") + "Id"
}

func postmanItem(spec object, req request, idVariables map[string]bool) object {
	operation := req.operation
	path := paramPattern.ReplaceAllStringFunc(req.path, func(match string) string {
		name := idVariable(req.path, match[1:len(match)-1])
		idVariables[name] = true
		return "{{" + name + "}}"
	})
	var query []any
	for _, param := range list(operation["parameters"]) {
		parameter := param.(object)
		if parameter["in"] != "query" {
			continue
		}
		required, _ := parameter["required"].(bool)
		query = append(query, object{"key": parameter["name"], "value": "", "disabled": !required, "description": parameter["description"]})
	}
	url := object{
		"raw":  "{{baseUrl}}" + path,
		"host": []any{"{{baseUrl}}"},
		"path": toAny(strings.Split(strings.Trim(path, "/"), "/")),
	}
	if len(query) > 0 {
		url["query"] = query
	}
	body := object(nil)
	header := []any{}
	if content := requestContent(spec, operation); content != nil {
		if schema, ok := content["application/json"].(object); ok {
			example, _ := json.MarshalIndent(exampleOf(spec, schema["schema"], 0), "", "  ")
			body = object{"mode": "raw", "raw": string(example), "options": object{"raw": object{"language": "json"}}}
			header = append(header, object{"key": "Content-Type", "value": "application/json"})
		} else if _, ok := content["multipart/form-data"]; ok {
			body = object{"mode": "formdata", "formdata": []any{object{"key": "file", "type": "file", "src": ""}}}
		}
	}
	item := object{
		"name": operation["summary"],
		"request": object{
			"method":      strings.ToUpper(req.method),
			"header":      header,
			"url":         url,
			"description": operation["description"],
		},
		"event": []any{testEvent(spec, req, operation)},
	}
	if body != nil {
		item["request"].(object)["body"] = body
	}
	return item
}

func requestContent(spec object, operation object) object {
	requestBody, _ := operation["requestBody"].(object)
	if requestBody == nil {
		return nil
	}
	requestBody = resolve(spec, requestBody)
	content, _ := requestBody["content"].(object)
	return content
}

func list(value any) []any {
	items, _ := value.([]any)
	return items
}

func toAny(values []string) []any {
	out := make([]any, len(values))
	for index, value := range values {
		out[index] = value
	}
	return out
}

func resolve(spec object, value object) object {
	ref, ok := value["$ref"].(string)
	if !ok {
		return value
	}
	node := any(spec)
	for _, part := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		node = node.(object)[part]
	}
	resolved, _ := node.(object)
	return resolved
}

// exampleOf builds a sample value from a schema (its example when it has one).
func exampleOf(spec object, value any, depth int) any {
	schema, _ := value.(object)
	if schema == nil || depth > 4 {
		return nil
	}
	schema = resolve(spec, schema)
	if example, ok := schema["example"]; ok {
		return example
	}
	if all := list(schema["allOf"]); len(all) > 0 {
		merged := object{}
		for _, part := range all {
			if sample, ok := exampleOf(spec, part, depth+1).(object); ok {
				for key, item := range sample {
					merged[key] = item
				}
			}
		}
		return merged
	}
	switch schema["type"] {
	case "object":
		sample := object{}
		properties, _ := schema["properties"].(object)
		for name, property := range properties {
			if value := exampleOf(spec, property, depth+1); value != nil {
				sample[name] = value
			}
		}
		return sample
	case "array":
		if item := exampleOf(spec, schema["items"], depth+1); item != nil {
			return []any{item}
		}
		return []any{}
	case "integer", "number":
		return 1
	case "boolean":
		return false
	case "string":
		if format, _ := schema["format"].(string); format == "date-time" {
			return "2026-01-01T08:00:00Z"
		}
		return "Postman sample"
	}
	return nil
}

// successSchema returns the success status and the schema of the envelope's data.
func successSchema(spec object, operation object) (string, object) {
	responses, _ := operation["responses"].(object)
	for _, code := range []string{"200", "201"} {
		response, ok := responses[code].(object)
		if !ok {
			continue
		}
		response = resolve(spec, response)
		content, _ := response["content"].(object)
		media, _ := content["application/json"].(object)
		schema, _ := media["schema"].(object)
		for _, part := range list(schema["allOf"]) {
			properties, _ := resolve(spec, part.(object))["properties"].(object)
			if data, ok := properties["data"].(object); ok {
				return code, resolve(spec, data)
			}
		}
		return code, nil
	}
	return "200", nil
}

func testEvent(spec object, req request, operation object) object {
	code, data := successSchema(spec, operation)
	lines := []string{
		fmt.Sprintf("pm.test('status is %s', () => pm.response.to.have.status(%s));", code, code),
		"pm.test('response time < 2000ms', () => pm.expect(pm.response.responseTime).to.be.below(2000));",
	}
	if strings.Contains(fmt.Sprint(operation["responses"]), "octet-stream") {
		return object{"listen": "test", "script": object{"type": "text/javascript", "exec": toAny(lines)}}
	}
	lines = append(lines,
		"const json = pm.response.json();",
		"pm.test('envelope fields', () => {",
		"  pm.expect(json).to.have.property('success', true);",
		"  pm.expect(json.statusCode).to.be.a('number');",
		"  pm.expect(json.message).to.be.a('string');",
		"});",
	)
	if data != nil {
		switch data["type"] {
		case "array":
			lines = append(lines, "pm.test('data is an array', () => pm.expect(json.data).to.be.an('array'));")
		case "object", nil:
			properties, _ := data["properties"].(object)
			keys := make([]string, 0, len(properties))
			for key := range properties {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			lines = append(lines, "pm.test('data is an object', () => pm.expect(json.data).to.be.an('object'));")
			if slices.Contains(keys, "id") {
				lines = append(lines, "pm.test('data.id is a number', () => pm.expect(json.data.id).to.be.a('number'));")
			}
			if len(keys) > 0 {
				shown := keys[:min(len(keys), 6)]
				quoted := make([]string, len(shown))
				for index, key := range shown {
					quoted[index] = "'" + key + "'"
				}
				lines = append(lines, fmt.Sprintf("pm.test('data has its fields', () => pm.expect(json.data).to.include.keys(%s));", strings.Join(quoted, ", ")))
			}
		}
	}
	if req.method == "post" && data != nil {
		if properties, _ := data["properties"].(object); properties["id"] != nil {
			lines = append(lines, fmt.Sprintf("if (json.data && json.data.id) pm.environment.set('%s', json.data.id);", idVariable(req.path+"/{id}", "id")))
		}
	}
	return object{"listen": "test", "script": object{"type": "text/javascript", "exec": toAny(lines)}}
}

func rawRequest(name, method, path, body string, tests []string) object {
	req := object{
		"method": method,
		"header": []any{object{"key": "Content-Type", "value": "application/json"}},
		"url":    object{"raw": "{{baseUrl}}" + path, "host": []any{"{{baseUrl}}"}, "path": toAny(strings.Split(strings.Trim(path, "/"), "/"))},
	}
	if body != "" {
		req["body"] = object{"mode": "raw", "raw": body, "options": object{"raw": object{"language": "json"}}}
	}
	return object{"name": name, "request": req, "event": []any{object{"listen": "test", "script": object{"type": "text/javascript", "exec": toAny(tests)}}}}
}

func statusTests(code int) []string {
	return []string{
		fmt.Sprintf("pm.test('status is %d', () => pm.response.to.have.status(%d));", code, code),
		"pm.test('response time < 2000ms', () => pm.expect(pm.response.responseTime).to.be.below(2000));",
		"pm.test('error envelope', () => { const json = pm.response.json(); pm.expect(json.success).to.equal(false); pm.expect(json.message).to.be.a('string'); });",
	}
}

func authFolder(spec object) object {
	login := statusTests(200)[:2]
	login = append(login,
		"const json = pm.response.json();",
		"pm.test('returns the user', () => { pm.expect(json.success).to.equal(true); pm.expect(json.data).to.be.an('object'); pm.expect(json.data.id).to.be.a('number'); });",
		"pm.test('sets the session cookie', () => pm.expect(pm.cookies.has('"+sessionCookie+"')).to.be.true);",
		"// No token to store: Postman's cookie jar sends ts_session on the next requests.",
		"if (json.data) pm.environment.set('currentUserId', json.data.id);",
	)
	return object{"name": "Auth", "item": []any{
		rawRequest("Login with a wrong password → 401", "POST", "/v1/auth/login", `{"username": "{{username}}", "password": "wrong-password", "remember": false}`, statusTests(401)),
		rawRequest("Login without a username → 400", "POST", "/v1/auth/login", `{"username": "", "password": ""}`, statusTests(400)),
		rawRequest("Login", "POST", "/v1/auth/login", `{"username": "{{username}}", "password": "{{password}}", "remember": false}`, login),
	}}
}

func negativeFolder() object {
	return object{"name": "Negative tests", "item": []any{
		rawRequest("Create shift with an empty body → 400", "POST", "/v1/shifts", `{}`, statusTests(400)),
		rawRequest("Get a shift that does not exist → 404", "GET", "/v1/shifts/999999999", "", statusTests(404)),
	}}
}

func logoutFolder() object {
	return object{"name": "Logout", "item": []any{
		rawRequest("Logout", "POST", "/v1/auth/logout", "", []string{
			"pm.test('status is 200', () => pm.response.to.have.status(200));",
			"pm.test('response time < 2000ms', () => pm.expect(pm.response.responseTime).to.be.below(2000));",
		}),
		rawRequest("Call without a session → 401", "GET", "/v1/machines", "", statusTests(401)),
	}}
}

// ---------------------------------------------------------------- output

func writeJSON(path string, value any) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
}

// writeYAML writes value as block-style YAML, keeping the key order of its JSON form.
func writeYAML(path string, value any) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Fatal(err)
	}
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		log.Fatal(err)
	}
	blockStyle(&node)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	if err := encoder.Encode(&node); err != nil {
		log.Fatal(err)
	}
}

func blockStyle(node *yaml.Node) {
	if node.Kind == yaml.MappingNode || node.Kind == yaml.SequenceNode {
		node.Style = 0
	}
	if node.Kind == yaml.ScalarNode && node.Style == yaml.DoubleQuotedStyle && node.Tag == "!!str" {
		node.Style = 0
	}
	for _, child := range node.Content {
		blockStyle(child)
	}
}
