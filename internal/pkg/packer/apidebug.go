package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransMaoAPIDebugToRaoAPIDebug(m *mao.APIDebug) *rao.APIDebug {
	rawNodes, err := m.Assertion.Values()
	if err != nil {
		fmt.Sprintln(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
	}

	var as []*rao.Assertion
	for _, node := range rawNodes {
		d, ok := node.DocumentOK()
		if !ok {
			fmt.Sprintln(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
		}

		var a rao.Assertion
		if err := bson.Unmarshal(d, &a); err != nil {
			fmt.Sprintln(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
		}

		as = append(as, &a)
	}

	return &rao.APIDebug{
		ApiID:                 m.ApiID,
		APIName:               m.APIName,
		Assertion:             as,
		EventID:               m.EventID,
		Regex:                 m.Regex.String(),
		RequestBody:           m.RequestBody.String(),
		RequestCode:           m.RequestCode,
		RequestHeader:         m.RequestHeader,
		RequestTime:           m.RequestTime,
		ResponseBody:          m.ResponseBody.String(),
		ResponseBytes:         m.ResponseBytes,
		ResponseHeader:        m.ResponseHeader,
		ResponseTime:          m.ResponseTime,
		ResponseLen:           m.ResponseLen,
		ResponseStatusMessage: m.ResponseStatusMessage,
		UUID:                  m.UUID,
	}
}
