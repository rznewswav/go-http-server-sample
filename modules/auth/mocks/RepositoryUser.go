package mocks

import (
	"newswav/http-server-sample/modules/auth"
	"newswav/http-server-sample/services/utils"
)

type UserRepositoryMock struct {
	utils.PopNextHandler
}

/*
Pass value to be mock-returned by reference. For example:

	repo := GetMockUserRepository(&User {})
*/
func GetMockUserRepository(mockReturnValues ...interface{}) auth.IUserRepository {
	mock := UserRepositoryMock{}
	mock.PrepareMockReturnValue(mockReturnValues...)
	return &mock
}

func (mock *UserRepositoryMock) GetUserByEmail(email string) (auth.SchemaUser, error) {
	nextValue := mock.Next()
	return *(nextValue.(*auth.SchemaUser)), nil
}
