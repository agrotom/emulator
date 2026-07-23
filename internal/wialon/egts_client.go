package wialon

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/agrotom/emulator/internal/codec/egts"
)

type OutgoingPacket struct {
	pid      uint16
	raw      []byte
	sentTime time.Time
	retires  int
}

type IngoingPacket struct {
	packet   *egts.EGTSFrame
	sentTime time.Time
	wait     chan *egts.EGTSFrame
}

type EgtsClient struct {
	conn net.Conn
	r    *egts.FrameReader

	pid         uint16
	expectedPid uint16

	outgoing  map[uint16]*OutgoingPacket
	recieved  map[uint16]*IngoingPacket
	muPending sync.RWMutex
	muClose   sync.RWMutex

	closed bool
}

type EgtsClientOptions struct {
	InitialPacketIdentifier uint16
}

func CreateEgtsClient(opts *EgtsClientOptions) *EgtsClient {
	if opts != nil {
		return &EgtsClient{
			conn:        nil,
			r:           nil,
			pid:         opts.InitialPacketIdentifier,
			expectedPid: opts.InitialPacketIdentifier + 1,
			outgoing:    make(map[uint16]*OutgoingPacket),
			recieved:    make(map[uint16]*IngoingPacket),
			muPending:   sync.RWMutex{},
			closed:      true,
		}
	}

	return &EgtsClient{
		conn:        nil,
		r:           nil,
		pid:         0,
		expectedPid: 1,
		outgoing:    make(map[uint16]*OutgoingPacket),
		recieved:    make(map[uint16]*IngoingPacket),
		muPending:   sync.RWMutex{},
		closed:      true,
	}
}

func (wc *EgtsClient) Dial(host string, port string) error {
	var err error

	wc.conn, err = net.Dial("tcp", net.JoinHostPort(host, port))

	if err != nil {
		return fmt.Errorf("error while dialing wialon client: %w", err)
	}

	wc.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	wc.r = egts.NewReader(wc.conn)

	wc.muClose.Lock()
	wc.closed = false
	wc.muClose.Unlock()

	go wc.readerLoop()
	//go wc.retransmissionLoop()

	return nil
}

func (wc *EgtsClient) SendPacket(frame *egts.EGTSFrame) error {
	wc.muClose.RLock()
	defer wc.muClose.RUnlock()

	if wc.closed {
		return fmt.Errorf("wialon client is closed: unable to sending packet")
	}

	data, err := frame.Encode()

	if err != nil {
		return fmt.Errorf("error while encoding egts packet: %w", err)
	}

	log.Printf(
		"TX %d bytes: %x",
		len(data),
		data,
	)

	_, err = wc.conn.Write(data)

	if err != nil {
		return fmt.Errorf("error while sending egts packet: %w", err)
	}

	wc.muPending.Lock()

	wc.outgoing[wc.pid] = &OutgoingPacket{
		wc.pid,
		data,
		time.Now(),
		0,
	}

	wc.recieved[wc.expectedPid] = &IngoingPacket{
		packet:   frame,
		sentTime: time.Now(),
	}

	wc.muPending.Unlock()

	return nil
}

func (wc *EgtsClient) retransmissionLoop() {
	for {
		wc.muClose.RLock()

		if wc.closed {
			log.Printf("wialon client is closed, stopping the retransmission loop...")
			wc.muClose.RUnlock()
			return
		}

		wc.muClose.RUnlock()
	}
}

func (wc *EgtsClient) readerLoop() {
	for {
		wc.muClose.RLock()

		if wc.closed {
			log.Printf("wialon client is closed, stopping the reader loop...")
			wc.muClose.RUnlock()
			return
		}

		packet, err := wc.r.ReadFrame()

		if err != nil {
			log.Printf("error while reading packet: %s", err.Error())
			wc.muClose.RUnlock()
			continue
		}

		wc.muPending.RLock()
		_, ok := wc.outgoing[packet.PacketIdentifier]
		wc.muPending.RUnlock()

		if !ok {
			log.Printf("unknown packet with pid: %d", packet.PacketIdentifier)
		} else {
			wc.muPending.Lock()
			wc.recieved[wc.expectedPid].wait <- packet
			delete(wc.outgoing, packet.PacketIdentifier)
			wc.muPending.Unlock()
		}

		log.Printf("packet recieved:\t%+v", packet)
		wc.muClose.RUnlock()
	}
}

func (wc *EgtsClient) Recv() *egts.EGTSFrame {
	return <-wc.recieved[wc.expectedPid].wait
}
