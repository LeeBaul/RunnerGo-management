package packer

import (
	"github.com/go-omnibus/proof"
	"go.mongodb.org/mongo-driver/bson"

	"kp-management/internal/pkg/dal/mao"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/rao"
)

func TransSaveTargetReqToMaoAPI(target *rao.SaveTargetReq) *mao.API {
	if target.Request == nil {
		proof.Error("target.request not found request")
		return nil
	}

	header, err := bson.Marshal(target.Request.Header)
	if err != nil {
		proof.Error("target.request.header bson marshal err", proof.WithError(err))
	}

	query, err := bson.Marshal(target.Request.Query)
	if err != nil {
		proof.Errorf("target.request.query bson marshal err", proof.WithError(err))
	}

	body, err := bson.Marshal(target.Request.Body)
	if err != nil {
		proof.Errorf("target.request.body bson marshal err", proof.WithError(err))
	}

	auth, err := bson.Marshal(target.Request.Auth)
	if err != nil {
		proof.Errorf("target.request.auth bson marshal err", proof.WithError(err))
	}

	assert, err := bson.Marshal(mao.Assert{Assert: target.Assert})
	if err != nil {
		proof.Errorf("target.request.assert bson marshal err", proof.WithError(err))
	}

	regex, err := bson.Marshal(mao.Regex{Regex: target.Regex})
	if err != nil {
		proof.Errorf("target.request.regex bson marshal err", proof.WithError(err))
	}

	return &mao.API{
		TargetID:    target.TargetID,
		URL:         target.URL,
		Header:      header,
		Query:       query,
		Body:        body,
		Auth:        auth,
		Description: target.Description,
		Assert:      assert,
		Regex:       regex,
	}
}

func TransTargetToRaoAPIDetail(target *model.Target, api *mao.API) *rao.APIDetail {
	var auth rao.Auth
	if err := bson.Unmarshal(api.Auth, &auth); err != nil {
		proof.Errorf("api.auth bson Unmarshal err", proof.WithError(err))
	}
	var body rao.Body
	if err := bson.Unmarshal(api.Body, &body); err != nil {
		proof.Errorf("api.body bson Unmarshal err", proof.WithError(err))
	}
	var header rao.Header
	if err := bson.Unmarshal(api.Header, &header); err != nil {
		proof.Errorf("api.header bson Unmarshal err", proof.WithError(err))
	}
	var query rao.Query
	if err := bson.Unmarshal(api.Query, &query); err != nil {
		proof.Errorf("api.query bson Unmarshal err", proof.WithError(err))
	}

	var assert mao.Assert
	if err := bson.Unmarshal(api.Assert, &assert); err != nil {
		proof.Errorf("api.assert bson Unmarshal err", proof.WithError(err))
	}

	var regex mao.Regex
	if err := bson.Unmarshal(api.Regex, &regex); err != nil {
		proof.Errorf("api.regex bson Unmarshal err", proof.WithError(err))
	}

	return &rao.APIDetail{
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
		Assert:         assert.Assert,
		Regex:          regex.Regex,
	}
}

func TransTargetsToRaoAPIDetails(targets []*model.Target, apis []*mao.API) []*rao.APIDetail {
	ret := make([]*rao.APIDetail, 0, len(targets))

	for _, target := range targets {
		for _, api := range apis {
			if api.TargetID == target.ID {
				ret = append(ret, TransTargetToRaoAPIDetail(target, api))
			}
		}
	}

	return ret
}
