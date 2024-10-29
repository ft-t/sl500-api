package sl500_api_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ft-t/sl500-api"
)

func TestCanReadCard(t *testing.T) {
	reader, err := sl500_api.NewConnection("COM3", sl500_api.Baud.Baud19200, true, 3*time.Second)
	assert.NoError(t, err)

	_, err = reader.RfAntennaSta(sl500_api.AntennaOn)
	assert.NoError(t, err)

	_, err = reader.RfLight(sl500_api.ColorYellow)
	assert.NoError(t, err)

	resp, err := reader.RfRequest(sl500_api.RequestAll)
	assert.NoError(t, err)

	if len(resp) != 2 {
		t.Fatalf("Wrong length for rf request %v", resp)
	}

	cardId, _ := reader.RfAnticoll()
	fmt.Println(cardId)

	cardCapacity, _ := reader.RfSelect(cardId)

	if len(cardCapacity) == 0 {
		t.Errorf("Capacity length is wrong %v", cardCapacity)
	}

	_, err = reader.RfM1Authentication2(sl500_api.AuthModeKeyA, 0, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	assert.NoError(t, err)

	blockData, err := reader.RfM1Read(0)
	assert.NoError(t, err)

	if len(blockData) != 16 {
		t.Fatalf("Wrong block data length %v", blockData)
	}

	if !reflect.DeepEqual(blockData[:4], cardId) {
		t.Fatalf("CardId missmatch")
	}

	_, err = reader.RfLight(sl500_api.ColorGreen)
	assert.NoError(t, err)
}
