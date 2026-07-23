package test

/*import (
	"testing"

	"github.com/agrotom/emulator/internal/codec/egts/auth"
	"github.com/agrotom/emulator/internal/codec/egts/builder"
	"github.com/agrotom/emulator/internal/codec/egts/common"
	"github.com/agrotom/emulator/internal/wialon"
)

func createClientAndDial() (*wialon.EgtsClient, error) {
	client := wialon.CreateEgtsClient(nil)

	return client, client.Dial(Host, Port)
}

func WialonClientDial(t *testing.T) {
	_, err := createClientAndDial()

	if err != nil {
		t.Errorf("error while dialing wialon client: %s", err.Error())
		return
	}
}

func WialonClientAuthorize(t *testing.T) {
	var err error

	client, err := createClientAndDial()

	if err != nil {
		t.Errorf("error while dialing to wialon server: %s", err.Error())
		return
	}

	builder := builder.NewFrameBuilder()

	builder.AddAuthRecord(common.EgtsSrTermIdentity, &auth.AuthService{
		TerminalIdentifier:                         0,
		InternationalMobileEquipmentIdentityExists: true,
		InternationalMobileEquipmentIdentity:       UnitID,
	})

	packet := builder.Build()

	if err = client.SendPacket(packet); err != nil {
		t.Errorf("error while sending packet to wialon server: %s", err.Error())
		return
	}

	client.Recv()

}
*/
