package dynamic

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	protocol "dynamic-protobuf/protocol"
)

var (
	// protocol
	protocolID = protocol.OSS_LOG_TYPE_OSS_LOG_PLAYER_LOGIN
	msg        = &protocol.OSS_LOG_DATA{
		PlayerLogin: &protocol.PlayerLogin{
			Player: &protocol.Player{
				OpenId:   1,
				Nickname: "nickname",
				Sex:      1,
			},
			Time: &protocol.Time{
				EventTime: time.Now().String()[:10],
			},
			PlatId: 12,
			ZoneId: 32,
		},
	}
)

func TestDynamicParser(t *testing.T) {
	parser := NewDynamicParser()
	parser.SetImportPath("/Users/qinhan/go/src/test/proto")
	if err := parser.ParseFiles("common.proto", "sausageshoot.proto", "protocol.proto"); err != nil {
		t.Error(err)
		return
	}

	// serialize protobuf data
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Error(err)
		return
	}

	tableName, ret, err := parser.ParseToMap("OSS_LOG_DATA", int32(protocolID), data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(tableName)
	t.Log(ret)
}
