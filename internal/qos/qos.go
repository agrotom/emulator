package qos

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

type QoSSample struct {
	Timestamp      time.Time
	Jitter         time.Duration
	Delay          time.Duration
	LossCount      int
	ConnBreakCount int
}

type QoSData struct {
	Samples []QoSSample
}

func (qos *QoSData) AddData(timestamp time.Time, jitter, delay time.Duration, lossCount, connBreakCount int) {
	qos.Samples = append(qos.Samples, QoSSample{timestamp, jitter, delay, lossCount, connBreakCount})
}

func (qos *QoSData) WriteAll(w io.Writer) error {
	var (
		err error
	)

	writer := csv.NewWriter(w)

	for _, sample := range qos.Samples {
		err = writer.Write([]string{strconv.FormatInt(sample.Timestamp.UnixNano(), 10),
			strconv.FormatInt(int64(sample.Jitter), 10),
			strconv.FormatInt(int64(sample.Delay), 10),
			strconv.FormatInt(int64(sample.LossCount), 10),
			strconv.FormatInt(int64(sample.ConnBreakCount), 10)})

		if err != nil {
			return err
		}
	}

	writer.Flush()

	return nil
}
