package channel

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/agrotom/emulator/internal/codec"
	"github.com/agrotom/emulator/internal/qos"
	"github.com/agrotom/emulator/internal/wialon"
)

type ControlChannel struct {
	wialon.Client

	curLosses int
	curBreaks int

	opts ControlChannelOptions

	qos qos.QoSData
}

func CreateControlChannel(tcp wialon.Client, opts ControlChannelOptions) *ControlChannel {
	cch := &ControlChannel{
		Client:    tcp,
		opts:      opts,
		curLosses: 0,
		curBreaks: 0,
		qos:       qos.QoSData{},
	}

	return cch
}

func (cch *ControlChannel) GetQoSData() *qos.QoSData {
	return &cch.qos
}

func (cch *ControlChannel) SendPacket(ctx context.Context, packet codec.Packet) ([]byte, error) {
	var jitter int64

	nowTime := time.Now()

	if cch.opts.Jitter > 0 {
		jitter = rand.Int64N(2*int64(cch.opts.Jitter)) - int64(cch.opts.Jitter)
	}

	delay := cch.opts.Delay + time.Duration(jitter)

	defer cch.qos.AddData(nowTime, time.Duration(jitter)*time.Millisecond, delay*time.Millisecond, cch.curLosses, cch.curBreaks)

	if delay > 0 {
		//log.Printf("channel: delay=%d and jitter=%d", cch.opts.Delay, jitter)
		time.Sleep(delay)
	}

	if cch.curLosses < cch.opts.MaxLossCount && rand.Float32() < cch.opts.LossPercent {
		//log.Printf("channel: loosing packet (%d of %d)", cch.curLosses, cch.opts.MaxLossCount)
		cch.curLosses++

		return nil, wialon.ErrTimeout
	}

	if cch.curBreaks < cch.opts.MaxConnBreakCount && rand.Float32() < cch.opts.ConnBreakPercent {
		//log.Printf("channel: closing channel (%d of %d)", cch.curBreaks, cch.opts.MaxConnBreakCount)
		cch.curBreaks++
		cch.Close()

		return nil, wialon.ErrDisconnected
	}

	return cch.Client.SendPacket(ctx, packet)
}
