package mocks

import (
	"example-loan.com/loan-app/app/domain/dao"
	"example-loan.com/loan-app/app/domain/dto"

	"github.com/stretchr/testify/mock"
)

type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) RequestLoan(req *dto.RequestLoanRequest) (dao.Loan, error) {
	args := m.Called(req)

	var loan dao.Loan

	if args.Get(0) != nil {
		loan = args.Get(0).(dao.Loan)
	}

	return loan, args.Error(1)
}

func (m *MockLoanRepository) ApproveLoan(req *dto.ApproveLoanRequest) error {
	args := m.Called(req)

	return args.Error(0)
}
