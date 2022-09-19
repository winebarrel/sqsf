package sqsf

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func marshalMessage(message types.Message, decode bool) ([]byte, error) {
	if !decode {
		return json.MarshalIndent(message, "", "    ")
	}

	body := aws.ToString(message.Body)
	var m map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(body)))
	decoder.UseNumber()
	err := decoder.Decode(&m)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, body)
	}

	return json.MarshalIndent(m, "", "    ")
}
