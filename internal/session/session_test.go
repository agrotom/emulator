package session

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/agrotom/emulator/internal/channel"
	"github.com/agrotom/emulator/internal/codec/wialonips"
	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/wialon"
)

var (
	UnitID   string
	UnitID2  string
	Host     string
	Port     string
	Password string
)

var StartCoords = mathutil.Vector2f{X: 51.546529, Y: 46.037097}
var EndCoords = mathutil.Vector2f{X: 51.536257, Y: 46.02405}

const StepM int = 100
const StepMS time.Duration = 500

var SimOptions = SimulationOptions{
	StartPos:        &StartCoords,
	EndPos:          &EndCoords,
	MaxPacketToSend: 20,
	Sattelites:      1,
	StepDistance:    StepM,
	StepMillis:      StepMS,
}

func TestMain(m *testing.M) {
	UnitID = os.Getenv("EMULATOR_TEST_UNIT_ID")
	UnitID2 = os.Getenv("EMULATOR_TEST_UNIT_ID_2")
	Host = os.Getenv("EMULATOR_TEST_HOST")
	Port = os.Getenv("EMULATOR_TEST_PORT")
	Password = os.Getenv("EMULATOR_TEST_PASSWORD")

	os.Exit(m.Run())
}

func createSession(unitID string) (Session, error) {
	tcp, err := wialon.CreateWialonClient(&wialonips.WialonIPSReader{}, Host, Port)

	if err != nil {
		return nil, err
	}

	session, err := CreateWialonSession(DefaultOptions(unitID, Password), tcp, &wialonips.WialonPacketBuilder{})

	return session, err
}

func createBrokenSession(unitID string, opts channel.ControlChannelOptions) (Session, error) {
	tcp, err := wialon.CreateWialonClient(&wialonips.WialonIPSReader{}, Host, Port)

	if err != nil {
		return nil, err
	}

	cch := channel.CreateControlChannel(tcp, opts)
	session, err := CreateWialonSession(DefaultOptions(unitID, Password), cch, &wialonips.WialonPacketBuilder{})

	return session, err
}

func TestCreateSingleSession(t *testing.T) {
	_, err := createSession(UnitID)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}
}

func TestLoginSingleSession(t *testing.T) {
	session, err := createSession(UnitID)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}
}

func TestSimulationSingleSession(t *testing.T) {
	session, err := createSession(UnitID)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}

	err = session.StartSimulation(ctx, SimOptions)

	if err != nil {
		t.Errorf("error while sending pos to a session: %s", err.Error())
		return
	}
}

func TestSimulationSingleDelayedSession(t *testing.T) {
	opts := channel.DefaultOptions()

	opts.Delay = time.Millisecond

	session, err := createBrokenSession(UnitID, opts)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}

	err = session.StartSimulation(ctx, SimOptions)

	if err != nil {
		t.Errorf("error while sending pos to a session: %s", err.Error())
		return
	}
}

func TestSimulationSingleJitterSession(t *testing.T) {
	opts := channel.DefaultOptions()

	opts.Delay = time.Millisecond
	opts.Jitter = time.Millisecond / 2

	session, err := createBrokenSession(UnitID, opts)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}

	err = session.StartSimulation(ctx, SimOptions)

	if err != nil {
		t.Errorf("error while sending pos to a session: %s", err.Error())
		return
	}
}

func TestSimulationSingleLoss30Session(t *testing.T) {
	opts := channel.DefaultOptions()

	opts.Delay = time.Millisecond * 25
	opts.Jitter = time.Millisecond * 50
	opts.LossPercent = 0.3
	opts.MaxLossCount = 4

	session, err := createBrokenSession(UnitID, opts)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}

	err = session.StartSimulation(ctx, SimOptions)

	if err != nil {
		t.Errorf("error while sending pos to a session: %s", err.Error())
		return
	}
}

func TestSimulationSingleBreak30Session(t *testing.T) {
	opts := channel.DefaultOptions()

	opts.Delay = time.Millisecond * 25
	opts.Jitter = time.Millisecond * 50
	opts.ConnBreakPercent = 0.30
	opts.MaxConnBreakCount = 3

	session, err := createBrokenSession(UnitID, opts)

	if err != nil {
		t.Errorf("error while creating a session: %s", err.Error())
		return
	}

	ctx := context.Background()

	err = session.Login(ctx)

	if err != nil {
		t.Errorf("error while login to a session: %s", err.Error())
		return
	}

	err = session.StartSimulation(ctx, SimOptions)

	if err != nil {
		t.Errorf("error while sending pos to a session: %s", err.Error())
		return
	}
}

func runAsyncSession(unitID string, t *testing.T) {
	go func() {
		session, err := createSession(unitID)

		if err != nil {
			t.Errorf("error while creating a session: %s", err.Error())
			return
		}

		ctx := context.Background()

		err = session.Login(ctx)

		if err != nil {
			t.Errorf("error while login to a session: %s", err.Error())
			return
		}

		err = session.StartSimulation(ctx, SimOptions)

		if err != nil {
			t.Errorf("error while sending pos to a session: %s", err.Error())
			return
		}
	}()
}

func TestSimulationMultiSession(t *testing.T) {
	runAsyncSession(UnitID, t)
	runAsyncSession(UnitID2, t)
}
