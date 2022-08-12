package packer

import (
	"encoding/json"
	"fmt"

	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransTargetReqToAPI(target *rao.CreateTargetReq) *model.API {
	if target.Request == nil {
		fmt.Sprintln(fmt.Errorf("target.request not found request"))
		return nil
	}

	headerByte, err := json.Marshal(target.Request.Header)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.header json marshal err %w", err))
	}

	queryByte, err := json.Marshal(target.Request.Query)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.query json marshal err %w", err))
	}

	bodyByte, err := json.Marshal(target.Request.Body)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.body json marshal err %w", err))
	}

	authByte, err := json.Marshal(target.Request.Auth)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.auth json marshal err %w", err))
	}

	return &model.API{
		TargetID:    target.TargetID,
		URL:         target.URL,
		Header:      string(headerByte),
		Query:       string(queryByte),
		Body:        string(bodyByte),
		Auth:        string(authByte),
		Description: target.Description,
	}
}
