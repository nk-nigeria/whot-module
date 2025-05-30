package mock

import (
	"testing"

	pb "github.com/nakamaFramework/cgp-common/proto/whot"
	"github.com/nakamaFramework/whot-module/entity"
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
