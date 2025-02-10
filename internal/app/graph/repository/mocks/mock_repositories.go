// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -destination=mocks/mock_repositories.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "ozon_test_task/internal/app/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockPostRepo is a mock of PostRepo interface.
type MockPostRepo struct {
	ctrl     *gomock.Controller
	recorder *MockPostRepoMockRecorder
	isgomock struct{}
}

// MockPostRepoMockRecorder is the mock recorder for MockPostRepo.
type MockPostRepoMockRecorder struct {
	mock *MockPostRepo
}

// NewMockPostRepo creates a new mock instance.
func NewMockPostRepo(ctrl *gomock.Controller) *MockPostRepo {
	mock := &MockPostRepo{ctrl: ctrl}
	mock.recorder = &MockPostRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostRepo) EXPECT() *MockPostRepoMockRecorder {
	return m.recorder
}

// AddPost mocks base method.
func (m *MockPostRepo) AddPost(ctx context.Context, post *models.Post) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPost", ctx, post)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddPost indicates an expected call of AddPost.
func (mr *MockPostRepoMockRecorder) AddPost(ctx, post any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPost", reflect.TypeOf((*MockPostRepo)(nil).AddPost), ctx, post)
}

// GetPostByID mocks base method.
func (m *MockPostRepo) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostByID", ctx, postID)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostByID indicates an expected call of GetPostByID.
func (mr *MockPostRepoMockRecorder) GetPostByID(ctx, postID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostByID", reflect.TypeOf((*MockPostRepo)(nil).GetPostByID), ctx, postID)
}

// GetPosts mocks base method.
func (m *MockPostRepo) GetPosts(ctx context.Context, limit, after int) ([]*models.Post, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosts", ctx, limit, after)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPosts indicates an expected call of GetPosts.
func (mr *MockPostRepoMockRecorder) GetPosts(ctx, limit, after any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosts", reflect.TypeOf((*MockPostRepo)(nil).GetPosts), ctx, limit, after)
}

// SetCommentsAllowed mocks base method.
func (m *MockPostRepo) SetCommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCommentsAllowed", ctx, postID, commentsAllowed)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCommentsAllowed indicates an expected call of SetCommentsAllowed.
func (mr *MockPostRepoMockRecorder) SetCommentsAllowed(ctx, postID, commentsAllowed any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCommentsAllowed", reflect.TypeOf((*MockPostRepo)(nil).SetCommentsAllowed), ctx, postID, commentsAllowed)
}

// MockCommentRepo is a mock of CommentRepo interface.
type MockCommentRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepoMockRecorder
	isgomock struct{}
}

// MockCommentRepoMockRecorder is the mock recorder for MockCommentRepo.
type MockCommentRepoMockRecorder struct {
	mock *MockCommentRepo
}

// NewMockCommentRepo creates a new mock instance.
func NewMockCommentRepo(ctrl *gomock.Controller) *MockCommentRepo {
	mock := &MockCommentRepo{ctrl: ctrl}
	mock.recorder = &MockCommentRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentRepo) EXPECT() *MockCommentRepoMockRecorder {
	return m.recorder
}

// AddComment mocks base method.
func (m *MockCommentRepo) AddComment(ctx context.Context, comment *models.Comment) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", ctx, comment)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddComment indicates an expected call of AddComment.
func (mr *MockCommentRepoMockRecorder) AddComment(ctx, comment any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockCommentRepo)(nil).AddComment), ctx, comment)
}

// GetCommentsByPostID mocks base method.
func (m *MockCommentRepo) GetCommentsByPostID(ctx context.Context, postID, limit, after int) ([]*models.Comment, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentsByPostID", ctx, postID, limit, after)
	ret0, _ := ret[0].([]*models.Comment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCommentsByPostID indicates an expected call of GetCommentsByPostID.
func (mr *MockCommentRepoMockRecorder) GetCommentsByPostID(ctx, postID, limit, after any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentsByPostID", reflect.TypeOf((*MockCommentRepo)(nil).GetCommentsByPostID), ctx, postID, limit, after)
}

// GetReplaysByCommentID mocks base method.
func (m *MockCommentRepo) GetReplaysByCommentID(ctx context.Context, commentID, limit, after int) ([]*models.Comment, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReplaysByCommentID", ctx, commentID, limit, after)
	ret0, _ := ret[0].([]*models.Comment)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReplaysByCommentID indicates an expected call of GetReplaysByCommentID.
func (mr *MockCommentRepoMockRecorder) GetReplaysByCommentID(ctx, commentID, limit, after any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReplaysByCommentID", reflect.TypeOf((*MockCommentRepo)(nil).GetReplaysByCommentID), ctx, commentID, limit, after)
}

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
	isgomock struct{}
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockUserRepo) AddUser(ctx context.Context, user *models.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockUserRepoMockRecorder) AddUser(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUserRepo)(nil).AddUser), ctx, user)
}

// GetUserByID mocks base method.
func (m *MockUserRepo) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserRepoMockRecorder) GetUserByID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserRepo)(nil).GetUserByID), ctx, userID)
}

// GetUserByLoginWithCred mocks base method.
func (m *MockUserRepo) GetUserByLoginWithCred(ctx context.Context, login string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLoginWithCred", ctx, login)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLoginWithCred indicates an expected call of GetUserByLoginWithCred.
func (mr *MockUserRepoMockRecorder) GetUserByLoginWithCred(ctx, login any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLoginWithCred", reflect.TypeOf((*MockUserRepo)(nil).GetUserByLoginWithCred), ctx, login)
}
