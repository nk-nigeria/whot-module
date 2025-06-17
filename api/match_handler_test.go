package api

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	_ "github.com/lib/pq"
	"github.com/nk-nigeria/whot-module/mock"
	"github.com/nk-nigeria/whot-module/pkg/log"
	"google.golang.org/protobuf/proto"
)

func TestMatch(t *testing.T) {
	t.Logf("test match")

	//userIds := []string{"user1", "user2", "user3"}
	//pairs := combinations.Combinations(userIds, 2)
	//log.GetLogger().Info("combination %v", pairs)
	//for _, pair := range pairs {
	//	t.Logf("compare %v with %v", pair[0], pair[1])
	//}

	marshaler := &proto.MarshalOptions{}
	unmarshaler := &proto.UnmarshalOptions{
		DiscardUnknown: false,
	}

	m := NewMatchHandler(marshaler, unmarshaler)
	var params = make(map[string]interface{})
	params["bet"] = int32(100)
	params["name"] = "name"
	params["password"] = "password"

	logger := log.GetLogger()
	dispatcher := mock.MockDispatcher{}
	nk := mock.MockModule{}
	connStr := "postgresql://postgres:localdb@127.0.0.1/nakama?sslmode=disable"
	mdb, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Logf("failed to open db: %v", err)
	}

	if err := mdb.Ping(); err != nil {
		t.Logf("failed to connect to db: %v", err)
	}
	defer mdb.Close()
	s, _, _ := m.MatchInit(context.Background(), logger, mdb, nil, params)

	ctx := context.TODO()

	// mock event routine
	var stop = make(chan bool)
	go func() {
		t.Logf("start mock loop")
		for i := 0; i < 2*15; i++ {
			t.Logf("log %d", i)
			time.Sleep(time.Millisecond * 500)
			m.MatchLoop(ctx, logger, nil, &nk, dispatcher, 0, s, nil)
		}

		t.Logf("current state %v", m.GetState())

		stop <- true
	}()

	go func() {
		t.Logf("start mock join leave")
		presences := make([]runtime.Presence, 1)
		presences[0] = &mock.MockPresence{
			UserId: "user1",
		}

		m.MatchJoin(ctx, logger, mdb, &nk, dispatcher, 0, s, presences)

		time.Sleep(time.Second * 2)
		presences = make([]runtime.Presence, 1)
		presences[0] = &mock.MockPresence{
			UserId: "user2",
		}
		m.MatchJoin(ctx, logger, mdb, &nk, dispatcher, 0, s, presences)
	}()

	t.Logf("wait for finish")
	<-stop
	t.Logf("wait for finish done")
}
