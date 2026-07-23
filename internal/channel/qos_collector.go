package channel

import (
	"context"
	"errors"
	"time"

	"github.com/agrotom/emulator/internal/codec"
	"github.com/agrotom/emulator/internal/qos"
	"github.com/agrotom/emulator/internal/wialon"
)

type QoSCollector struct {
	wialon.Client

	startTime time.Time
	rtt       time.Duration

	lastRTT time.Duration
	avgRTT  time.Duration
	jitter  time.Duration

	curLosses int
	curBreaks int

	qos qos.QoSData
}

func CreateQoSCollector(tcp wialon.Client) *QoSCollector {
	return &QoSCollector{
		tcp,
		time.Time{},
		0,
		0,
		0,
		0,
		0,
		0,
		qos.QoSData{},
	}
}

func (qosc *QoSCollector) GetQoSData() qos.QoSData {
	return qosc.qos
}

func (qosc *QoSCollector) Update(rtt time.Duration) {
	if qosc.avgRTT == 0 {
		qosc.avgRTT = rtt
	} else {
		qosc.avgRTT = (qosc.avgRTT*15 + rtt) / 16
	}

	if qosc.lastRTT != 0 {
		diff := rtt - qosc.lastRTT
		if diff < 0 {
			diff = -diff
		}
		qosc.jitter = (qosc.jitter*15 + diff) / 16
	}

	qosc.lastRTT = rtt
}

func (qosc *QoSCollector) SendPacket(ctx context.Context, packet codec.Packet) ([]byte, error) {
	startTime := time.Now()
	result, err := qosc.Client.SendPacket(ctx, packet)

	if errors.Is(err, wialon.ErrDisconnected) {
		qosc.curBreaks++
		return result, err
	}

	if errors.Is(err, wialon.ErrTimeout) {
		qosc.curLosses++
		return result, err
	}

	rtt := time.Since(startTime)

	qosc.Update(rtt)

	qosc.qos.AddData(startTime, qosc.jitter/time.Millisecond, rtt/time.Millisecond, qosc.curLosses, qosc.curBreaks)

	return result, err
}
