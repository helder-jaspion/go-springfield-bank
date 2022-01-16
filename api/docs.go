// Package api GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package api

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Helder Alves",
            "url": "https://github.com/helder-jaspion/go-springfield-bank/"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/helder-jaspion/go-springfield-bank/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/account/{id}/balance": {
            "get": {
                "description": "Get the balance of an account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Get account balance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.AccountBalanceOutput"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            }
        },
        "/accounts": {
            "get": {
                "description": "Fetch all the accounts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Fetch accounts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/usecase.AccountFetchOutput"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            },
            "post": {
                "description": "Creates a new account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Create account",
                "parameters": [
                    {
                        "description": "Account",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.AccountCreateInput"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Idempotency key",
                        "name": "X-Idempotency-Key",
                        "in": "header"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/usecase.AccountCreateOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Authenticates the user/account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Credentials",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.AuthLoginInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/usecase.AuthTokenOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            }
        },
        "/transfers": {
            "get": {
                "security": [
                    {
                        "Access token": []
                    }
                ],
                "description": "Fetch the transfers the current account is related to",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transfers"
                ],
                "summary": "Fetch transfers",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/usecase.TransferFetchOutput"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Access token": []
                    }
                ],
                "description": "Creates a new transfer from the current account to another. Debits the amount from origin account and credit it to the destination account.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Transfers"
                ],
                "summary": "Create transfer",
                "parameters": [
                    {
                        "description": "Transfer",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/usecase.TransferCreateInput"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Idempotency key",
                        "name": "X-Idempotency-Key",
                        "in": "header"
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/usecase.TransferCreateOutput"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/io.ErrorOutput"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "io.ErrorOutput": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string",
                    "example": "something wrong happened"
                }
            }
        },
        "usecase.AccountBalanceOutput": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 9999.99
                },
                "id": {
                    "type": "string",
                    "example": "16b1d860-43d3-4970-bb54-ec395908599a"
                }
            }
        },
        "usecase.AccountCreateInput": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "default": 0,
                    "example": 9999.99
                },
                "cpf": {
                    "type": "string",
                    "example": "999.999.999-99"
                },
                "name": {
                    "type": "string",
                    "example": "Bart Simpson"
                },
                "secret": {
                    "type": "string",
                    "example": "S3cr3t"
                }
            }
        },
        "usecase.AccountCreateOutput": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 9999.99
                },
                "cpf": {
                    "type": "string",
                    "example": "999.999.999-99"
                },
                "created_at": {
                    "type": "string",
                    "example": "2020-12-31T23:59:59.999999-03:00"
                },
                "id": {
                    "type": "string",
                    "example": "16b1d860-43d3-4970-bb54-ec395908599a"
                },
                "name": {
                    "type": "string",
                    "example": "Bart Simpson"
                }
            }
        },
        "usecase.AccountFetchOutput": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "number",
                    "example": 9999.99
                },
                "cpf": {
                    "type": "string",
                    "example": "999.999.999-99"
                },
                "created_at": {
                    "type": "string",
                    "example": "2020-12-31T23:59:59.999999-03:00"
                },
                "id": {
                    "type": "string",
                    "example": "16b1d860-43d3-4970-bb54-ec395908599a"
                },
                "name": {
                    "type": "string",
                    "example": "Bart Simpson"
                }
            }
        },
        "usecase.AuthLoginInput": {
            "type": "object",
            "properties": {
                "cpf": {
                    "type": "string",
                    "example": "999.999.999-99"
                },
                "secret": {
                    "type": "string",
                    "example": "S3cr3t"
                }
            }
        },
        "usecase.AuthTokenOutput": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTI1Mjk3NTksImlhdCI6MTYxMjUyODg1OSwic3ViIjoiNmMzYjhhNTUtNmI4MC00MTM3LTlkZmYtNTAzY2FmNTc2NTE0In0.uJc24K1GLMwbswa_4DqCpNzqhdLYAklVyGaZEcmlKn8"
                }
            }
        },
        "usecase.TransferCreateInput": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string",
                    "example": "ce8ba94a-2c5f-4e00-80a1-6fcb0ce7382d"
                },
                "amount": {
                    "type": "number",
                    "example": 9999.99
                }
            }
        },
        "usecase.TransferCreateOutput": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string",
                    "example": "ce8ba94a-2c5f-4e00-80a1-6fcb0ce7382d"
                },
                "account_origin_id": {
                    "type": "string",
                    "example": "16b1d860-43d3-4970-bb54-ec395908599a"
                },
                "amount": {
                    "type": "number",
                    "example": 9999.99
                },
                "created_at": {
                    "type": "string",
                    "example": "2020-12-31T23:59:59.999999-03:00"
                },
                "id": {
                    "type": "string",
                    "example": "e82706ef-9ffb-45a2-8081-547accd818c4"
                }
            }
        },
        "usecase.TransferFetchOutput": {
            "type": "object",
            "properties": {
                "account_destination_id": {
                    "type": "string",
                    "example": "ce8ba94a-2c5f-4e00-80a1-6fcb0ce7382d"
                },
                "account_origin_id": {
                    "type": "string",
                    "example": "16b1d860-43d3-4970-bb54-ec395908599a"
                },
                "amount": {
                    "type": "number",
                    "example": 9999.99
                },
                "created_at": {
                    "type": "string",
                    "example": "2020-12-31T23:59:59.999999-03:00"
                },
                "id": {
                    "type": "string",
                    "example": "e82706ef-9ffb-45a2-8081-547accd818c4"
                }
            }
        }
    },
    "securityDefinitions": {
        "Access token": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "GO Springfield Bank API",
	Description: "GO Springfield Bank API simulates a digital bank where you can create and fetch accounts, login with your account and transfer money to other accounts.\n### Authorization\nYou can get the access_token returned from `/login`, click the **Authorize** button and input this format `Bearer <access_token>`. After this, the `Authorization` header will be sent along in your next requests.\nThe JWT access token has short expiration, so maybe you have to log in again to get a new `access_token`.\n### X-Idempotency-Key\nIf you send the `X-Idempotency-Key` header along with a request, that request's response will be cached. So, if you send the same request with the same `X-Idempotency-Key` again, the server will respond the cached response, so no processing will be done twice.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
