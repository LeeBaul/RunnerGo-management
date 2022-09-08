package packer

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransMaoAPIDebugToRaoAPIDebug(m *mao.APIDebug) *rao.APIDebug {
	rawNodes, err := m.Assertion.Values()
	if err != nil {
		fmt.Println(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
	}

	var as []*rao.Assertion
	for _, node := range rawNodes {
		d, ok := node.DocumentOK()
		if !ok {
			fmt.Println(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
		}

		var a rao.Assertion
		if err := bson.Unmarshal(d, &a); err != nil {
			fmt.Println(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
		}

		as = append(as, &a)
	}

	var regex interface{}
	if err := bson.Unmarshal(m.Regex, &regex); err != nil {
		fmt.Println(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
	}

	r, err := json.Marshal(regex)
	if err != nil {
		fmt.Println(fmt.Errorf("apidebug.assertion json unmarshal err %w", err))
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
