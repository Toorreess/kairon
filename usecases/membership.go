package usecases

import (
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
)

type MembershipUsecase interface {
	Read(id string) (model.Membership, error)
	Create(cm model.Membership) (model.Membership, error)
	Update(id string, changes map[string]any) (model.Membership, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Membership, error)
}

type MembershipUsecaseImp struct {
	membershipRepository repositories.MembershipRepository
}

func NewMembershipUsecase(dr repositories.MembershipRepository) MembershipUsecase {
	return &MembershipUsecaseImp{
		membershipRepository: dr,
	}
}

func (cu *MembershipUsecaseImp) Read(id string) (model.Membership, error) {
	return cu.membershipRepository.Read(id)
}

func (cu *MembershipUsecaseImp) Create(cm model.Membership) (model.Membership, error) {
	return cu.membershipRepository.Create(cm)
}

func (cu *MembershipUsecaseImp) Update(id string, changes map[string]any) (model.Membership, error) {
	return cu.membershipRepository.Update(id, changes)
}

func (cu *MembershipUsecaseImp) Delete(id string) error {
	return cu.membershipRepository.Delete(id)
}

func (cu *MembershipUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Membership, error) {
	return cu.membershipRepository.List(queryOpts)
}
