package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var memberIndex string = "Member"

type MemberRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.Member, error)
	Create(cm model.Member) (model.Member, error)
	Update(id string, changes map[string]any) (model.Member, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Member, error)
	Index() string
}

type MemberRepositoryImp struct {
	DB *db.Connection
}

func NewMemberRepository(dbConn *db.Connection) MemberRepository {
	return &MemberRepositoryImp{
		DB: dbConn,
	}
}

func (cs *MemberRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *MemberRepositoryImp) Index() string {
	return memberIndex
}

func (cs *MemberRepositoryImp) Read(id string) (model.Member, error) {
	member := model.Member{}
	resMap, err := cs.DB.Read(memberIndex, id, model.Member{})
	if err != nil {
		return member, err
	}

	err = utils.Map2Struct(resMap, &member)
	return member, err
}

func (cs *MemberRepositoryImp) Create(cm model.Member) (model.Member, error) {
	member := model.Member{}
	resMap, err := cs.DB.Create(memberIndex, cm)
	if err != nil {
		return member, err
	}

	err = utils.Map2Struct(resMap, &member)
	return member, err
}

func (cs *MemberRepositoryImp) Update(id string, changes map[string]any) (model.Member, error) {
	member := model.Member{}
	resMap, err := cs.DB.Update(memberIndex, id, model.Member{}, changes)
	if err != nil {
		return member, err
	}

	err = utils.Map2Struct(resMap, &member)
	return member, err
}

func (cs *MemberRepositoryImp) Delete(id string) error {
	return cs.DB.Delete(memberIndex, id, model.Member{})
}

func (cs *MemberRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.Member, error) {
	members := []model.Member{}
	res, err := cs.DB.List(memberIndex, model.Member{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		member := model.Member{}
		err = utils.Map2Struct(v, &member)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
