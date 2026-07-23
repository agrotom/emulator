package wialon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/agrotom/emulator/internal/codec"
	"github.com/agrotom/emulator/internal/qos"
)

type TCPClient struct {
	conn   net.Conn
	reader codec.PacketReader

	host, port string

	isClosed       bool
	isDisconnected bool

	QoS qos.QoSData
}

func CreateWialonClient(reader codec.PacketReader, host, port string) (*TCPClient, error) {
	if reader == nil {
		return nil, ErrNoReader
	}

	return &TCPClient{
		reader:   reader,
		host:     host,
		port:     port,
		isClosed: true,
	}, nil
}

func (wc *TCPClient) Close() error {
	if wc.conn == nil {
		return nil
	}

	wc.isClosed = true

	return wc.conn.Close()
}

func (wc *TCPClient) Dial(ctx context.Context) error {
	d := net.Dialer{}

	conn, err := d.DialContext(ctx, "tcp", net.JoinHostPort(wc.host, wc.port))

	if err != nil {
		return fmt.Errorf("%w: %w", ErrConnection, err)
	}

	wc.conn = conn
	wc.reader.SetReader(conn)
	wc.isClosed = false

	return nil
}

func (wc *TCPClient) write(ctx context.Context, data []byte) error {
	var err error

	if deadline, ok := ctx.Deadline(); ok {

		if err = wc.conn.SetWriteDeadline(deadline); err != nil {
			return fmt.Errorf("%w (writing): %w", ErrSetDeadline, err)
		}

		defer wc.conn.SetWriteDeadline(time.Time{})
	}

	_, err = wc.conn.Write(data)

	return err
}

func (wc *TCPClient) read(ctx context.Context) ([]byte, error) {
	var err error

	if deadline, ok := ctx.Deadline(); ok {

		if err = wc.conn.SetReadDeadline(deadline); err != nil {
			return nil, fmt.Errorf("%w (reading): %w", ErrSetDeadline, err)
		}

		defer wc.conn.SetReadDeadline(time.Time{})
	}

	resp, err := wc.reader.Read()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (wc *TCPClient) SendPacket(ctx context.Context, packet codec.Packet) ([]byte, error) {
	var (
		err           error
		encodedPacket []byte
		result        []byte
	)

	if wc.conn == nil {
		return nil, ErrDisconnected
	}

	if encodedPacket, err = packet.Encode(); err != nil {
		return nil, fmt.Errorf("encoding packet error: %w", err)
	}

	if err = wc.write(ctx, encodedPacket); err != nil {
		return nil, fmt.Errorf("sending error: %w", err)
	}

	if result, err = wc.read(ctx); err != nil {
		var netErr net.Error

		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("reading error: %w", ErrTimeout)
		}

		return nil, fmt.Errorf("reading error: %w", err)
	}

	return []byte(result), nil
}
