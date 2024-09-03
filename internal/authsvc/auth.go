package authsvc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	authPB "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
)

type AuthSvc struct {
	authServerUrl string
	httpClient    *http.Client
}

func NewAuthSvc(authServerUrl string) *AuthSvc {
	httpClient := http.Client{Timeout: time.Duration(1) * time.Second}

	return &AuthSvc{
		authServerUrl: authServerUrl,
		httpClient:    &httpClient,
	}
}

type AuthRequest struct {
	Token string `json:"token"`
}

var NoAuthPaths [2]string

func init() {
	NoAuthPaths = [...]string{"favicon.ico", "index.html"}
}

func isNoAuthPath(path string) bool {

	for _, noAuthPath := range NoAuthPaths {
		if strings.TrimSpace(path) == noAuthPath {
			return true
		}
	}

	return false
}

func (aSvc *AuthSvc) validateToken(token string) error {

	authRequest := AuthRequest{
		Token: token,
	}

	reqBodyBuffer := new(bytes.Buffer)
	if err := json.NewEncoder(reqBodyBuffer).Encode(authRequest); err != nil {
		log.Println("failed to construct json request ", err)
		return err
	}

	req, err := http.NewRequest("POST", aSvc.authServerUrl, reqBodyBuffer)

	if err != nil {
		log.Println("failed to construct request ", err)
		return err
	}

	if _, err := aSvc.httpClient.Do(req); err != nil {
		log.Println("failed to authorized ", err)
		return err
	}

	return nil

}

func (aSvc *AuthSvc) Check(ctx context.Context, req *authPB.CheckRequest) (*authPB.CheckResponse, error) {

	requestPath := req.Attributes.Request.Http.Path[1:]

	headers := req.GetAttributes().GetRequest().GetHttp().GetHeaders()

	if isNoAuthPath(requestPath) {
		return &authPB.CheckResponse{}, nil
	}

	log.Println("will authorize ", requestPath)
	// if its a no auth uri, let it go

	// check bearer token

	bearerToken := headers["authorization"]

	if len(bearerToken) == 0 {
		log.Println("failed to locate bearer token on ", requestPath)
		// if its a no auth uri, let it go
		return nil, fmt.Errorf("failed to locate bearer token ")
	}

	jwtToken := bearerToken[7:]

	if err := aSvc.validateToken(jwtToken); err != nil {
		return nil, err
	}

	log.Println("will authorize ", requestPath)
	return &authPB.CheckResponse{}, nil
}
