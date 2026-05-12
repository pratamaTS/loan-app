package controller_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example-loan.com/loan-app/app/controller"
	"example-loan.com/loan-app/app/domain/dao"
	serviceMock "example-loan.com/loan-app/app/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestLoanController(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success request loan", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		reqBody := `{
			"user_id":"Bruce",
			"mrp":10000000,
			"dp":2000000,
			"vehicle_year":2022,
			"police_number":"B 1234 BYE",
			"machine_number":"ENG12345"
		}`

		mockService.
			On("RequestLoan", mock.Anything).
			Return(dao.Loan{}, nil)

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/request", ctrl.RequestLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/request", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid json request", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/request", ctrl.RequestLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/request", bytes.NewBufferString(`{invalid-json}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("service error request loan", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		reqBody := `{
			"user_id":"Bruce",
			"mrp":10000000,
			"dp":2000000,
			"vehicle_year":2022,
			"police_number":"B 1234 BYE",
			"machine_number":"ENG12345"
		}`

		mockService.
			On("RequestLoan", mock.Anything).
			Return(dao.Loan{}, errors.New("internal error"))

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/request", ctrl.RequestLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/request", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)

		mockService.AssertExpectations(t)
	})
}

func TestApproveLoanController(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("success approve loan", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		reqBody := `{
			"user_id":"Bruce",
			"police_number":"B 1234 BYE"
		}`

		mockService.
			On("ApproveLoan", mock.Anything).
			Return(nil)

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/approve", ctrl.ApproveLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/approve", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid json approve loan", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/approve", ctrl.ApproveLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/approve", bytes.NewBufferString(`{invalid-json}`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})

	t.Run("service error approve loan", func(t *testing.T) {
		mockService := new(serviceMock.MockLoanService)

		reqBody := `{
			"user_id":"Bruce",
			"police_number":"B 1234 BYE"
		}`

		mockService.
			On("ApproveLoan", mock.Anything).
			Return(errors.New("internal error"))

		ctrl := controller.LoanControllerInit(mockService)

		router := gin.Default()
		router.POST("/loan/approve", ctrl.ApproveLoan)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/loan/approve", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)

		mockService.AssertExpectations(t)
	})
}
