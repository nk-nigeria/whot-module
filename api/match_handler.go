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

	"github.com/heroiclabs/nakama-common/runtime"
	pb1 "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/whot-module/api/presenter"
	"github.com/nk-nigeria/whot-module/constant"
	"github.com/nk-nigeria/whot-module/entity"
	"github.com/nk-nigeria/whot-module/pkg/packager"
	"github.com/nk-nigeria/whot-module/usecase/engine"
	"github.com/nk-nigeria/whot-module/usecase/processor"
	gsm "github.com/nk-nigeria/whot-module/usecase/state_machine"
	"github.com/qmuntal/stateless"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

func NewMatchHandler(marshaler *proto.MarshalOptions, unmarshaler *proto.UnmarshalOptions) *MatchHandler {
	return &MatchHandler{
		processor: processor.NewMatchProcessor(marshaler, unmarshaler, engine.NewWhotEngine()),
		machine:   gsm.NewGameStateMachine(),
	}
}

func (m *MatchHandler) GetState() stateless.State {
	return m.machine.MustState()
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	rawLabel, ok := params["label"].(string) // đọc label từ param
	if !ok {
		logger.Error("match init label not found")
		return nil, 0, ""
	}

	matchInfo := &pb1.Match{}
	err := proto.Unmarshal([]byte(rawLabel), matchInfo)
	if err != nil {
		logger.Error("failed to unmarshal match label: %v", err)
		return nil, 0, ""
	}

	matchInfo.Name = entity.ModuleName
	matchInfo.MaxSize = entity.MaxPresences
	matchInfo.MockCodeCard = 0

	matchState := entity.NewMatchState(matchInfo)
	// matchInfo.Size = int32(matchState.Presences.Size())

	logger.Info("match init: %+v", matchInfo)

	matchJson, err := protojson.Marshal(matchInfo)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, constant.TickRate, ""
	}

	logger.Info("match init label= %s", string(matchJson))

	// fire idle event
	procPkg := packager.NewProcessorPackage(&matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(packager.GetContextWithProcessorPackager(procPkg))

	return &matchState, constant.TickRate, string(matchJson)
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
		// fire finish event

		return nil
	}

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate state= %v , shutdown in graceSeconds = %v", state, graceSeconds)
	m.processor.ProcessMatchTerminate(ctx, logger, nk, db, dispatcher, state.(*entity.MatchState))
	return state
}
