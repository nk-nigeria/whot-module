package mock

import (
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/nakamaFramework/whot-module/entity"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestParseMockCard(t *testing.T) {
	fileMock := "/home/sondq/Documents/myspace/cgb-chinese-poker-module/mock/mock_card/natural_special.txt"
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
