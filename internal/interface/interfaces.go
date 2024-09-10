// Package usecase implements application business logic. Each logic group in own file.
package _interface

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Usecase -.
	Usecase interface {
	}

	// Repo -.
	Repo interface {
	}

	// WebAPI -.
	WebAPI interface {
	}
)
