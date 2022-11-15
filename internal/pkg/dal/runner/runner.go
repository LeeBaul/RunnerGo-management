package runner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-omnibus/proof"
	"github.com/go-resty/resty/v2"

	"kp-management/internal/pkg/conf"
	"kp-management/internal/pkg/dal/rao"
)

type RunAPIResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type StopRunnerReq struct {
	TeamID    int64    `json:"team_id"`
	PlanID    int64    `json:"plan_id"`
	ReportIds []string `json:"report_ids"`
}

func RunAPI(ctx context.Context, body *rao.APIDetail) (string, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	proof.Infof("body %s", bodyByte)

	var ret RunAPIResp
	_, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyByte).
		SetResult(&ret).
		Post(conf.Conf.Clients.Runner.RunAPI)

	if err != nil {
		return "", err
	}

	if ret.Code != 200 {
		return "", fmt.Errorf("ret code not 200")
	}

	return ret.Data, nil
}

func RunScene(ctx context.Context, body *rao.SceneFlow) (string, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	proof.Infof("body %s", bodyByte)

	var ret RunAPIResp
	_, err = resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyByte).
		SetResult(&ret).
		Post(conf.Conf.Clients.Runner.RunScene)

	if err != nil {
		return "", err
	}

	if ret.Code != 200 {
		return "", fmt.Errorf("ret code not 200")
	}

	return ret.Data, nil
}

func StopScene(ctx context.Context, req *rao.StopSceneReq) error {
	var ret RunAPIResp
	_, err := resty.New().R().
		SetBody(req).
		SetResult(&ret).
		Post(conf.Conf.Clients.Runner.StopScene)

	if err != nil {
		return err
	}

	if ret.Code != 200 {
		return fmt.Errorf("ret code not 200")
	}

	return nil
}
