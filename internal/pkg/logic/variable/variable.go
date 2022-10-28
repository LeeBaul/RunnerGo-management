package variable

import (
	"context"

	"kp-management/internal/pkg/biz/consts"
	"kp-management/internal/pkg/dal"
	"kp-management/internal/pkg/dal/model"
	"kp-management/internal/pkg/dal/query"
	"kp-management/internal/pkg/dal/rao"
	"kp-management/internal/pkg/packer"
)

func SaveVariable(ctx context.Context, req *rao.SaveVariableReq) error {
	tx := query.Use(dal.DB()).Variable

	_, err := tx.WithContext(ctx).Where(tx.ID.Eq(req.VarID)).Assign(
		tx.TeamID.Value(req.TeamID),
		tx.Var.Value(req.Var),
		tx.Val.Value(req.Val),
		tx.Description.Value(req.Description),
	).FirstOrCreate()

	return err
}

func DeleteVariable(ctx context.Context, teamID, varID int64) error {
	tx := query.Use(dal.DB()).Variable

	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.ID.Eq(varID)).Delete()
	return err
}

func ListGlobalVariables(ctx context.Context, teamID int64, limit, offset int) ([]*rao.Variable, int64, error) {
	tx := query.Use(dal.DB()).Variable

	v, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.Type.Eq(consts.VariableTypeGlobal)).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	return packer.TransModelVariablesToRaoVariables(v), cnt, nil
}

func SyncGlobalVariables(ctx context.Context, teamID int64, variables []*rao.Variable) error {
	vs := packer.TransRaoVariablesToModelVariables(teamID, variables)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.Variable.WithContext(ctx).Where(tx.Variable.TeamID.Eq(teamID), tx.Variable.Type.Eq(consts.VariableTypeGlobal)).Unscoped().Delete(); err != nil {
			return err
		}

		return tx.Variable.WithContext(ctx).CreateInBatches(vs, 10)
	})
}

func ListSceneVariables(ctx context.Context, teamID, sceneID int64, limit, offset int) ([]*rao.Variable, int64, error) {
	tx := dal.GetQuery().Variable

	v, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.SceneID.Eq(sceneID), tx.Type.Eq(consts.VariableTypeScene)).Limit(limit).Offset(offset).Find()

	//v, cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.SceneID.Eq(sceneID), tx.Type.Eq(consts.VariableTypeScene)).FindByPage(offset, limit)
	if err != nil {
		return nil, 0, err
	}

	cnt, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.SceneID.Eq(sceneID), tx.Type.Eq(consts.VariableTypeScene)).Count()
	if err != nil {
		return nil, 0, err
	}

	return packer.TransModelVariablesToRaoVariables(v), cnt, nil
}

func SyncSceneVariables(ctx context.Context, teamID, sceneID int64, variables []*rao.Variable) error {
	vs := packer.TransSceneRaoVariablesToModelVariables(teamID, sceneID, variables)

	return query.Use(dal.DB()).Transaction(func(tx *query.Query) error {
		if _, err := tx.Variable.WithContext(ctx).Where(tx.Variable.TeamID.Eq(teamID), tx.Variable.Type.Eq(consts.VariableTypeScene)).Unscoped().Delete(); err != nil {
			return err
		}

		return tx.Variable.WithContext(ctx).CreateInBatches(vs, 10)
	})
}

func ImportSceneVariables(ctx context.Context, req *rao.ImportVariablesReq, userID int64) error {

	tx := dal.GetQuery().VariableImport
	return tx.WithContext(ctx).Create(&model.VariableImport{
		TeamID:     req.TeamID,
		SceneID:    req.SceneID,
		Name:       req.Name,
		URL:        req.URL,
		UploaderID: userID,
	})
}

func DeleteImportSceneVariables(ctx context.Context, req *rao.DeleteImportSceneVariablesReq) error {
	tx := dal.GetQuery().VariableImport
	_, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(req.TeamID), tx.SceneID.Eq(req.SceneID), tx.Name.Eq(req.Name)).Delete()
	return err
}

func ListImportSceneVariables(ctx context.Context, teamID, sceneID int64) ([]*rao.Import, error) {
	tx := dal.GetQuery().VariableImport
	vi, err := tx.WithContext(ctx).Where(tx.TeamID.Eq(teamID), tx.SceneID.Eq(sceneID)).Limit(5).Find()
	if err != nil {
		return nil, err
	}

	return packer.TransImportVariablesToRaoImportVariables(vi), nil
}
