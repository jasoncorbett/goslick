package slickqa

import (
	context "golang.org/x/net/context"
)


type SlickAuthService struct {}

func (s *SlickAuthService) IsAuthorized(context.Context, *IsAuthorizedRequest) (*IsAuthorizedResponse, error) {
	return &IsAuthorizedResponse{
		Allowed: true,
	}, nil
}
