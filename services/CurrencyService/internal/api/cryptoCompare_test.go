package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/stretchr/testify/assert"
)

func Test_CryptoCompare_GetRateFromCryptoCompare(t *testing.T) {
	tests := []struct {
		name       string
		from       string
		to         string
		handler    http.HandlerFunc
		wantRate   float64
		wantErr    error
		wantAnyErr bool
	}{
		{
			name: "success",
			from: "USD", to: "RUB",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]float64{"RUB": 85.0})
			},
			wantRate:   85.0,
			wantAnyErr: false,
		},
		{
			name: "error: api returns non-200",
			from: "USD", to: "RUB",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:    domain.ErrFetchFromExternalAPI,
			wantAnyErr: false,
		},
		{
			name: "error: rate not found in response",
			from: "USD", to: "RUB",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]float64{"EUR": 0.9})
			},
			wantErr:    domain.ErrRateNotFound,
			wantAnyErr: false,
		},
		{
			name: "error: invalid json",
			from: "USD", to: "RUB",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("not json"))
			},
			wantAnyErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			c := NewCryptoCompare(
				config.CompareFinAPIURLType(server.URL),
				config.CompareFinAPIKeyType("test-key"),
			)

			rate, err := c.GetRateFromCryptoCompare(context.Background(), tt.from, tt.to)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, 0.0, rate)
			} else if tt.wantAnyErr {
				assert.Error(t, err)
				assert.Equal(t, 0.0, rate)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRate, rate)
			}
		})
	}
}
