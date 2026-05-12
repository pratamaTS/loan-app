package repository_test

import (
	"errors"
	"sync"
	"testing"

	"example-loan.com/loan-app/app/domain/dao"
	dto "example-loan.com/loan-app/app/domain/dto"
	repositoryMock "example-loan.com/loan-app/app/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanRepository(t *testing.T) {

	t.Run("success request loan", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.RequestLoanRequest{
			UserID:        "Bruce",
			MRP:           10000000,
			DP:            2000000,
			VehicleYear:   2022,
			PoliceNumber:  "B 1234 BYE",
			MachineNumber: "ENG12345",
		}

		expected := dao.Loan{}

		mockRepo.
			On("RequestLoan", req).
			Return(expected, nil)

		loan, err := mockRepo.RequestLoan(req)

		assert.NoError(t, err)
		assert.Equal(t, expected, loan)

		mockRepo.AssertExpectations(t)
	})

	t.Run("success approve loan", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.ApproveLoanRequest{
			UserID:       "Bruce",
			PoliceNumber: "B 1234 BYE",
		}

		mockRepo.
			On("ApproveLoan", req).
			Return(nil)

		err := mockRepo.ApproveLoan(req)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("request loan error", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.RequestLoanRequest{
			UserID:        "Bruce",
			MRP:           10000000,
			DP:            2000000,
			VehicleYear:   2022,
			PoliceNumber:  "B 1234 BYE",
			MachineNumber: "ENG12345",
		}

		mockRepo.
			On("RequestLoan", req).
			Return(dao.Loan{}, errors.New("database error"))

		loan, err := mockRepo.RequestLoan(req)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, dao.Loan{}, loan)

		mockRepo.AssertExpectations(t)
	})

	t.Run("race condition approve loan concurrent", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.ApproveLoanRequest{
			UserID:       "Bruce",
			PoliceNumber: "B 1234 BYE",
		}

		mockRepo.
			On("ApproveLoan", mock.Anything).
			Return(nil)

		wg := sync.WaitGroup{}
		concurrency := 20
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				_ = mockRepo.ApproveLoan(req)
			}()
		}

		wg.Wait()

		mockRepo.AssertNumberOfCalls(t, "ApproveLoan", concurrency)
	})

	t.Run("race condition request loan concurrent", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.RequestLoanRequest{
			UserID:        "Bruce",
			MRP:           10000000,
			DP:            2000000,
			VehicleYear:   2022,
			PoliceNumber:  "B 1234 BYE",
			MachineNumber: "ENG12345",
		}

		mockRepo.
			On("RequestLoan", mock.Anything).
			Return(dao.Loan{}, nil)

		wg := sync.WaitGroup{}
		concurrency := 20
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				_, _ = mockRepo.RequestLoan(req)
			}()
		}

		wg.Wait()

		mockRepo.AssertNumberOfCalls(t, "RequestLoan", concurrency)
	})
}
