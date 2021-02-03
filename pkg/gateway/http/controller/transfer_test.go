package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/helder-jaspion/go-springfield-bank/pkg/appcontext"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/model"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase"
	"github.com/helder-jaspion/go-springfield-bank/pkg/domain/usecase/mock"
	"github.com/kinbiko/jsonassert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_transferController_Create(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		trfUC  usecase.TransferUseCase
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
			name: "successful minimum input",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnCreate: func(ctx context.Context, transferInput usecase.TransferCreateInput) (*usecase.TransferCreateOutput, error) {
						ret := usecase.TransferCreateOutput{
							ID:                   "trf-uuid-1",
							AccountOriginID:      "uuid-1",
							AccountDestinationID: "uuid-2",
							Amount:               1,
							CreatedAt:            time.Time{},
						}

						return &ret, nil
					},
				},
				authUC: mock.AuthUseCase{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader([]byte(`{"account_destination_id":"uuid-2", "amount": 1}`)))

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 201,
			want:       `{"id": "trf-uuid-1", "account_origin_id":"uuid-1", "account_destination_id":"uuid-2", "amount": 1, "created_at": "<<PRESENCE>>"}`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnCreate: func(ctx context.Context, transferInput usecase.TransferCreateInput) (*usecase.TransferCreateOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader([]byte(`{"account_destination_id":"uuid-2", "amount": 1}`)))

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
		{
			name: "should return 400 with error msg when request body is missing",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnCreate: nil,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/transfers", nil)

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 400,
			want:       `{"code": 400, "message": "error reading input"}`,
		},
		{
			name: "should return 400 when destination is invalid",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnCreate: func(ctx context.Context, transferInput usecase.TransferCreateInput) (*usecase.TransferCreateOutput, error) {
						return nil, usecase.ErrTransferDestinationAccountRequired
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader([]byte(`{"account_destination_id":"", "amount": 1}`)))

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 400,
			want:       fmt.Sprintf(`{"code": 400, "message": "%s"}`, usecase.ErrTransferDestinationAccountRequired),
		},
		{
			name: "should return 401 when invalid token",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnCreate: nil,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewReader([]byte(`{"account_destination_id":"", "amount": 1}`)))
				}(),
			},
			wantStatus: 401,
			want:       fmt.Sprintf(`{"code": 401, "message": "%s"}`, usecase.ErrAuthInvalidAccessToken),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trfCtrl := NewTransferController(tt.fields.trfUC, tt.fields.authUC)

			trfCtrl.Create(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("Create() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := rec.Body.String()
			ja.Assertf(bodyStr, tt.want)
		})
	}
}

func Test_transferController_Fetch(t *testing.T) {
	t.Parallel()

	ja := jsonassert.New(t)

	type fields struct {
		trfUC  usecase.TransferUseCase
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
			name: "successful empty result",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]usecase.TransferFetchOutput, error) {
						return []usecase.TransferFetchOutput{}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/transfers", nil)

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 200,
			want:       `[]`,
		},
		{
			name: "successful one result",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]usecase.TransferFetchOutput, error) {
						return []usecase.TransferFetchOutput{
							{
								TransferCreateOutput: usecase.TransferCreateOutput{
									ID:                   "trf-uuid-1",
									AccountOriginID:      "uuid-1",
									AccountDestinationID: "uuid-2",
									Amount:               1,
									CreatedAt:            time.Time{},
								},
							},
						}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/transfers", nil)

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 200,
			want:       `[{"id": "trf-uuid-1", "account_origin_id": "uuid-1","account_destination_id": "uuid-2","amount": 1, "created_at": "<<PRESENCE>>"}]`,
		},
		{
			name: "should return 500 when usecase error",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnFetch: func(ctx context.Context, accountID model.AccountID) ([]usecase.TransferFetchOutput, error) {
						return nil, errors.New("any error")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req := httptest.NewRequest(http.MethodGet, "/transfers", nil)

					return req.WithContext(appcontext.WithAuthSubject(req.Context(), "uuid-1"))
				}(),
			},
			wantStatus: 500,
			want:       `{"code": 500, "message": "any error"}`,
		},
		{
			name: "should return 401 when invalid token error",
			fields: fields{
				trfUC: mock.TransferUseCase{
					OnFetch: nil,
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					return httptest.NewRequest(http.MethodGet, "/transfers", nil)
				}(),
			},
			wantStatus: 401,
			want:       fmt.Sprintf(`{"code": 401, "message": "%s"}`, usecase.ErrAuthInvalidAccessToken),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trfCtrl := NewTransferController(tt.fields.trfUC, tt.fields.authUC)

			trfCtrl.Fetch(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			// Check the response status code
			if statusCode := rec.Code; statusCode != tt.wantStatus {
				t.Errorf("Fetch() statusCode = %v, wantStatus %v", statusCode, tt.wantStatus)
			}

			// Check result response
			bodyStr := rec.Body.String()
			ja.Assertf(bodyStr, tt.want)
		})
	}
}
