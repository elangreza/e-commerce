package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/elangreza/e-commerce/gen"
	"github.com/elangreza/e-commerce/payment/internal/constanta"
	"github.com/elangreza/e-commerce/payment/internal/entity"
	"github.com/elangreza/e-commerce/payment/internal/service"
	"github.com/elangreza/e-commerce/payment/internal/service/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type PaymentServiceTestSuite struct {
	suite.Suite
	ctrl                   *gomock.Controller
	svc                    *service.PaymentService
	mockPaymentRepo        *mock.MockpaymentRepo
	mockOrderServiceClient *mock.MockOrderServiceClient
}

func (s *PaymentServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.mockPaymentRepo = mock.NewMockpaymentRepo(s.ctrl)
	s.mockOrderServiceClient = mock.NewMockOrderServiceClient(s.ctrl)

	s.svc = service.NewPaymentService(
		s.mockPaymentRepo,
		1*time.Second,
		s.mockOrderServiceClient,
	)
}

func (s *PaymentServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestWarehouseServiceSuite(t *testing.T) {
	suite.Run(t, new(PaymentServiceTestSuite))
}

func (s *PaymentServiceTestSuite) TestProcessPayment() {

	tests := []struct {
		name          string
		req           *gen.ProcessPaymentRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Success",
			req: &gen.ProcessPaymentRequest{
				OrderId: "1",
				TotalAmount: &gen.Money{
					Units: 1000,
				},
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					CreatePayment(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.ProcessPayment(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
			}
		})
	}
}

func (s *PaymentServiceTestSuite) TestRollbackPayment() {

	tests := []struct {
		name          string
		req           *gen.RollbackPaymentRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Success",
			req: &gen.RollbackPaymentRequest{
				TransactionId: "aaaa",
				Reason:        "aaaa",
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.WAITING,
					}, nil)

				s.mockPaymentRepo.EXPECT().
					UpdatePaymentStatusByTransactionID(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.RollbackPayment(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
			}
		})
	}
}

func (s *PaymentServiceTestSuite) TestGetPayment() {

	tests := []struct {
		name          string
		req           *gen.GetPaymentRequest
		setupMock     func()
		expectedError string
	}{
		{
			name: "Success",
			req: &gen.GetPaymentRequest{
				TransactionId: "aaaa",
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.WAITING,
					}, nil)

			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.GetPayment(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
			}
		})
	}
}

func (s *PaymentServiceTestSuite) TestUpdatePayment() {

	tests := []struct {
		name          string
		req           *gen.UpdatePaymentRequest
		setupMock     func()
		expectedError string
		expectedResp  *gen.UpdatePaymentResponse
	}{
		{
			name: "Success Paid",
			req: &gen.UpdatePaymentRequest{
				TransactionId: "aaaa",
				TotalAmount: &gen.Money{
					Units:        10000,
					CurrencyCode: "IDR",
				},
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.WAITING,
						TotalAmount: &gen.Money{
							Units:        10000,
							CurrencyCode: "IDR",
						},
					}, nil)
				s.mockPaymentRepo.EXPECT().
					UpdatePaymentStatusByTransactionID(gomock.Any(), constanta.PAID, gomock.Any()).
					Return(nil)

				s.mockOrderServiceClient.EXPECT().
					CallbackTransaction(gomock.Any(), gomock.Any()).
					Return(&gen.Empty{}, nil)

			},
			expectedError: "",
			expectedResp: &gen.UpdatePaymentResponse{
				Status: string(constanta.PAID),
			},
		},
		{
			name: "Success Failed",
			req: &gen.UpdatePaymentRequest{
				TransactionId: "aaaa",
				TotalAmount: &gen.Money{
					Units:        9000,
					CurrencyCode: "IDR",
				},
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.WAITING,
						TotalAmount: &gen.Money{
							Units:        10000,
							CurrencyCode: "IDR",
						},
					}, nil)
				s.mockPaymentRepo.EXPECT().
					UpdatePaymentStatusByTransactionID(gomock.Any(), constanta.FAILED, gomock.Any()).
					Return(nil)

				s.mockOrderServiceClient.EXPECT().
					CallbackTransaction(gomock.Any(), gomock.Any()).
					Return(&gen.Empty{}, nil)

			},
			expectedError: "",
			expectedResp: &gen.UpdatePaymentResponse{
				Status: string(constanta.FAILED),
			},
		},
		{
			name: "Already processed, return early",
			req: &gen.UpdatePaymentRequest{
				TransactionId: "aaaa",
				TotalAmount: &gen.Money{
					Units:        9000,
					CurrencyCode: "IDR",
				},
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.PAID,
						TotalAmount: &gen.Money{
							Units:        10000,
							CurrencyCode: "IDR",
						},
					}, nil)

			},
			expectedError: "",
			expectedResp: &gen.UpdatePaymentResponse{
				Status: string(constanta.PAID),
			},
		},
		{
			name: "currency code not match",
			req: &gen.UpdatePaymentRequest{
				TransactionId: "aaaa",
				TotalAmount: &gen.Money{
					Units:        9000,
					CurrencyCode: "USD",
				},
			},
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetPaymentByTransactionID(gomock.Any(), gomock.Any()).
					Return(&entity.Payment{
						ID:     uuid.New(),
						Status: constanta.WAITING,
						TotalAmount: &gen.Money{
							Units:        10000,
							CurrencyCode: "IDR",
						},
					}, nil)

			},
			expectedError: "currency code not match",
			expectedResp:  nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.UpdatePayment(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedResp)
			}
		})
	}
}

func (s *PaymentServiceTestSuite) TestRemoveExpiredPayment() {
	tests := []struct {
		name          string
		req           time.Duration
		setupMock     func()
		expectedError string
		expectedResp  int
	}{
		{
			name: "Success",
			req:  1 * time.Minute,
			setupMock: func() {
				s.mockPaymentRepo.EXPECT().
					GetExpiredPayments(gomock.Any(), gomock.Any()).
					Return([]entity.Payment{
						{
							ID:     uuid.New(),
							Status: constanta.PAID,
							TotalAmount: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
						},
						{
							ID:     uuid.New(),
							Status: constanta.WAITING,
							TotalAmount: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
						},
						{
							ID:     uuid.New(),
							Status: constanta.WAITING,
							TotalAmount: &gen.Money{
								Units:        10000,
								CurrencyCode: "IDR",
							},
						},
					}, nil)
				s.mockPaymentRepo.EXPECT().
					UpdatePaymentStatusByTransactionID(gomock.Any(), constanta.FAILED, gomock.Any()).
					Return(nil).Times(2)

				s.mockOrderServiceClient.EXPECT().
					CallbackTransaction(gomock.Any(), gomock.Any()).
					Return(&gen.Empty{}, nil).Times(2)

			},
			expectedError: "",
			expectedResp:  3,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.svc.RemoveExpiredPayment(context.Background(), tt.req)

			if tt.expectedError != "" {
				s.Error(err)
				s.Contains(err.Error(), tt.expectedError)
				s.Nil(resp)
			} else {
				s.NoError(err)
				s.NotNil(resp)
				s.Equal(resp, tt.expectedResp)
			}
		})
	}
}
