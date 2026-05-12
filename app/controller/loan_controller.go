package controller

import (
	"example-loan.com/loan-app/app/domain/dto"
	"example-loan.com/loan-app/app/service"
	"github.com/gin-gonic/gin"
)

type LoanController interface {
	RequestLoan(c *gin.Context)
	ApproveLoan(c *gin.Context)
}

type LoanControllerImpl struct {
	loanService service.LoanService
}

func LoanControllerInit(loanService service.LoanService) *LoanControllerImpl {
	return &LoanControllerImpl{
		loanService: loanService,
	}
}

func (ctrl *LoanControllerImpl) RequestLoan(c *gin.Context) {
	var req dto.RequestLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	loan, err := ctrl.loanService.RequestLoan(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": loan})
}

func (ctrl *LoanControllerImpl) ApproveLoan(c *gin.Context) {
	var req dto.ApproveLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := ctrl.loanService.ApproveLoan(&req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Loan approved successfully"})
}
