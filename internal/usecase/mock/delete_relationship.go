package usecasemock

import (
	"context"
	"reflect"

	"github.com/daisuke-harada/date-courses-go/internal/usecase"
	"go.uber.org/mock/gomock"
)

type MockDeleteRelationshipInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteRelationshipInputPortMockRecorder
}

type MockDeleteRelationshipInputPortMockRecorder struct {
	mock *MockDeleteRelationshipInputPort
}

func NewMockDeleteRelationshipInputPort(ctrl *gomock.Controller) *MockDeleteRelationshipInputPort {
	mock := &MockDeleteRelationshipInputPort{ctrl: ctrl}
	mock.recorder = &MockDeleteRelationshipInputPortMockRecorder{mock}
	return mock
}

func (m *MockDeleteRelationshipInputPort) EXPECT() *MockDeleteRelationshipInputPortMockRecorder {
	return m.recorder
}

func (m *MockDeleteRelationshipInputPort) Execute(arg0 context.Context, arg1 usecase.DeleteRelationshipInput) (*usecase.DeleteRelationshipOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*usecase.DeleteRelationshipOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDeleteRelationshipInputPortMockRecorder) Execute(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockDeleteRelationshipInputPort)(nil).Execute), arg0, arg1)
}
