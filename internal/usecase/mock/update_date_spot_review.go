package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockUpdateDateSpotReviewInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateDateSpotReviewInputPortMockRecorder
}

type MockUpdateDateSpotReviewInputPortMockRecorder struct {
	mock *MockUpdateDateSpotReviewInputPort
}

func NewMockUpdateDateSpotReviewInputPort(ctrl *gomock.Controller) *MockUpdateDateSpotReviewInputPort {
	mock := &MockUpdateDateSpotReviewInputPort{ctrl: ctrl}
	mock.recorder = &MockUpdateDateSpotReviewInputPortMockRecorder{mock}
	return mock
}

func (m *MockUpdateDateSpotReviewInputPort) EXPECT() *MockUpdateDateSpotReviewInputPortMockRecorder {
	return m.recorder
}

func (m *MockUpdateDateSpotReviewInputPort) Execute(arg0 context.Context, arg1 usecase.UpdateDateSpotReviewInput) (*usecase.UpdateDateSpotReviewOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.UpdateDateSpotReviewOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockUpdateDateSpotReviewInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockUpdateDateSpotReviewInputPort)(nil).Execute), arg0, arg1)
}
