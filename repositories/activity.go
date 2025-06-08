package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var activityIndex string = "Activity"

type ActivityRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.Activity, error)
	Create(cm model.Activity) (model.Activity, error)
	Update(id string, changes map[string]any) (model.Activity, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Activity, error)
	Index() string
}

type ActivityRepositoryImp struct {
	DB *db.Connection
}

func NewActivityRepository(dbConn *db.Connection) ActivityRepository {
	return &ActivityRepositoryImp{
		DB: dbConn,
	}
}

func (cs *ActivityRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *ActivityRepositoryImp) Index() string {
	return activityIndex
}

func (cs *ActivityRepositoryImp) Read(id string) (model.Activity, error) {
	activity := model.Activity{}
	resMap, err := cs.DB.Read(activityIndex, id, model.Activity{})
	if err != nil {
		return activity, err
	}

	err = utils.Map2Struct(resMap, &activity)
	return activity, err
}

func (cs *ActivityRepositoryImp) Create(cm model.Activity) (model.Activity, error) {
	activity := model.Activity{}
	resMap, err := cs.DB.Create(activityIndex, cm)
	if err != nil {
		return activity, err
	}

	err = utils.Map2Struct(resMap, &activity)
	return activity, err
}

func (cs *ActivityRepositoryImp) Update(id string, changes map[string]any) (model.Activity, error) {
	activity := model.Activity{}
	resMap, err := cs.DB.Update(activityIndex, id, model.Activity{}, changes)
	if err != nil {
		return activity, err
	}

	err = utils.Map2Struct(resMap, &activity)
	return activity, err
}

func (cs *ActivityRepositoryImp) Delete(id string) error {
	return cs.DB.Delete(activityIndex, id, model.Activity{})
}

func (cs *ActivityRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.Activity, error) {
	activitys := []model.Activity{}
	res, err := cs.DB.List(activityIndex, model.Activity{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		activity := model.Activity{}
		err = utils.Map2Struct(v, &activity)
		if err != nil {
			return nil, err
		}

		activitys = append(activitys, activity)
	}

	return activitys, nil
}
