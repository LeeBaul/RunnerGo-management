package packer

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransTargetReqToAPI(target *rao.CreateTargetReq) *mao.API {
	if target.Request == nil {
		fmt.Sprintln(fmt.Errorf("target.request not found request"))
		return nil
	}

	header, err := bson.Marshal(target.Request.Header)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.header json marshal err %w", err))
	}

	query, err := bson.Marshal(target.Request.Query)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.query json marshal err %w", err))
	}

	body, err := bson.Marshal(target.Request.Body)
	if err != nil {
		fmt.Sprintln(fmt.Errorf("target.request.body json marshal err %w", err))
	}

	auth, err := bson.Marshal(target.Request.Auth)
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

func TransTargetToAPIDetail(targets []*model.Target, apis []*mao.API) []*rao.APIDetail {
	ret := make([]*rao.APIDetail, 0, len(targets))

	for _, target := range targets {
		for _, api := range apis {
			if api.TargetID == target.ID {

				var auth rao.Auth
				bson.Unmarshal(api.Auth, &auth)

				var body rao.Body
				bson.Unmarshal(api.Body, &body)

				var header rao.Header
				bson.Unmarshal(api.Header, &header)

				var query rao.Query
				bson.Unmarshal(api.Query, &query)

				ret = append(ret, &rao.APIDetail{
					TargetID:   target.ID,
					ParentID:   target.ParentID,
					TeamID:     target.TeamID,
					TargetType: target.TargetType,
					Name:       target.Name,
					Method:     target.Method,
					URL:        api.URL,
					Sort:       target.Sort,
					TypeSort:   target.TypeSort,
					Request: &rao.Request{
						URL:         api.URL,
						Description: api.Description,
						Auth:        &auth,
						Body:        &body,
						Header:      &header,
						Query:       &query,
						Event:       nil,
						Cookie:      nil,
						Resful:      nil,
					},
					Response:       nil,
					Version:        target.Version,
					Description:    api.Description,
					CreatedTimeSec: target.CreatedAt.Unix(),
					UpdatedTimeSec: target.UpdatedAt.Unix(),
				})
			}
		}
	}

	return ret
}
