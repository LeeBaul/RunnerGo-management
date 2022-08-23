package packer

import (
	"fmt"

	"github.com/bytedance/sonic"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/rao"
)

func TransTargetReqToAPI(target *rao.CreateTargetReq) *mao.API {
	if target.Request == nil {
		fmt.Sprintln(fmt.Errorf("target.request not found request"))
		return nil
	}

	header, err := sonic.MarshalString(target.Request.Header)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.header json marshal err %w", err))
	}

	query, err := sonic.MarshalString(target.Request.Query)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.query json marshal err %w", err))
	}

	body, err := sonic.MarshalString(target.Request.Body)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.body json marshal err %w", err))
	}

	auth, err := sonic.MarshalString(target.Request.Auth)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.auth json marshal err %w", err))
	}

	return &mao.API{
		TargetID:    target.TargetID,
		URL:         target.URL,
		Header:      header,
		Query:       query,
		Body:        body,
		Auth:        auth,
		Description: target.Description,
	}
}
