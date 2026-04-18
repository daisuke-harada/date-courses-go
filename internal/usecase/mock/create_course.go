package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockCreateCourseInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockCreateCourseInputPortMockRecorder
}

type MockCreateCourseInputPortMockRecorder struct {
	mock *MockCreateCourseInputPort
}

func NewMockCreateCourseInputPort(ctrl *gomock.Controller) *MockCreateCourseInputPort {
	mock := &MockCreateCourseInputPort{ctrl: ctrl}
	mock.recorder = &MockCreateCourseInputPortMockRecorder{mock}
	return mock
}

func (m *MockCreateCourseInputPort) EXPECT() *MockCreateCourseInputPortMockRecorder {
	return m.recorder
}

func (m *MockCreateCourseInputPort) Execute(arg0 context.Context, arg1 usecase.CreateCourseInput) (*usecase.CreateCourseOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.CreateCourseOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockCreateCourseInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCreateCourseInputPort)(nil).Execute), arg0, arg1)
}
