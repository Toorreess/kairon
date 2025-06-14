package usecases

import (
	"context"
	"fmt"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
	"slices"
)

type ActivityUsecase interface {
	Read(id string) (model.Activity, error)
	Create(cm model.Activity) (model.Activity, error)
	Update(id string, changes map[string]any) (model.Activity, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Activity, error)
	Reserve(memberID, activityID string) error
}

type ActivityUsecaseImp struct {
	activityRepository repositories.ActivityRepository
	memberRepository   repositories.MemberRepository
}

func NewActivityUsecase(dr repositories.ActivityRepository, mr repositories.MemberRepository) ActivityUsecase {
	return &ActivityUsecaseImp{
		activityRepository: dr,
		memberRepository:   mr,
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

func (cu *ActivityUsecaseImp) Reserve(memberID, activityID string) error {
	cm, err := cu.memberRepository.Read(memberID)
	if err != nil {
		return err
	}

	if slices.Contains(cm.ActivityList, activityID) {
		return fmt.Errorf("member already has this activity reserved")
	}

	am, err := cu.activityRepository.Read(activityID)
	if err != nil {
		return err
	}

	if am.MaxCapacity <= 0 {
		return fmt.Errorf("activity has no available capacity")
	}

	err = cu.activityRepository.RunTransaction(context.Background(), func(tx db.DBTransaction) error {
		changesActivity := map[string]any{
			"max_capacity": am.MaxCapacity - 1,
		}
		if _, err := cu.activityRepository.Update(activityID, changesActivity); err != nil {
			return fmt.Errorf("failed to update activity capacity: %w", err)
		}

		changesMember := map[string]any{
			"activity_list": append(cm.ActivityList, activityID),
		}
		if _, err := cu.memberRepository.Update(memberID, changesMember); err != nil {
			return fmt.Errorf("failed to update member's activity list: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to complete reservation: %w", err)
	}

	return nil
}
