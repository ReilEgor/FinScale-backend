package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	mocks "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/mocks/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestRouter(t *testing.T, uc *mocks.CurrencyUseCase) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	h := NewHandler(uc)
	r := gin.New()
	r.POST("/convert", h.ConvertCurrency)
	return r
}

func Test_Handler_ConvertCurrency(t *testing.T) {
	tests := []struct {
		name       string
		body       any
		mockSetup  func(uc *mocks.CurrencyUseCase)
		wantStatus int
	}{
		{
			name: "success",
			body: ConvertCurrencyRequest{From: "USD", To: "RUB", Amount: 1.0},
			mockSetup: func(uc *mocks.CurrencyUseCase) {
				uc.On("ConvertCurrency", mock.Anything, "USD", "RUB", 1.0).Return(85.0, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "error: invalid json",
			body:       "not json",
			mockSetup:  func(uc *mocks.CurrencyUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error: amount <= 0",
			body:       ConvertCurrencyRequest{From: "USD", To: "RUB", Amount: -1.0},
			mockSetup:  func(uc *mocks.CurrencyUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error: from contains digits",
			body:       ConvertCurrencyRequest{From: "US1", To: "RUB", Amount: 1.0},
			mockSetup:  func(uc *mocks.CurrencyUseCase) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "error: rate not found",
			body: ConvertCurrencyRequest{From: "USD", To: "RUB", Amount: 1.0},
			mockSetup: func(uc *mocks.CurrencyUseCase) {
				uc.On("ConvertCurrency", mock.Anything, "USD", "RUB", 1.0).Return(0.0, domain.ErrRateNotFound).Once()
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name: "error: external API unavailable",
			body: ConvertCurrencyRequest{From: "USD", To: "RUB", Amount: 1.0},
			mockSetup: func(uc *mocks.CurrencyUseCase) {
				uc.On("ConvertCurrency", mock.Anything, "USD", "RUB", 1.0).Return(0.0, domain.ErrFetchFromExternalAPI).Once()
			},
			wantStatus: http.StatusBadGateway,
		},
		{
			name: "error: internal server error",
			body: ConvertCurrencyRequest{From: "USD", To: "RUB", Amount: 1.0},
			mockSetup: func(uc *mocks.CurrencyUseCase) {
				uc.On("ConvertCurrency", mock.Anything, "USD", "RUB", 1.0).Return(0.0, errors.New("unexpected")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ucMock := mocks.NewCurrencyUseCase(t)
			tt.mockSetup(ucMock)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/convert", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			newTestRouter(t, ucMock).ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
