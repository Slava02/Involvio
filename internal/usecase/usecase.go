package usecase

import "github.com/Slava02/Involvio/internal/interface"

type Usecase struct {
}

type Deps struct {
	Repo _interface.Repo
}

func New(deps *Deps) *Usecase {
	return &Usecase{}
}
