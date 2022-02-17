package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"

	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase/mock"
)

func Test_authController_Login(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		authUC usecase.AuthUseCase
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		want       string
	}{
		{
			name: "successful",
			fields: fields{
				authUC: mock.AuthUseCase{
					OnLogin: func(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error) {
						ret := usecase.AuthTokenOutput{
							AccessToken: "my_access_token",
						}

						return &ret, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`{"cpf":"12345611", "secret":"secret"}`)))
				}(),
			},
			wantStatus: 200,
			want:       `{"access_token": "my_access_token"}`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				authUC: mock.AuthUseCase{
					OnLogin: func(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`{"cpf":"12345678911", "secret": "secret"}`)))
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
		{
			name: "should return 401 when invalid credentials",
			fields: fields{
				authUC: mock.AuthUseCase{
					OnLogin: func(ctx context.Context, loginInput usecase.AuthLoginInput) (*usecase.AuthTokenOutput, error) {
						return nil, usecase.ErrAuthInvalidCredentials
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte(`{"cpf":"12345678911", "secret": "secret"}`)))
				}(),
			},
			wantStatus: 401,
			want:       fmt.Sprintf(`{"code": 401, "message": "%s"}`, usecase.ErrAuthInvalidCredentials),
		},
		{
			name: "should return 400 with error msg when request body is missing",
			fields: fields{
				authUC: mock.AuthUseCase{
					OnLogin: nil,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/login", nil),
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authCtrl := NewAuthController(tt.fields.authUC)

			authCtrl.Login(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("Login() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := rec.Body.String()
			ja.Assertf(bodyStr, tt.want)
		})
	}
}
