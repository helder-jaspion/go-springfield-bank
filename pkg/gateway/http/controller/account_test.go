package controller

import (
	"bytes"
	"context"
	"errors"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase/mock"
	"github.com/kinbiko/jsonassert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_accountController_Create(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		accountUC usecase.AccountUseCase
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
			name: "successful minimum input",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnCreate: func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
						ret := usecase.AccountCreateOutput{
							ID:        "uuid-1",
							Name:      "Bart Simpson",
							CPF:       "123.456.789-11",
							Balance:   0,
							CreatedAt: time.Time{},
						}

						return &ret, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"name":"Bart Simpson", "cpf":"12345611", "secret":"secret"}`)))
				}(),
			},
			wantStatus: 201,
			want:       `{"id": "uuid-1", "name": "Bart Simpson", "cpf": "123.456.789-11", "balance": 0, "created_at": "<<PRESENCE>>"}`,
		},
		{
			name: "successful maximum input",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnCreate: func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
						ret := usecase.AccountCreateOutput{
							ID:        "uuid-1",
							Name:      "Bart Simpson",
							CPF:       "123.456.789-11",
							Balance:   5.96,
							CreatedAt: time.Time{},
						}

						return &ret, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"name":"Bart Simpson", "cpf":"12345611", "balance":5.96, "secret": "secret"}`)))
				}(),
			},
			wantStatus: 201,
			want:       `{"id": "uuid-1", "name": "Bart Simpson", "cpf": "123.456.789-11", "balance": 5.96, "created_at": "<<PRESENCE>>"}`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnCreate: func(ctx context.Context, accountInput usecase.AccountCreateInput) (*usecase.AccountCreateOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader([]byte(`{"name":"Bart Simpson", "cpf":"12345611", "balance":5.96, "secret": "secret"}`)))
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
		{
			name: "should return 400 with error msg when request body is missing",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnCreate: nil,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/accounts", nil),
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAccountController(tt.fields.accountUC)

			a.Create(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("Create() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := string(rec.Body.Bytes())
			ja.Assertf(bodyStr, tt.want)
		})
	}
}

func Test_accountController_Fetch(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		accountUC usecase.AccountUseCase
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
			name: "successful empty result",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnFetch: func(ctx context.Context) ([]usecase.AccountFetchOutput, error) {
						return []usecase.AccountFetchOutput{}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts", nil)
				}(),
			},
			wantStatus: 200,
			want:       `[]`,
		},
		{
			name: "successful one result",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnFetch: func(ctx context.Context) ([]usecase.AccountFetchOutput, error) {
						return []usecase.AccountFetchOutput{
							{
								AccountCreateOutput: usecase.AccountCreateOutput{
									ID:        "uuid-1",
									Name:      "Bart Simpson",
									CPF:       "123.456.789-11",
									Balance:   0,
									CreatedAt: time.Time{},
								},
							},
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts", nil)
				}(),
			},
			wantStatus: 200,
			want:       `[{"id": "uuid-1", "name": "Bart Simpson", "cpf": "123.456.789-11", "balance": 0, "created_at": "<<PRESENCE>>"}]`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnFetch: func(ctx context.Context) ([]usecase.AccountFetchOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts", nil)
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAccountController(tt.fields.accountUC)

			a.Fetch(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("Fetch() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := string(rec.Body.Bytes())
			ja.Assertf(bodyStr, tt.want)
		})
	}
}

func Test_accountController_GetBalance(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		accountUC usecase.AccountUseCase
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
			name: "successful result zero balance",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error) {
						return &usecase.AccountBalanceOutput{
							ID:      "uuid-1",
							Balance: 0,
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts/uuid-1/balance", nil)
				}(),
			},
			wantStatus: 200,
			want:       `{"id":"uuid-1", "balance":0}`,
		},
		{
			name: "successful result positive balance",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error) {
						return &usecase.AccountBalanceOutput{
							ID:      "uuid-1",
							Balance: 10.59,
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts/uuid-1/balance", nil)
				}(),
			},
			wantStatus: 200,
			want:       `{"id":"uuid-1", "balance":10.59}`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				accountUC: mock.AccountUseCase{
					OnGetBalance: func(ctx context.Context, id model.AccountID) (*usecase.AccountBalanceOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/accounts/uuid-1/balance", nil)
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAccountController(tt.fields.accountUC)

			a.GetBalance(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("GetBalance() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := string(rec.Body.Bytes())
			ja.Assertf(bodyStr, tt.want)
		})
	}
}
