package mock

import "github.com/heroiclabs/nakama-common/runtime"

type MockPresence struct {
	UserId string
}

func (m *MockPresence) GetHidden() bool {
	panic("implement me")
}

func (m *MockPresence) GetPersistence() bool {
	panic("implement me")
}

func (m *MockPresence) GetUsername() string {
	return "username1"
}

func (m *MockPresence) GetStatus() string {
	panic("implement me")
}

func (m *MockPresence) GetReason() runtime.PresenceReason {
	panic("implement me")
}

func (m *MockPresence) GetUserId() string {
	return m.UserId
}
func (m *MockPresence) GetSessionId() string {
	return ""
}

func (m *MockPresence) GetNodeId() string {
	return ""
}
