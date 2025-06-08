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
	memberRepository repositories.MemberRepository
}

func NewMemberUsecase(dr repositories.MemberRepository) MemberUsecase {
	return &MemberUsecaseImp{
		memberRepository: dr,
	}
}

func (cu *MemberUsecaseImp) Read(id string) (model.Member, error) {
	return cu.memberRepository.Read(id)
}

func (cu *MemberUsecaseImp) Create(cm model.Member) (model.Member, error) {
	return cu.memberRepository.Create(cm)
}

func (cu *MemberUsecaseImp) Update(id string, changes map[string]any) (model.Member, error) {
	return cu.memberRepository.Update(id, changes)
}

func (cu *MemberUsecaseImp) Delete(id string) error {
	return cu.memberRepository.Delete(id)
}

func (cu *MemberUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.Member, error) {
	return cu.memberRepository.List(queryOpts)
}
