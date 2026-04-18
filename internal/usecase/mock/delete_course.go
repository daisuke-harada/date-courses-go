package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockDeleteCourseInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteCourseInputPortMockRecorder
}

type MockDeleteCourseInputPortMockRecorder struct {
	mock *MockDeleteCourseInputPort
}

func NewMockDeleteCourseInputPort(ctrl *gomock.Controller) *MockDeleteCourseInputPort {
	mock := &MockDeleteCourseInputPort{ctrl: ctrl}
	mock.recorder = &MockDeleteCourseInputPortMockRecorder{mock}
	return mock
}

func (m *MockDeleteCourseInputPort) EXPECT() *MockDeleteCourseInputPortMockRecorder {
	return m.recorder
}

func (m *MockDeleteCourseInputPort) Execute(arg0 context.Context, arg1 usecase.DeleteCourseInput) (*usecase.DeleteCourseOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.DeleteCourseOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDeleteCourseInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockDeleteCourseInputPort)(nil).Execute), arg0, arg1)
}
