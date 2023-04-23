package security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/seantcanavan/lambda_jwt_router/lambda_jwt"
	"net/http"
	"os"
	"time"
)

type LoginReq struct {
	Email    string `json:"email,omitempty"`
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRes struct {
	JWT string `json:"jwt,omitempty"`
}

func Login(lReq *LoginReq) (*LoginRes, int, error) {
	lReq, httpStatus, err := validateLoginReq(lReq)
	if err != nil {
		return nil, httpStatus, err
	}

	extendedClaims := lambda_jwt.ExtendStandard(jwt.StandardClaims{
		Audience:  lReq.ID,
		ExpiresAt: time.Now().Add(time.Hour * 1440).Unix(), // 60 days = 1440 hours
		IssuedAt:  time.Now().Unix(),
		Issuer:    os.Getenv("APP_NAME") + "-" + os.Getenv("STAGE"),
		NotBefore: time.Now().Add(time.Hour * -1).Unix(),
		Subject:   lReq.ID,
	})

	fmt.Println(fmt.Sprintf("about to sign [%+v]", extendedClaims))
	userJWT, err := lambda_jwt.Sign(extendedClaims)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &LoginRes{
		JWT: userJWT,
	}, http.StatusOK, nil
}
