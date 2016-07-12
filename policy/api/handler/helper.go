package handler

import (
	"encoding/json"

	"github.com/kazoup/gabs"
	proto "github.com/kazoup/platform/policy/srv/proto/policy"
	"github.com/micro/go-micro/errors"
)

// PolicyToStringJSON returns a JSON stringify policy.
func PolicyToStringJSON(policy *proto.ReadResponse) ([]byte, error) {
	var output map[string]interface{}
	var filterInterface interface{}

	marshaledSrvRsp, err := json.Marshal(policy)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Unmarshall srv response into a map[string]interface{}
	if err := json.Unmarshal(marshaledSrvRsp, &output); err != nil {
		return nil, errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Unmarshall JSON like string into interface
	// We do not know the structure, so github.com/Jeffail/gabs does not help for this task
	if err := json.Unmarshal([]byte(policy.Filter), &filterInterface); err != nil {
		return nil, errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Assign the interface
	output["filter"] = filterInterface

	b, err := json.Marshal(output)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	return b, nil
}

// PolicyToStringJSON returns a JSON stringify Array of policies.
func PoliciesToStringJSON(policies []*proto.ReadResponse) ([]byte, error) {
	result := gabs.New()
	result.Array("result")

	for _, v := range policies {
		var policy interface{}
		b, err := PolicyToStringJSON(v)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(b, &policy)
		result.ArrayAppend(policy, "result")
	}

	return result.Bytes(), nil
}
