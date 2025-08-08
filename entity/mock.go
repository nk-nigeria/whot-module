package entity

import pb "github.com/nk-nigeria/cgp-common/proto"

type MockCard struct {
	UserId string `json:"userId"`
	Card   string `json:"card"`
}

type InputChinsePokerMock struct {
	Cards []MockCard `json:"cards"`
}

type ChinsePokerMock struct {
	Id     int                  `json:"id"`
	Name   string               `json:"name"`
	Input  InputChinsePokerMock `json:"input"`
	Output pb.UpdateFinish      `json:"output"`
}
