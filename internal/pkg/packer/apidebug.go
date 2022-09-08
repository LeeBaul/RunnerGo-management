package packer

import (
	"encoding/json"

	"github.com/go-omnibus/proof"
	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransMaoAPIDebugToRaoAPIDebug(m *mao.APIDebug) *rao.APIDebug {
	rawNodes, err := m.Assertion.Values()
	if err != nil {
		proof.Errorf("api_debug.assertion get values err", proof.WithError(err))
	}

	var as []*rao.Assertion
	for _, node := range rawNodes {
		d, ok := node.DocumentOK()
		if !ok {
			proof.Errorf("api_debug.assertion DocumentOK err", proof.WithError(err))
		}

		var a rao.Assertion
		if err := bson.Unmarshal(d, &a); err != nil {
			proof.Errorf("api_debug.assertion bson unmarshal err", proof.WithError(err))
		}

		as = append(as, &a)
	}

	var regex interface{}
	if err := bson.Unmarshal(m.Regex, &regex); err != nil {
		proof.Errorf("api_debug.regex bson unmarshal err", proof.WithError(err))
	}

	r, err := json.Marshal(regex)
	if err != nil {
		proof.Errorf("api_debug.regex json marshal err", proof.WithError(err))
	}

	return &rao.APIDebug{
		ApiID:                 m.ApiID,
		APIName:               m.APIName,
		Assertion:             as,
		EventID:               m.EventID,
		Regex:                 string(r),
		RequestBody:           m.RequestBody,
		RequestCode:           m.RequestCode,
		RequestHeader:         m.RequestHeader,
		RequestTime:           m.RequestTime,
		ResponseBody:          m.ResponseBody,
		ResponseBytes:         m.ResponseBytes,
		ResponseHeader:        m.ResponseHeader,
		ResponseTime:          m.ResponseTime,
		ResponseLen:           m.ResponseLen,
		ResponseStatusMessage: m.ResponseStatusMessage,
		UUID:                  m.UUID,
	}
}
