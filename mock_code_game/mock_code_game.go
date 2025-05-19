package mockcodegame

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
)

var MapMockCodeListCard = make(map[int][]*pb.ListCard)

func InitMapMockCodeListCard() {
	// path := "./mock_in_game"
	// files, err := ioutil.ReadDir(path)
	// if err != nil {
	// 	return
	// }
	listUrlMock := []string{
		"http://103.226.250.195:9000/chinese-poker-mock/1.json",
		"http://103.226.250.195:9000/chinese-poker-mock/2.json",
		"http://103.226.250.195:9000/chinese-poker-mock/3.json",
		"http://103.226.250.195:9000/chinese-poker-mock/4.json",
	}
	for _, urlStr := range listUrlMock {
		// if f.IsDir() {
		// 	continue
		// }
		// nameFile := f.Name()
		// if !strings.HasSuffix(nameFile, ".json") {
		// 	continue
		// }
		// fileMock := filepath.Join(path, f.Name())
		// data, err := os.ReadFile(fileMock) // just pass the file name
		data, err := downloadFile(urlStr)
		if err != nil {
			return
		}
		cpMock := &entity.ChinsePokerMock{}
		err = json.Unmarshal(data, &cpMock)
		if err != nil {
			return
		}
		for _, u := range cpMock.Input.Cards {
			listCard := &pb.ListCard{
				Cards: entity.ParseListCard(u.Card),
			}
			MapMockCodeListCard[cpMock.Id] = append(MapMockCodeListCard[cpMock.Id], listCard)
		}
	}
	fmt.Printf("init map mock code card, len = %d \r\n", len(MapMockCodeListCard))
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
