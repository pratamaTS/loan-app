package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"example-loan.com/loan-app/app/domain/dao"
	"example-loan.com/loan-app/app/domain/dto"
	"example-loan.com/loan-app/app/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoanRepository interface {
	RequestLoan(req *dto.RequestLoanRequest) (dao.Loan, error)
	ApproveLoan(req *dto.ApproveLoanRequest) error
}

type LoanRepositoryImpl struct {
	loanCol *mongo.Collection
}

func LoanRepositoryInit(mongoClient *mongo.Client) *LoanRepositoryImpl {
	dbName := "loan_service_db"
	loanCol := mongoClient.Database(dbName).Collection("loans")

	return &LoanRepositoryImpl{loanCol: loanCol}
}

func (r *LoanRepositoryImpl) RequestLoan(req *dto.RequestLoanRequest) (dao.Loan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	if strings.TrimSpace(req.UserID) == "" {
		return dao.Loan{}, errors.New("user_id required")
	}

	if req.MRP <= 0 {
		return dao.Loan{}, errors.New("mrp required")
	}

	if req.DP <= 0 {
		return dao.Loan{}, errors.New("dp required")
	}

	if req.VehicleYear <= 0 {
		return dao.Loan{}, errors.New("vehicle_year required")
	}

	if strings.TrimSpace(req.PoliceNumber) == "" {
		return dao.Loan{}, errors.New("police_number required")
	}

	if strings.TrimSpace(req.MachineNumber) == "" {
		return dao.Loan{}, errors.New("machine_number required")
	}

	if req.DP > req.MRP {
		return dao.Loan{}, errors.New("dp cannot exceed mrp")
	}

	session, err := r.loanCol.Database().Client().StartSession()
	if err != nil {
		return dao.Loan{}, err
	}

	defer session.EndSession(ctx)
	var loan dao.Loan

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		count, err := r.loanCol.CountDocuments(sessCtx, bson.M{"police_number": req.PoliceNumber})
		if err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, errors.New("police number already exists")
		}

		now := time.Now()

		loan = dao.Loan{
			ID:            helpers.GenerateID(),
			UserID:        req.UserID,
			MRP:           req.MRP,
			DP:            req.DP,
			VehicleYear:   req.VehicleYear,
			PoliceNumber:  req.PoliceNumber,
			MachineNumber: req.MachineNumber,
			Status:        helpers.StatusSubmitted,
			RequestDate:   now.Format(time.RFC3339),
			CreatedAt:     now.Format(time.RFC3339),
			ApprovedDate:  "",
		}

		_, err = r.loanCol.InsertOne(sessCtx, loan)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, errors.New("police number already exists")
			}

			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return dao.Loan{}, err
	}

	return loan, nil
}

func (r *LoanRepositoryImpl) ApproveLoan(req *dto.ApproveLoanRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("user_id required")
	}

	if strings.TrimSpace(req.PoliceNumber) == "" {
		return errors.New("police_number required")
	}

	session, err := r.loanCol.Database().Client().StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		var loan dao.Loan

		err := r.loanCol.FindOne(sessCtx, bson.M{"user_id": req.UserID, "police_number": req.PoliceNumber}).Decode(&loan)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, errors.New("loan_not_found")
			}

			return nil, err
		}

		if loan.Status == helpers.StatusApproved {
			return nil, errors.New("loan already approved")
		}

		now := time.Now()

		result, err := r.loanCol.UpdateOne(sessCtx,
			bson.M{
				"id":     loan.ID,
				"status": helpers.StatusSubmitted,
			},
			bson.M{
				"$set": bson.M{
					"status":        helpers.StatusApproved,
					"updated_at":    now,
					"approved_date": now.Format(time.RFC3339)},
			},
		)
		if err != nil {
			return nil, err
		}

		if result.ModifiedCount == 0 {
			return nil, errors.New("loan already approved")
		}

		return nil, nil
	})

	return err
}
