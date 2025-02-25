// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "dragonball-test/internal/character/domain"

	mock "github.com/stretchr/testify/mock"
)

// CharacterRepository is an autogenerated mock type for the CharacterRepository type
type CharacterRepository struct {
	mock.Mock
}

// FindByName provides a mock function with given fields: ctx, name
func (_m *CharacterRepository) FindByName(ctx context.Context, name string) (*domain.Character, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for FindByName")
	}

	var r0 *domain.Character
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.Character, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.Character); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Character)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, character
func (_m *CharacterRepository) Save(ctx context.Context, character *domain.Character) error {
	ret := _m.Called(ctx, character)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Character) error); ok {
		r0 = rf(ctx, character)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCharacterRepository creates a new instance of CharacterRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCharacterRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CharacterRepository {
	mock := &CharacterRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
