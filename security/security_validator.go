package security

import (
	"errors"
	"fmt"
	"net/http"
)

func validateLoginReq(lReq *LoginReq) (*LoginReq, int, error) {
	if lReq.Email == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required parameter Email is empty [%s]", lReq.Email)
	}

	if lReq.ID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required parameter ID is empty [%s]", lReq.ID)
	}

	if lReq.Password == "" {
		return nil, http.StatusBadRequest, errors.New("required parameter Password is empty")
	}

	return lReq, http.StatusOK, nil
}
