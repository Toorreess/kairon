package usecases

import (
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
)

type ActivityUsecase interface {
	Read(id string) (model.Activity, error)
	Create(cm model.Activity) (model.Activity, error)
	Update(id string, changes map[string]any) (model.Activity, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Activity, error)
}

type ActivityUsecaseImp struct {
	activityRepository repositories.ActivityRepository
}

func NewActivityUsecase(dr repositories.ActivityRepository) ActivityUsecase {
	return &ActivityUsecaseImp{
		activityRepository: dr,
	}
}

func (cu *ActivityUsecaseImp) Read(id string) (model.Activity, error) {
	return cu.activityRepository.Read(id)
}

func (cu *ActivityUsecaseImp) Create(cm model.Activity) (model.Activity, error) {
	return cu.activityRepository.Create(cm)
}

func (cu *ActivityUsecaseImp) Update(id string, changes map[string]any) (model.Activity, error) {
	return cu.activityRepository.Update(id, changes)
}

func (cu *ActivityUsecaseImp) Delete(id string) error {
	return cu.activityRepository.Delete(id)
}

func (cu *ActivityUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Activity, error) {
	return cu.activityRepository.List(queryOpts)
}
