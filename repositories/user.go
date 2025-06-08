package repositories

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/utils"
)

var userIndex string = "User"

type UserRepository interface {
	RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error
	Read(id string) (model.User, error)
	Create(cm model.User) (model.User, error)
	Update(id string, changes map[string]any) (model.User, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.User, error)
	Index() string
}

type UserRepositoryImp struct {
	DB *db.Connection
}

func NewUserRepository(dbConn *db.Connection) UserRepository {
	return &UserRepositoryImp{
		DB: dbConn,
	}
}

func (cs *UserRepositoryImp) RunTransaction(ctx context.Context, f func(tx db.DBTransaction) error) error {
	return cs.DB.RunTransaction(ctx, f)
}

func (cs *UserRepositoryImp) Index() string {
	return userIndex
}

func (cs *UserRepositoryImp) Read(id string) (model.User, error) {
	user := model.User{}
	resMap, err := cs.DB.Read(userIndex, id, model.User{})
	if err != nil {
		return user, err
	}

	err = utils.Map2Struct(resMap, &user)
	return user, nil
}

func (cs *UserRepositoryImp) Create(cm model.User) (model.User, error) {
	user := model.User{}
	resMap, err := cs.DB.Create(userIndex, cm)
	if err != nil {
		return user, err
	}

	err = utils.Map2Struct(resMap, &user)
	return user, nil
}

func (cs *UserRepositoryImp) Update(id string, changes map[string]any) (model.User, error) {
	user := model.User{}
	resMap, err := cs.DB.Update(userIndex, id, model.User{}, changes)
	if err != nil {
		return user, err
	}

	err = utils.Map2Struct(resMap, &user)
	return user, nil
}

func (cs *UserRepositoryImp) Delete(id string) error {
	err := cs.DB.Delete(userIndex, id, model.User{})
	if err != nil {
		return err
	}

	return nil
}

func (cs *UserRepositoryImp) List(queryOpts infrastructure.QueryOpts) ([]model.User, error) {
	users := []model.User{}
	res, err := cs.DB.List(userIndex, model.User{}, queryOpts)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		user := model.User{}
		err = utils.Map2Struct(v, &user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
