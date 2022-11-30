package usecase

import (
	"fmt"
	"testing"

	"assignment-golang-backend/internal/custom_error"
	"assignment-golang-backend/internal/entity"
	"assignment-golang-backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	NewUserService(mocks.NewIUserRepository(t))
}

func Test_userService_FindByID(t *testing.T) {
	mockUser := &entity.User{
		Base: entity.Base{
			ID: 1,
		},
		Name:  "name",
		Email: "email@email.com",
	}
	tests := []struct {
		name           string
		userRepository *mocks.IUserRepository
		mock           func(ur *mocks.IUserRepository)
		id             int
		want           *entity.User
		wantErr        bool
		expectedErr    error
	}{
		{
			name:           "Error | No User found from repository",
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ur *mocks.IUserRepository) {
				ur.On("FindByID", 1).Return(nil, 0, nil)
			},
			id:          1,
			want:        nil,
			wantErr:     true,
			expectedErr: &custom_error.NoDataFound{DataType: "user"},
		},
		{
			name:           "Error | Other errors from repository",
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ur *mocks.IUserRepository) {
				ur.On("FindByID", 1).Return(nil, 1, fmt.Errorf("error"))
			},
			id:          1,
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("error"),
		},
		{
			name:           "Success",
			userRepository: mocks.NewIUserRepository(t),
			mock: func(ur *mocks.IUserRepository) {
				ur.On("FindByID", 1).Return(mockUser, 1, nil)
			},
			id:          1,
			want:        mockUser,
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				userRepository: tt.userRepository,
			}

			tt.mock(tt.userRepository)

			got, err := s.FindByID(tt.id)

			if !tt.wantErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
