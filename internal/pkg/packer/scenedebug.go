package packer

import (
	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransMaoSceneDebugsToRaoSceneDebugs(ms []*mao.SceneDebug) []*rao.SceneDebug {
	ret := make([]*rao.SceneDebug, 0)

	for _, m := range ms {

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

		ret = append(ret, &rao.SceneDebug{
			ApiID:          m.ApiID,
			APIName:        m.APIName,
			Assertion:      as,
			EventID:        m.EventID,
			NextList:       m.NextList,
			Regex:          rs,
			RequestBody:    m.RequestBody,
			RequestCode:    m.RequestCode,
			RequestHeader:  m.RequestHeader,
			ResponseBody:   m.ResponseBody,
			ResponseBytes:  m.ResponseBytes,
			ResponseHeader: m.ResponseHeader,
			Status:         m.Status,
			UUID:           m.UUID,
		})
	}

	return ret
}
