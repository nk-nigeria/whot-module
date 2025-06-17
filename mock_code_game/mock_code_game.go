package mockcodegame

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/heroiclabs/nakama-common/runtime"

	pb "github.com/nk-nigeria/cgp-common/proto/whot"
	"github.com/nk-nigeria/whot-module/entity"
)

var MapMockCodeListCard = make(map[int][]*pb.ListCard)

func InitMapMockCodeListCard(logger runtime.Logger) {
	listUrlMock := []string{
		"https://raw.githubusercontent.com/huy24112001/whot-json/refs/heads/main/1.json",
		"https://raw.githubusercontent.com/huy24112001/whot-json/refs/heads/main/2.json",
		"https://raw.githubusercontent.com/huy24112001/whot-json/refs/heads/main/3.json",
	}
	for _, urlStr := range listUrlMock {
		data, err := downloadFile(urlStr)
		if err != nil {
			logger.Error("Failed to download: %s\n", urlStr)
			continue
		}
		cpMock := &entity.ChinsePokerMock{}
		err = json.Unmarshal(data, &cpMock)
		if err != nil {
			logger.Error("Failed to parse: %s\n", urlStr)
			continue
		}
		for _, u := range cpMock.Input.Cards {
			listCard := &pb.ListCard{
				Cards: entity.ParseListCard(u.Card),
			}
			MapMockCodeListCard[cpMock.Id] = append(MapMockCodeListCard[cpMock.Id], listCard)
		}
	}
	logger.Info("init map mock code card, len = %d\n", len(MapMockCodeListCard))
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 {
		return nil, errors.New("status not ok")
	}
	if resp.Body == nil {
		return nil, errors.New("body is nil")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
