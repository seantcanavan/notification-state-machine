package audit

import (
	"fmt"
	"net/http"
)

func validateCreateReq(cReq *CreateReq) (*CreateReq, int, error) {
	if cReq.JobID == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("missing requied parameter JobID [%s]", cReq.JobID)
	}

	if cReq.NextStatus == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("missing required parameter NextStatus [%+v]", cReq.NextStatus)
	}

	if !cReq.NextStatus.Valid() {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid required parameter NextStatus [%+v]", cReq.NextStatus)
	}

	if cReq.Operation == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("missing required parameter Operation [%+v]", cReq.Operation)
	}

	if !cReq.Operation.Valid() {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid required parameter Operation [%+v]", cReq.Operation)
	}

	if cReq.PreviousStatus == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("missing required parameter PreviousStatus [%+v]", cReq.PreviousStatus)
	}

	if !cReq.PreviousStatus.Valid() {
		return nil, http.StatusBadRequest, fmt.Errorf("invalid required parameter PreviousStatus [%+v]", cReq.PreviousStatus)
	}

	return cReq, http.StatusOK, nil
}
