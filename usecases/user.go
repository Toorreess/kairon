package usecases

import (
	"context"
	db "kairon/adapters/database"
	"kairon/cmd/api/infrastructure"
	"kairon/domain/model"
	"kairon/repositories"
	"kairon/utils"
	"fmt"

	"firebase.google.com/go/auth"
)

type UserUsecase interface {
	Read(id string) (model.User, error)
	Create(cm model.User) (model.User, error)
	Update(id string, changes map[string]any) (model.User, error)
	Delete(id string) error
	List(queryOpts infrastructure.QueryOpts) ([]model.User, error)
}

type UserUsecaseImp struct {
	userRepository repositories.UserRepository
	authClient     *auth.Client
}

func NewUserUsecase(ur repositories.UserRepository, ac *auth.Client) UserUsecase {
	return &UserUsecaseImp{
		userRepository: ur,
		authClient:     ac,
	}
}

func (cu *UserUsecaseImp) Read(id string) (model.User, error) {
	return cu.userRepository.Read(id)
}

func (cu *UserUsecaseImp) Create(cm model.User) (model.User, error) {
	var user model.User

	ctx := context.Background()
	err := cu.userRepository.RunTransaction(context.Background(), func(tx db.DBTransaction) error {
		fbUser, err := cu.authClient.GetUserByEmail(ctx, cm.Email)
		if err != nil {
			params := (&auth.UserToCreate{}).
				Email(cm.Email).
				EmailVerified(true).
				DisplayName(cm.Name).
				Disabled(false)

			fbUser, err = cu.authClient.CreateUser(ctx, params)
			if err != nil {
				return err
			}
		}

		if err := cu.authClient.SetCustomUserClaims(ctx, fbUser.UID, cm.UserClaims()); err != nil {
			return err
		}

		cmr, err := tx.CreateWithID(cu.userRepository.Index(), fbUser.UID, cm)
		if err != nil {
			return err
		}

		return utils.Map2Struct(cmr, &user)
	})

	return user, err
}

func (cu *UserUsecaseImp) Update(id string, changes map[string]any) (model.User, error) {
	var user model.User
	var param auth.UserToUpdate

	if email, ok := changes["email"]; ok {
		qo := infrastructure.QueryOpts{
			QueryString: fmt.Sprintf("email:\"%s\"", email),
			Limit:       1,
			Offset:      0,
		}

		usrList, err := cu.userRepository.List(qo)
		if err != nil {
			return user, err
		}

		if len(usrList) > 0 && usrList[0].ID != id {
			return user, fmt.Errorf("an user with this email already exists")
		}

		param.Email(email.(string))
	}

	if name, ok := changes["name"]; ok {
		param.DisplayName(name.(string))
	}

	res, err := cu.userRepository.Update(id, changes)
	if err != nil {
		return user, fmt.Errorf("error updating user %s", id)
	}

	param.CustomClaims(res.UserClaims())

	ctx := context.Background()
	if _, err := cu.authClient.UpdateUser(ctx, id, &param); err != nil {
		return user, err
	}

	return res, err
}

func (cu *UserUsecaseImp) Delete(id string) error {
	if err := cu.userRepository.Delete(id); err != nil {
		return err
	}

	ctx := context.Background()
	params := (&auth.UserToUpdate{}).Disabled(true)
	_, err := cu.authClient.UpdateUser(ctx, id, params)

	return err
}

func (cu *UserUsecaseImp) List(queryOpts infrastructure.QueryOpts) ([]model.User, error) {
	return cu.userRepository.List(queryOpts)
}
