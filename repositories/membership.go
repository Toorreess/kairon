package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var membershipIndex string = "Membership"

type MembershipRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.Membership, error)
	Create(cm model.Membership) (model.Membership, error)
	Update(id string, changes map[string]any) (model.Membership, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.Membership, error)
	Index() string
}

type MembershipRepositoryImp struct {
	DB *db.Connection
}

func NewMembershipRepository(dbConn *db.Connection) MembershipRepository {
	return &MembershipRepositoryImp{
		DB: dbConn,
	}
}

func (cs *MembershipRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *MembershipRepositoryImp) Index() string {
	return membershipIndex
}

func (cs *MembershipRepositoryImp) Read(id string) (model.Membership, error) {
	membership := model.Membership{}
	resMap, err := cs.DB.Read(membershipIndex, id, model.Membership{})
	if err != nil {
		return membership, err
	}

	err = utils.Map2Struct(resMap, &membership)
	return membership, err
}

func (cs *MembershipRepositoryImp) Create(cm model.Membership) (model.Membership, error) {
	membership := model.Membership{}
	resMap, err := cs.DB.Create(membershipIndex, cm)
	if err != nil {
		return membership, err
	}

	err = utils.Map2Struct(resMap, &membership)
	return membership, err
}

func (cs *MembershipRepositoryImp) Update(id string, changes map[string]any) (model.Membership, error) {
	membership := model.Membership{}
	resMap, err := cs.DB.Update(membershipIndex, id, model.Membership{}, changes)
	if err != nil {
		return membership, err
	}

	err = utils.Map2Struct(resMap, &membership)
	return membership, err
}

func (cs *MembershipRepositoryImp) Delete(id string) error {
	return cs.DB.Delete(membershipIndex, id, model.Membership{})
}

func (cs *MembershipRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.Membership, error) {
	memberships := []model.Membership{}
	res, err := cs.DB.List(membershipIndex, model.Membership{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		membership := model.Membership{}
		err = utils.Map2Struct(v, &membership)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, nil
}
