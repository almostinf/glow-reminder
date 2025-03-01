package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	logger_mocks "github.com/almostinf/glow-reminder/pkg/logger/mocks"
	"github.com/almostinf/glow-reminder/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	invalidURL = "test"
	validURL   = `postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable`

	connTimeout  = time.Microsecond
	connAttempts = 30
	maxPoolSize  = 1
)

func postgresHelper(t *testing.T) *logger_mocks.MockLogger {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	log := logger_mocks.NewMockLogger(mockCtrl)

	return log
}

func TestNewPostgres(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	type args struct {
		ctx context.Context
		cfg postgres.Config
	}

	parseConfigErr := errors.New("cannot parse `test`: failed to parse as DSN (invalid dsn)")

	testcases := []struct {
		name        string
		args        args
		mock        func(*logger_mocks.MockLogger)
		expectedErr error
	}{
		{
			name: "failed to parse pool config url",
			args: args{
				ctx: ctx,
				cfg: postgres.Config{
					URL:          invalidURL,
					ConnAttempts: connAttempts,
					ConnTimeout:  connTimeout,
					MaxPoolSize:  maxPoolSize,
				},
			},
			mock:        func(*logger_mocks.MockLogger) {},
			expectedErr: fmt.Errorf("failed to connect to postgres pool: %w", parseConfigErr),
		},
		{
			name: "successfully connected to postgres pool",
			args: args{
				ctx: ctx,
				cfg: postgres.Config{
					URL:          validURL,
					ConnAttempts: connAttempts,
					ConnTimeout:  connTimeout,
					MaxPoolSize:  maxPoolSize,
				},
			},
			mock:        func(*logger_mocks.MockLogger) {},
			expectedErr: nil,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()

			mockLog := postgresHelper(t)
			testcase.mock(mockLog)

			_, err := postgres.New(
				testcase.args.ctx,
				testcase.args.cfg,
				mockLog,
			)

			if testcase.expectedErr != nil {
				assert.Errorf(t, err, testcase.expectedErr.Error())
			} else {
				assert.Equal(t, testcase.expectedErr, err)
			}
		})
	}
}
