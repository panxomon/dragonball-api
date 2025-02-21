package create_test

import (
	"context"
	"dragonball-test/internal/character/application/create"
	"dragonball-test/internal/character/domain"
	"dragonball-test/internal/character/domain/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CreateCharacterService(t *testing.T) {
	type test struct {
		name           string
		repositoryMock func() *mocks.CharacterRepository
		expectedResult *domain.Character
		expectedError  error
	}

	cases := []test{
		{
			name: "when character is found in database",
			repositoryMock: func() *mocks.CharacterRepository {
				repo := new(mocks.CharacterRepository)
				repo.On("FindByName", mock.Anything, "Goku").Return(&domain.Character{
					ID:    1,
					Name:  "Goku",
					Race:  "Saiyan",
					Image: "goku_image.png",
				}, nil)
				return repo
			},
			expectedResult: &domain.Character{
				ID:    1,
				Name:  "Goku",
				Race:  "Saiyan",
				Image: "goku_image.png",
			},
			expectedError: nil,
		},
		{
			name: "when character is not found in database and saving fails",
			repositoryMock: func() *mocks.CharacterRepository {
				repo := new(mocks.CharacterRepository)
				repo.On("FindByName", mock.Anything, "Goku").Return(nil, nil)
				repo.On("Save", mock.Anything, mock.Anything).Return(errors.New("failed to save"))
				return repo
			},
			expectedResult: nil,
			expectedError:  errors.New("failed to save"),
		},
		{
			name: "when character is not found in database and saving succeeds",
			repositoryMock: func() *mocks.CharacterRepository {
				repo := new(mocks.CharacterRepository)
				repo.On("FindByName", mock.Anything, "Goku").Return(nil, nil)
				repo.On("Save", mock.Anything, mock.Anything).Return(nil)
				return repo
			},
			expectedResult: &domain.Character{
				Name:  "Goku",
				Race:  "",
				Image: "",
			},
			expectedError: nil,
		},
		{
			name: "when repository returns error on finding character",
			repositoryMock: func() *mocks.CharacterRepository {
				repo := new(mocks.CharacterRepository)
				repo.On("FindByName", mock.Anything, "Goku").Return(nil, errors.New("repository error"))
				return repo
			},
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := tt.repositoryMock()
			service := create.NewCreateCharacterService(repositoryMock)

			character, err := service.CreateCharacter(context.Background(), "Goku")

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, character)
			}
			repositoryMock.AssertExpectations(t)
		})
	}
}
