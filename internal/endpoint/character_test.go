package endpoint_test

import (
	"bytes"
	"context"
	"dragonball-test/internal/character/application"
	"dragonball-test/internal/character/application/create"
	"dragonball-test/internal/character/domain"
	"dragonball-test/internal/endpoint"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCreateCharacterService struct {
	mock.Mock
}

func (m *MockCreateCharacterService) Handle(ctx context.Context, cmd create.CreateCharacter) (*domain.Character, error) {
	args := m.Called(ctx, cmd)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Character), nil
}

func Test_CharacterEndpoint_Invoke(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type test struct {
		name           string
		createService  *MockCreateCharacterService
		requestBody    string
		mockResponse   *domain.Character
		mockError      error
		expectedStatus int
		expectedBody   string
	}

	cases := []test{
		{
			name: "valid request - character created successfully",
			createService: func() *MockCreateCharacterService {
				mockService := new(MockCreateCharacterService)
				mockService.On("Handle", mock.Anything, create.CreateCharacter{Name: "Goku"}).Return(&domain.Character{
					ID:    1,
					Name:  "Goku",
					Race:  "Saiyan",
					Image: "goku_image.png",
				}, nil)
				return mockService
			}(),
			requestBody: `{"name": "Goku"}`,
			mockResponse: &domain.Character{
				ID:    1,
				Name:  "Goku",
				Race:  "Saiyan",
				Image: "goku_image.png",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"character":{"ID":1, "Name":"Goku", "image":"goku_image.png", "race":"Saiyan"}}`,
		},
		{
			name: "error on character creation",
			createService: func() *MockCreateCharacterService {
				mockService := new(MockCreateCharacterService)
				mockService.On("Handle", mock.Anything, create.CreateCharacter{Name: "Goku"}).Return(nil, assert.AnError)
				return mockService
			}(),
			requestBody:    `{"name": "Goku"}`,
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to fetch character"}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			characterApp := &application.App{
				Commands: application.Commands{
					CreateCharacter: tt.createService,
				},
			}
			ep := endpoint.NewCharacterEndpoint(characterApp)
			router.POST("/characters", ep.Invoke)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/characters", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
