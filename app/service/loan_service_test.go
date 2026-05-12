package service_test

import (
	"errors"
	"sync"
	"testing"

	"example-loan.com/loan-app/app/domain/dao"
	"example-loan.com/loan-app/app/domain/dto"
	"example-loan.com/loan-app/app/helpers"
	repositoryMock "example-loan.com/loan-app/app/repository/mocks"
	"example-loan.com/loan-app/app/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestLoan(t *testing.T) {

	validReq := &dto.RequestLoanRequest{
		UserID:        "Bruce",
		MRP:           10000000,
		DP:            2000000,
		VehicleYear:   2022,
		PoliceNumber:  "B 1234 BYE",
		MachineNumber: "ENG12345",
	}

	t.Run("success request loan", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		expectedLoan := dao.Loan{}

		mockRepo.
			On("RequestLoan", validReq).
			Return(expectedLoan, nil)

		svc := service.LoanServiceInit(mockRepo)

		loan, err := svc.RequestLoan(validReq)

		assert.NoError(t, err)
		assert.Equal(t, expectedLoan, loan)

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		mockRepo.
			On("RequestLoan", validReq).
			Return(dao.Loan{}, errors.New("db error"))

		svc := service.LoanServiceInit(mockRepo)

		loan, err := svc.RequestLoan(validReq)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		assert.Equal(t, dao.Loan{}, loan)

		mockRepo.AssertExpectations(t)
	})

	t.Run("negative - nil request", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		svc := service.LoanServiceInit(mockRepo)

		loan, err := svc.RequestLoan(nil)

		assert.Error(t, err)
		assert.Equal(t, "request payload required", err.Error())
		assert.Equal(t, dao.Loan{}, loan)
	})

	t.Run("negative - invalid dp greater than mrp", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		req := &dto.RequestLoanRequest{
			UserID:        "Bruce",
			MRP:           1000000,
			DP:            5000000,
			VehicleYear:   2022,
			PoliceNumber:  "B 1234 BYE",
			MachineNumber: "ENG12345",
		}

		svc := service.LoanServiceInit(mockRepo)

		loan, err := svc.RequestLoan(req)

		assert.Error(t, err)
		assert.Equal(t, helpers.ErrInvalidDP, err.Error())
		assert.Equal(t, dao.Loan{}, loan)
	})

	t.Run("race condition request loan concurrent", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		mockRepo.
			On("RequestLoan", mock.Anything).
			Return(dao.Loan{}, nil)

		svc := service.LoanServiceInit(mockRepo)

		wg := sync.WaitGroup{}
		concurrency := 30
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				_, _ = svc.RequestLoan(validReq)
			}()
		}

		wg.Wait()

		mockRepo.AssertNumberOfCalls(t, "RequestLoan", concurrency)
	})
}

func TestApproveLoan(t *testing.T) {

	validReq := &dto.ApproveLoanRequest{
		UserID:       "Bruce",
		PoliceNumber: "B 1234 BYE",
	}

	t.Run("success approve loan", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		mockRepo.
			On("ApproveLoan", validReq).
			Return(nil)

		svc := service.LoanServiceInit(mockRepo)

		err := svc.ApproveLoan(validReq)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		mockRepo.
			On("ApproveLoan", validReq).
			Return(errors.New("db error"))

		svc := service.LoanServiceInit(mockRepo)

		err := svc.ApproveLoan(validReq)

		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})

	t.Run("negative - nil request", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		svc := service.LoanServiceInit(mockRepo)

		err := svc.ApproveLoan(nil)

		assert.Error(t, err)
		assert.Equal(t, "request payload required", err.Error())
	})

	t.Run("race condition approve loan concurrent", func(t *testing.T) {
		mockRepo := new(repositoryMock.MockLoanRepository)

		mockRepo.
			On("ApproveLoan", mock.Anything).
			Return(nil)

		svc := service.LoanServiceInit(mockRepo)

		wg := sync.WaitGroup{}
		concurrency := 30
		wg.Add(concurrency)

		for i := 0; i < concurrency; i++ {
			go func() {
				defer wg.Done()
				_ = svc.ApproveLoan(validReq)
			}()
		}

		wg.Wait()

		mockRepo.AssertNumberOfCalls(t, "ApproveLoan", concurrency)
	})
}
