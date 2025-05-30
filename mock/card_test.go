package mock

import (
	"testing"

	pb "github.com/nakama-nigeria/cgp-common/proto/whot"
	"github.com/nakama-nigeria/whot-module/entity"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestParseMockCard(t *testing.T) {
	fileMock := "./mock_card/natural_special.txt"
	list := entity.ParseMockCard(fileMock)
	x := pb.UpdateFinish{}
	pp := &protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	data, _ := pp.Marshal(&x)
	t.Logf("%s", string(data))
	t.Logf("%v", list)

}
