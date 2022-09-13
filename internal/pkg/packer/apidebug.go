package packer

import (
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransMaoAPIDebugToRaoAPIDebug(m *mao.APIDebug) *rao.APIDebug {

	var as []*rao.DebugAssertion
	for _, a := range m.Assertion {
		as = append(as, &rao.DebugAssertion{
			Code:      a.Code,
			IsSucceed: a.IsSucceed,
			Msg:       a.Msg,
		})
	}

	var rs []*rao.DebugRegex
	for _, r := range m.Regex {
		rs = append(rs, &rao.DebugRegex{
			Code: r.Code,
		})
	}

	return &rao.APIDebug{
		ApiID:                 m.ApiID,
		APIName:               m.APIName,
		Assertion:             as,
		EventID:               m.EventID,
		Regex:                 rs,
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
