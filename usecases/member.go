package usecases

import (
	"fmt"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
	"net/smtp"
)

type MemberUsecase interface {
	Read(id string) (model.Member, error)
	Create(cm model.Member) (model.Member, error)
	Update(id string, changes map[string]any) (model.Member, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Member, error)
	SendEmail(host, sender, password string, port int, receiver, subject, body string) error
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

func (cu *MemberUsecaseImp) SendEmail(host, sender, password string, port int, receiver, subject, body string) error {
	auth := smtp.PlainAuth("", sender, password, host)

	header := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := fmt.Sprintf("Subject: %s\n%s\n\n%s", subject, header, body)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		sender,
		[]string{receiver},
		[]byte(msg),
	)

	return err
}
