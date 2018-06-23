package sl500_api_test

import (
	"testing"
	"sl500-api"
	"fmt"
	"reflect"
)

func TestCanReadCard(t *testing.T) {
	reader, err := sl500_api.NewConnection("COM3", sl500_api.Baud.Baud19200)
	reader.RfAntennaSta(sl500_api.AntennaOn)
	if err != nil {
		t.Fatal(err)
	}
	reader.RfLight(sl500_api.ColorYellow)
	resp, _ := reader.RfRequest(sl500_api.RequestAll)

	if len(resp) != 2 {
		t.Fatalf("Wrong length for rf request %v", resp)
	}

	cardId, _ := reader.RfAnticoll()
	fmt.Println(cardId)

	cardCapacity, _ := reader.RfSelect(cardId)

	if len(cardCapacity) == 0 {
		t.Errorf("Capacity length is wrong %v", cardCapacity)
	}

	reader.RfM1Authentication2(sl500_api.AuthModeKeyA, 0, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	blockData, err := reader.RfM1Read(0)

	if err != nil {
		t.Fatal(err)
	}

	if len(blockData) != 16 {
		t.Fatalf("Wrong block data length %v", blockData)
	}
	if !reflect.DeepEqual(blockData[:4], cardId) {
		t.Fatalf("CardId missmatch")
	}
	reader.RfLight(sl500_api.ColorGreen)

	if err != nil {
		t.Error(err)
	}

}
