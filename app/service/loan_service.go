package service

import (
	"errors"
	"strings"
	"time"

	"example-loan.com/loan-app/app/domain/dao"
	"example-loan.com/loan-app/app/domain/dto"
	"example-loan.com/loan-app/app/helpers"
	"example-loan.com/loan-app/app/repository"
)

type LoanService interface {
	RequestLoan(req *dto.RequestLoanRequest) (dao.Loan, error)
	ApproveLoan(req *dto.ApproveLoanRequest) error
}

type LoanServiceImpl struct {
	loanRepo repository.LoanRepository
}

func LoanServiceInit(loanRepo repository.LoanRepository) *LoanServiceImpl {
	return &LoanServiceImpl{
		loanRepo: loanRepo,
	}
}

func (s *LoanServiceImpl) RequestLoan(req *dto.RequestLoanRequest) (dao.Loan, error) {
	if req == nil {
		return dao.Loan{}, errors.New("request payload required")
	}

	if req.UserID == "" {
		return dao.Loan{}, errors.New("user_id required")
	}

	if req.PoliceNumber == "" {
		return dao.Loan{}, errors.New("police_number required")
	}

	if req.MachineNumber == "" {
		return dao.Loan{}, errors.New("machine_number required")
	}

	if req.MRP <= 0 {
		return dao.Loan{}, errors.New("mrp must be greater than 0")
	}

	if req.DP <= 0 {
		return dao.Loan{}, errors.New("dp must be greater than 0")
	}

	req.UserID = strings.TrimSpace(req.UserID)
	req.PoliceNumber = helpers.NormalizePoliceNumber(req.PoliceNumber)
	if !helpers.PoliceNumberRegex.MatchString(req.PoliceNumber) {
		return dao.Loan{}, errors.New("invalid police number format")
	}

	req.MachineNumber = helpers.NormalizeMachineNumber(req.MachineNumber)
	if !helpers.MachineNumberRegex.MatchString(req.MachineNumber) {
		return dao.Loan{}, errors.New("invalid machine number format")
	}

	if req.DP > req.MRP {
		return dao.Loan{}, errors.New("dp cannot be greater than mrp")
	}

	minDP := req.MRP * 0.10
	if req.DP < minDP {
		return dao.Loan{}, errors.New("minimum dp is 10% of mrp")
	}

	if req.MRP > 10000000000 {
		return dao.Loan{}, errors.New("mrp too large")
	}

	currentYear := time.Now().Year()
	if req.VehicleYear < 1920 {
		return dao.Loan{}, errors.New("vehicle year invalid")
	}

	if req.VehicleYear > currentYear+1 {
		return dao.Loan{}, errors.New("vehicle year cannot exceed current year")
	}

	loan, err := s.loanRepo.RequestLoan(req)
	if err != nil {
		switch err.Error() {
		case "police number already exists":
			return dao.Loan{}, errors.New("duplicate vehicle application")

		case "duplicate key error":
			return dao.Loan{}, errors.New("duplicate vehicle application")

		default:
			return dao.Loan{}, err
		}
	}

	return loan, nil
}

func (s *LoanServiceImpl) ApproveLoan(req *dto.ApproveLoanRequest) error {
	if req == nil {
		return errors.New(
			"request payload required",
		)
	}

	if req.UserID == "" {
		return errors.New("user_id required")
	}

	if req.PoliceNumber == "" {
		return errors.New("police_number required")
	}

	req.UserID = strings.TrimSpace(req.UserID)
	req.PoliceNumber = helpers.NormalizePoliceNumber(req.PoliceNumber)
	if !helpers.PoliceNumberRegex.MatchString(req.PoliceNumber) {
		return errors.New("invalid police number format")
	}

	err := s.loanRepo.ApproveLoan(req)
	if err != nil {
		switch err.Error() {

		case "loan_not_found":
			return errors.New("loan not found")

		case "loan already approved":
			return errors.New("loan already approved")

		default:
			return err
		}
	}

	return nil
}
