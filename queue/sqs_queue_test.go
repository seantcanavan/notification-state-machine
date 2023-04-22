package queue

import (
	"fmt"
	"github.com/jgroeneveld/trial/assert"
	"github.com/seantcanavan/lambda_jwt_router/lambda_util"
	"testing"
)

func TestGenerateDeduplicationID(t *testing.T) {
	ID := lambda_util.GenerateRandomString(24)
	Resource := CreateReq
	message := &Message{
		ID:       ID,
		Resource: Resource,
		Type:     CreateReq,
	}
	dedupID := generateDeduplicationID(message)
	fmt.Println(fmt.Sprintf(dedupID))
	fmt.Println(fmt.Sprintf("len [%d]", len(dedupID)))
	assert.True(t, len(dedupID) <= 128)
}
