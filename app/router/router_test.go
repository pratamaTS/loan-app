package router_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"example-loan.com/loan-app/app/config"
	"example-loan.com/loan-app/app/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockController struct {
	RequestLoanFunc func(c *gin.Context)
	ApproveLoanFunc func(c *gin.Context)
}

func (m *MockController) RequestLoan(c *gin.Context) {
	m.RequestLoanFunc(c)
}

func (m *MockController) ApproveLoan(c *gin.Context) {
	m.ApproveLoanFunc(c)
}

func TestRouter(t *testing.T) {

	gin.SetMode(gin.TestMode)

	mockCtrl := &MockController{}

	initApp := &config.Initialization{
		LoanController: mockCtrl,
	}

	r := router.InitRouter(initApp)

	t.Run("health check", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/health", nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("request loan route", func(t *testing.T) {

		mockCtrl.RequestLoanFunc = func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		}

		body := `{
			"user_id":"Bruce",
			"mrp":10000000,
			"dp":2000000,
			"vehicle_year":2022,
			"police_number":"B 1234 BYE",
			"machine_number":"ENG12345"
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/loan/request", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("approve loan route", func(t *testing.T) {

		mockCtrl.ApproveLoanFunc = func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "approved"})
		}

		body := `{
			"user_id":"Bruce",
			"police_number":"B 1234 BYE"
		}`

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/loan/approve", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}
