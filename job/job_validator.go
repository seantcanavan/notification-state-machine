package job

import (
	"fmt"
	"net/http"
)

func validateCreateReq(cReq *CreateReq) (*CreateReq, int, error) {
	if cReq.From == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field From cannot be empty [%s]", cReq.From)
	}

	if cReq.Template == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Template cannot be empty [%s]", cReq.Template)
	}

	if cReq.To == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field To cannot be empty [%s]", cReq.To)
	}

	if cReq.Type == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Type cannot be empty [%s]", cReq.Type)
	}

	if !cReq.Type.Valid() {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Type is invalid [%+v]", cReq.Type)
	}

	if cReq.Variables == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Variables cannot be nil [%+v]", cReq.Variables)
	}

	return cReq, http.StatusOK, nil
}

func validateUpdateReq(uReq *UpdateReq) (*UpdateReq, int, error) {
	if uReq.ID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field ID cannot be empty [%s]", uReq.ID)
	}

	if uReq.Status == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Status cannot be empty [%+v]", uReq.Status)
	}

	if !uReq.Status.Valid() {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Status cannot be empty [%+v]", uReq.Status)
	}

	if uReq.Variables == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("required field Variables cannot be nil [%+v]", uReq.Variables)
	}

	if uReq.Email == nil && uReq.SMS == nil && uReq.Snail == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("optional fields Email [%+v] SMS [%+v] and Snail [%+v] cannot all be empty", uReq.Email, uReq.SMS, uReq.Snail)
	}

	return uReq, http.StatusOK, nil
}
