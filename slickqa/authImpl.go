package slickqa

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"strings"
	"github.com/jasoncorbett/goslick/jwtauth"
)


type SlickAuthService struct {}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (s *SlickAuthService) IsAuthorized(ctx context.Context, req *IsAuthorizedRequest) (*IsAuthorizedResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	notAuthorized := &IsAuthorizedResponse{
		Allowed: false,
	}

	if !ok {
		return notAuthorized, nil
	}
	auth, ok := md["authorization"]
	if !ok {
		return notAuthorized, nil
	}
	if len(auth) < 1 {
		return notAuthorized, nil
	}

	if !strings.HasPrefix(auth[0], "Bearer ") {
		return notAuthorized, nil
	}

	token := auth[0][7:]
	permissions, err := jwtauth.PermissionsFromJWT(token)
	if err != nil {
		return notAuthorized, nil
	}

	if stringInSlice(req.Permission, permissions) {
		return &IsAuthorizedResponse{
			Allowed: true,
		}, nil
	} else {
		return notAuthorized, nil
	}

}
