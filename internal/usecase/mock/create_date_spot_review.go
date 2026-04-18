package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockCreateDateSpotReviewInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockCreateDateSpotReviewInputPortMockRecorder
}

type MockCreateDateSpotReviewInputPortMockRecorder struct {
	mock *MockCreateDateSpotReviewInputPort
}

func NewMockCreateDateSpotReviewInputPort(ctrl *gomock.Controller) *MockCreateDateSpotReviewInputPort {
	mock := &MockCreateDateSpotReviewInputPort{ctrl: ctrl}
	mock.recorder = &MockCreateDateSpotReviewInputPortMockRecorder{mock}
	return mock
}

func (m *MockCreateDateSpotReviewInputPort) EXPECT() *MockCreateDateSpotReviewInputPortMockRecorder {
	return m.recorder
}

func (m *MockCreateDateSpotReviewInputPort) Execute(arg0 context.Context, arg1 usecase.CreateDateSpotReviewInput) (*usecase.CreateDateSpotReviewOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.CreateDateSpotReviewOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockCreateDateSpotReviewInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateDateSpotReviewInputPort)(nil).Execute), arg0, arg1)
}
