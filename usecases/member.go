package usecases

import (
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
)

type MemberUsecase interface {
	Read(id string) (model.Member, error)
	Create(cm model.Member) (model.Member, error)
	Update(id string, changes map[string]any) (model.Member, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Member, error)
}

type MemberUsecaseImp struct {
	memberRepository     repositories.MemberRepository
	membershipRepository repositories.MembershipRepository
}

func NewMemberUsecase(dr repositories.MemberRepository, msr repositories.MembershipRepository) MemberUsecase {
	return &MemberUsecaseImp{
		memberRepository:     dr,
		membershipRepository: msr,
	}
}

func (cu *MemberUsecaseImp) Read(id string) (model.Member, error) {
	return cu.memberRepository.Read(id)
}

func (cu *MemberUsecaseImp) Create(cm model.Member) (model.Member, error) {
	if _, err := cu.membershipRepository.Read(cm.MembershipID); err != nil {
		return model.Member{}, err
	}

	return cu.memberRepository.Create(cm)
}

func (cu *MemberUsecaseImp) Update(id string, changes map[string]any) (model.Member, error) {
	if ms, ok := changes["membership_id"]; ok {
		if msStr, ok := ms.(string); ok {
			if _, err := cu.membershipRepository.Read(msStr); err != nil {
				return model.Member{}, err
			}
		}
	}

	return cu.memberRepository.Update(id, changes)
}

func (cu *MemberUsecaseImp) Delete(id string) error {
	return cu.memberRepository.Delete(id)
}

func (cu *MemberUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Member, error) {
	return cu.memberRepository.List(queryOpts)
}
