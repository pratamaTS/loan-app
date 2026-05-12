//go:build wireinject
// +build wireinject

package config

import (
	"example-loan.com/loan-app/app/controller"
	"example-loan.com/loan-app/app/repository"
	"example-loan.com/loan-app/app/service"
	"github.com/google/wire"
)

var db = wire.NewSet(ConnectToMongoDB)

var repositorySet = wire.NewSet(
	repository.LoanRepositoryInit, wire.Bind(new(repository.LoanRepository), new(*repository.LoanRepositoryImpl)),
)

var serviceSet = wire.NewSet(
	service.LoanServiceInit, wire.Bind(new(service.LoanService), new(*service.LoanServiceImpl)),
)

var controllerSet = wire.NewSet(
	controller.LoanControllerInit, wire.Bind(new(controller.LoanController), new(*controller.LoanControllerImpl)),
)

func InitializeApp() (*Initialization, error) {
	wire.Build(db, repositorySet, serviceSet, controllerSet, Init)
	return &Initialization{}, nil
}
