package session

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/agrotom/emulator/internal/codec"
	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/track"
	"github.com/agrotom/emulator/internal/wialon"
)

type Session interface {
	Login(context.Context) error
	StartSimulation(context.Context, SimulationOptions) error
	SetOnSuccess(func())
}

type wialonSession struct {
	tcp wialon.Client

	pb codec.PacketBuilder

	opts SessionOptions

	successEvent []func()
}

func CreateWialonSession(opts SessionOptions, tcp wialon.Client, pb codec.PacketBuilder) (Session, error) {
	var err error

	if tcp == nil {
		return nil, ErrNoClient
	}

	if pb == nil {
		return nil, ErrNoPacketBuilder
	}

	ws := &wialonSession{
		tcp:  tcp,
		pb:   pb,
		opts: opts,
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*ws.opts.MaxTimeout)
	defer cancel()

	if err = ws.tcp.Dial(timeoutCtx); err != nil {
		return nil, err
	}

	return ws, nil
}

func (ws *wialonSession) SetOnSuccess(f func()) {
	ws.successEvent = append(ws.successEvent, f)
}

func (ws *wialonSession) Login(ctx context.Context) error {
	var (
		err error
	)

	packet := ws.pb.LoginPacket(&codec.LoginOptions{
		IMEI:     ws.opts.UnitID,
		Password: ws.opts.Password,
	})

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*ws.opts.MaxTimeout)
	defer cancel()

	if err = ws.trySendPacket(timeoutCtx, packet); err != nil {
		return err
	}

	return nil
}

func (ws *wialonSession) trySendPacket(ctx context.Context, packet codec.Packet) error {
	var (
		resp []byte
		err  error
	)

	for reconnect := 0; reconnect <= ws.opts.MaxReconnectTries; reconnect++ {
		if reconnect > 0 {
			time.Sleep(time.Millisecond * ws.opts.ResendWaitTime)
			log.Printf("%s: reconnecting to server... (%d of %d)", ws.opts.UnitID, reconnect, ws.opts.MaxReconnectTries)
			if err = ws.tcp.Dial(ctx); err != nil {
				continue
			}
		}

		for retry := 0; retry < ws.opts.MaxResendTries; retry++ {
			if retry > 0 {
				time.Sleep(time.Millisecond * ws.opts.ResendWaitTime)
				log.Printf("%s: resending packet... (%d of %d)", ws.opts.UnitID, retry, ws.opts.MaxResendTries)
			}

			timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*ws.opts.MaxTimeout)

			resp, err = ws.tcp.SendPacket(timeoutCtx, packet)
			cancel()

			if err == nil {
				return packet.ValidateResponse(resp)
			}

			if errors.Is(err, wialon.ErrDisconnected) {
				break
			}

			if errors.Is(err, wialon.ErrTimeout) {
				continue
			}

			return err
		}
	}

	return ErrSendPacket
}

func (ws *wialonSession) StartSimulation(ctx context.Context, opts SimulationOptions) error {
	var (
		err    error
		coords []mathutil.Vector2f
	)

	if opts.StartPos == nil {
		startPos := track.RandomPoint(opts.Bounds.ChooseRandomBounds())

		if startPos, err = track.GetOSRMNearestPoint(startPos); err != nil {
			return err
		}

		log.Printf("%s: start pos is not given. setting it to (%f, %f)", ws.opts.UnitID, startPos.X, startPos.Y)

		opts.StartPos = &startPos
	}

	if opts.EndPos == nil {
		endPos := track.RandomDestination(*opts.StartPos, opts.MinDist, opts.MaxDist)

		if endPos, err = track.GetOSRMNearestPoint(endPos); err != nil {
			return err
		}

		log.Printf("%s: end pos is not given. setting it to (%f, %f)", ws.opts.UnitID, endPos.X, endPos.Y)

		opts.EndPos = &endPos
	}

	if coords, err = track.GetOSRMRoute(*opts.StartPos, *opts.EndPos); err != nil {
		return fmt.Errorf("%s: error while getting osrm route: %w", ws.opts.UnitID, err)
	}

	coords = track.Interpolate(coords, opts.StepDistance)

	records := track.Simulate(coords, opts.StepMillis)

	for i, rd := range records {
		if opts.MaxPacketToSend >= 0 && i >= opts.MaxPacketToSend {
			log.Printf("%s: reached maximum packets to send count (%d)", ws.opts.UnitID, opts.MaxPacketToSend)
			break
		}

		packet := ws.pb.DataPacket(&codec.DataOptions{
			Sattelites: opts.Sattelites,
			NavRecord:  rd,
		})

		if err = ws.trySendPacket(ctx, packet); err != nil {
			return err
		}

		log.Printf("sending packet...: %+v", rd)
	}

	if ws.successEvent != nil {
		for _, sucEv := range ws.successEvent {
			sucEv()
		}
	}

	return nil
}
