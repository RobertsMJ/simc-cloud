package health

import "context"

type Checker interface {
	Name() string
	Check(context.Context) error
}

type Handler struct {
	checkers []Checker
}

func NewHandler(checkers []Checker) *Handler {

	return &Handler{
		checkers: checkers,
	}
}

func (h *Handler) Handle() error {
	return nil
}
