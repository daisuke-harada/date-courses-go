package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockDeleteDateSpotReviewInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteDateSpotReviewInputPortMockRecorder
}

type MockDeleteDateSpotReviewInputPortMockRecorder struct {
	mock *MockDeleteDateSpotReviewInputPort
}

func NewMockDeleteDateSpotReviewInputPort(ctrl *gomock.Controller) *MockDeleteDateSpotReviewInputPort {
	mock := &MockDeleteDateSpotReviewInputPort{ctrl: ctrl}
	mock.recorder = &MockDeleteDateSpotReviewInputPortMockRecorder{mock}
	return mock
}

func (m *MockDeleteDateSpotReviewInputPort) EXPECT() *MockDeleteDateSpotReviewInputPortMockRecorder {
	return m.recorder
}

func (m *MockDeleteDateSpotReviewInputPort) Execute(arg0 context.Context, arg1 usecase.DeleteDateSpotReviewInput) (*usecase.DeleteDateSpotReviewOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.DeleteDateSpotReviewOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDeleteDateSpotReviewInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockDeleteDateSpotReviewInputPort)(nil).Execute), arg0, arg1)
}
