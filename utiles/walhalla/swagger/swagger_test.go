package swagger

import (
	"reflect"
	"testing"
)

const testingYaml = `
swagger: '2.0'
info:
  version: "2018-10-23"
  title: Wave Application
host: localhost:9600
basePath: /
tags:
  - name: "user"
schemes:
  - http
consumes:
  - application/json
produces:
  - application/json
paths:
  /user/signup:
    post:
      tags:
      - "user"
      summary: "Creating an account."
      description: ""
      operationId: "signupUser"
      parameters:
        - in: "body"
          name: "body"
          description: "Key user information."
          required: true
          schema:
            $ref: "#/definitions/UserExtended"
      responses:
        201:
          description: "Created"
        403:
          description: "Forbidden"
          schema:
            $ref: "#/definitions/UserExtended"
  
definitions:
  UserCredentials:
    type: object 
    required:
      - "username"
      - "password"
    properties:
      username:
        type: string
        example: florence
      password:
        type: string
        example: pass123

  UserExtended:
    type: object 
    required:
      - "username"
      - "password"
      - "avatarSource"
    properties:
      username:
        type: string
        example: florence
      password:
        type: string
        example: pass123
      avatarSource:
        type: string
        example: <avatarSource>
`

func TestParceSwaggerYaml(t *testing.T) {
	expected := ParsedData{
		Info: Info{
			Version: "2018-10-23",
			Title:   "Wave Application",
		},
		Operations: []Operation{
			{"SignupUser", "user", "UserSignupUserHandler", "SignupUserHandlerFunc"},
		},
		Subcategories: []string{
			"user",
		},
		Sub2Operation: map[string][]string{
			"user": []string{"SignupUser"},
		},
	}

	data := []byte(testingYaml)
	parsed := ParceSwaggerYaml(data)
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("\nTaken:   %s\nExpected:%s", parsed, expected)
	}
}
