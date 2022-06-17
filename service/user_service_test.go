package service

import (
	"context"
	"entry-task-rpc/pkg/rpc"
	"testing"

	"github.com/stretchr/testify/mock" //nolint:goimports
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) InsertUser(string, string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) CheckAuth(string, string) (bool, string, string, error) {
	args := m.Called()
	return args.Bool(0), args.String(1), args.String(2), args.Error(3)
}

func (m *MockDB) GetProfile(string) (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockDB) EditProfile(string, string, string) error {
	args := m.Called()
	return args.Error(0)
}

func TestUserInfo_Register(t *testing.T) {
	in := rpc.RegisterRequest{
		RequestID: "aubsjhbjakbkj",
		Name:      "test01",
		Password:  "123456",
	}

	DB := new(MockDB)
	DB.On("InsertUser").Return(nil)
	user := UserService{DB: DB}
	_, res := user.Register(context.Background(), &in)
	if res != nil {
		t.Errorf("Register failed, res:%s", res)
	}
}

func TestUserInfo_Auth(t *testing.T) {
	in := rpc.AuthRequest{
		RequestID: "aubsjhbjakbkj",
		Name:      "test01",
		Password:  "123456",
	}

	DB := new(MockDB)
	DB.On("CheckAuth").Return(true, "test", "123456", nil)
	user := UserService{DB: DB}
	_, res := user.Auth(context.Background(), &in)
	if res != nil {
		t.Errorf("Register failed, res:%s", res)
	}
}

func TestUserInfo_EditProfile(t *testing.T) {
	in := rpc.EditProfileRequest{
		RequestID: "aubsjhbjakbkj",
		Name:      "test01",
		Nickname:  "123456",
	}

	DB := new(MockDB)
	DB.On("EditProfile").Return(nil)
	user := UserService{DB: DB}
	_, res := user.EditProfile(context.Background(), &in)
	if res != nil {
		t.Errorf("Register failed, res:%s", res)
	}
}

func TestUserInfo_GetProfile(t *testing.T) {
	in := rpc.GetProfileRequest{
		RequestID: "aubsjhbjakbkj",
		Name:      "test01",
	}

	DB := new(MockDB)
	DB.On("GetProfile").Return("123", "456", nil)
	user := UserService{DB: DB}
	_, res := user.GetProfile(context.Background(), &in)
	if res != nil {
		t.Errorf("Register failed, res:%s", res)
	}
}
