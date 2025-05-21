// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/api/presenter"
	"github.com/nakamaFramework/whot-module/cgbdb"
	"github.com/nakamaFramework/whot-module/entity"
	"github.com/nakamaFramework/whot-module/pkg/packager"
	"github.com/nakamaFramework/whot-module/usecase/engine"
	"github.com/nakamaFramework/whot-module/usecase/processor"
	gsm "github.com/nakamaFramework/whot-module/usecase/state_machine"
	"github.com/qmuntal/stateless"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	tickRate = 2
)

// Compile-time check to make sure all required functions are implemented.
var _ runtime.Match = &MatchHandler{}

type MatchHandler struct {
	processor processor.UseCase
	machine   gsm.UseCase
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	//panic("implement me")
	s := state.(*entity.MatchState)
	return s, ""
}

func NewMatchHandler(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) *MatchHandler {
	fmt.Println("new match handler")
	return &MatchHandler{
		processor: processor.NewMatchProcessor(marshaler, unmarshaler, engine.NewWhotPokerEngine()),
		machine:   gsm.NewGameStateMachine(),
	}
}

func (m *MatchHandler) GetState() stateless.State {
	return m.machine.MustState()
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	bet, ok := params["bet"].(int32)
	if !ok {
		logger.Error("invalid match init parameter \"bet\"")
		return nil, 0, ""
	}

	name, ok := params["name"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"name\"")
		//return nil, 0, ""
	}

	password, ok := params["password"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"password\"")
		//return nil, 0, ""
	}

	open := int32(1)
	if password != "" {
		open = 0
	}

	mockCodeCard, _ := params["mock_code_card"].(int32)

	label := &entity.MatchLabel{
		Open:         open,
		Bet:          bet,
		Code:         entity.ModuleName,
		Name:         name,
		Password:     password,
		MaxSize:      entity.MaxPresences,
		MockCodeCard: mockCodeCard,
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, tickRate, ""
	}

	logger.Info("match init label= %s", string(labelJSON))

	matchState := entity.NewMathState(label)
	// init jp treasure
	jpTreasure, _ := cgbdb.GetJackpot(ctx, logger, db, entity.ModuleName)
	if jpTreasure != nil {
		matchState.SetJackpotTreasure(&pb.Jackpot{
			GameCode: jpTreasure.GetGameCode(),
			Chips:    jpTreasure.Chips,
		})
	}
	// fire idle event
	procPkg := packager.NewProcessorPackage(&matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(packager.GetContextWithProcessorPackager(procPkg))

	return &matchState, tickRate, string(labelJSON)
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.MatchState)

	err := m.machine.FireProcessEvent(packager.GetContextWithProcessorPackager(
		packager.NewProcessorPackage(
			s, m.processor,
			logger,
			nk,
			db,
			dispatcher,
			messages,
			ctx),
	))
	if err == presenter.ErrGameFinish {
		logger.Info("match need finish")

		return nil
	}

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
}
