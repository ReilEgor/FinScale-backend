package usecase

import (
	"context"

	"testing"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	mocks "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/mocks/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UseCase_ConvertCurrency(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		amount    float64
		mockSetup func(repo *mocks.CurrencyRepository, fetcher *mocks.CurrencyFetcher)
		expected  float64
		wantErr   error
	}{
		{
			name:   "success: rate found in cache",
			from:   "USD",
			to:     "RUB",
			amount: 1.0,
			mockSetup: func(repo *mocks.CurrencyRepository, fetcher *mocks.CurrencyFetcher) {
				repo.On("GetCurrency", mock.Anything, "USD", "RUB").Return(85.0, nil).Once()
			},
			expected: 85.0,
			wantErr:  nil,
		},
		{
			name:   "success: rate not in cache, saved after fetch",
			from:   "USD",
			to:     "RUB",
			amount: 1.0,
			mockSetup: func(repo *mocks.CurrencyRepository, fetcher *mocks.CurrencyFetcher) {
				fetcher.On("GetRateFromCryptoCompare", mock.Anything, "USD", "RUB").
					Return(85.0, nil).Once()
				repo.On("GetCurrency", mock.Anything, "USD", "RUB").Return(0.0, domain.ErrCurrencyRateNotFound).Once()
				repo.On("SaveCurrency", mock.Anything, "USD", "RUB", mock.AnythingOfType("float64")).Return(nil).Once()
			},
			expected: 85.0,
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := mocks.NewCurrencyRepository(t)
			fetcherMock := mocks.NewCurrencyFetcher(t)
			tt.mockSetup(repositoryMock, fetcherMock)
			uc := newTestUseCase(t, repositoryMock, fetcherMock)
			result, err := uc.ConvertCurrency(context.Background(), tt.from, tt.to, tt.amount)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, 0.0, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func Test_UseCase_SaveCurrency(t *testing.T) {

}
func Test_UseCase_GetRateFromCryptoCompare(t *testing.T) {

}

func newTestUseCase(t *testing.T, repo domain.CurrencyRepository, fetcher domain.CurrencyFetcher) *CurrencyUseCase {
	t.Helper()
	return NewCurrencyUseCase(
		repo,
		fetcher,
		config.CompareFinAPIURLType("https://min-api.cryptocompare.com"),
		config.CompareFinAPIKeyType("test-api-key"),
	)
}
