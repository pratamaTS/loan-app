package config

import (
	"example-loan.com/loan-app/app/controller"
	"example-loan.com/loan-app/app/repository"
	"example-loan.com/loan-app/app/service"
)

type Initialization struct {
	LoanRepository repository.LoanRepository
	LoanService    service.LoanService
	LoanController controller.LoanController
}

func Init(loanRepo repository.LoanRepository, loanService service.LoanService, loanController controller.LoanController) *Initialization {
	return &Initialization{
		LoanRepository: loanRepo,
		LoanService:    loanService,
		LoanController: loanController,
	}
}
