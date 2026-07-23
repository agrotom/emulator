package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/agrotom/emulator/internal/channel"
	"github.com/agrotom/emulator/internal/codec"
	"github.com/agrotom/emulator/internal/codec/wialonips"
	"github.com/agrotom/emulator/internal/config"
	"github.com/agrotom/emulator/internal/mathutil"
	"github.com/agrotom/emulator/internal/qos"
	"github.com/agrotom/emulator/internal/session"
	"github.com/agrotom/emulator/internal/wialon"
)

const (
	DefaultConfigPath string = "./config.json"
	DefaultQoSPath    string = "./qos/"
)

func LoadConfig(path string) config.Config {

	cfg, err := config.Load(path)

	if err != nil {
		log.Fatalf("unable to load config file: %s", err.Error())
	}

	log.Printf("Config %s is successfuly loaded", path)

	return cfg
}

func SaveQoSFile(qos qos.QoSData, qosPath, name string) error {
	var (
		qosFile *os.File
		err     error
	)

	if qosFile, err = CreateQoSFile(qosPath, name); err != nil {
		return fmt.Errorf("unable to create qos file: %w", err)
	}

	err = qos.WriteAll(qosFile)

	if err != nil {
		return fmt.Errorf("unable to save qos file: %w", err)
	}

	if err = qosFile.Close(); err != nil {
		return fmt.Errorf("unable to close file: %w", err)
	}

	return nil
}

func CreateSession(qosPath string, cfg config.SimulationConfig) session.Session {
	var (
		reader codec.PacketReader
		tcp    wialon.Client
		pb     codec.PacketBuilder
		ssn    session.Session
		err    error
	)

	ssnOpts := session.SessionOptions{
		UnitID:            cfg.Creds.IMEI,
		Password:          cfg.Creds.Password,
		MaxTimeout:        cfg.NetCfg.MaxTimeout,
		MaxReconnectTries: cfg.NetCfg.MaxReconnectTries,
		MaxResendTries:    cfg.NetCfg.MaxResendTries,
		ResendWaitTime:    cfg.NetCfg.ResendWaitTime,
	}

	switch cfg.Protocol {
	case config.WialonIPS:
		reader = wialonips.CreateWialonReader()
	default:
		log.Fatalf("unknown protocol: %s", cfg.Protocol)
	}

	tcp, err = wialon.CreateWialonClient(reader, cfg.Creds.Host, cfg.Creds.Port)

	qosCollector := channel.CreateQoSCollector(tcp)
	tcp = qosCollector

	defer func() {
		if ssn != nil {
			log.Printf("%s: setting qos data file path...", ssnOpts.UnitID)
			ssn.SetOnSuccess(func() {
				if err := SaveQoSFile(qosCollector.GetQoSData(), qosPath, ssnOpts.UnitID); err != nil {
					log.Printf("%s: unable to save qos file: %s", ssnOpts.UnitID, err.Error())
				} else {
					log.Printf("%s: qos file is saved", ssnOpts.UnitID)
				}
			})
		}
	}()

	if cfg.HasControlChannel {
		chanOpts := channel.ControlChannelOptions{
			Jitter:            cfg.ChannelCfg.Jitter,
			Delay:             cfg.ChannelCfg.Delay,
			LossPercent:       cfg.ChannelCfg.LossPercent,
			ConnBreakPercent:  cfg.ChannelCfg.ConnBreakPercent,
			MaxLossCount:      cfg.ChannelCfg.MaxLossCount,
			MaxConnBreakCount: cfg.ChannelCfg.MaxConnBreakCount,
		}

		cchan := channel.CreateControlChannel(tcp, chanOpts)
		tcp = cchan

		defer func() {
			if ssn != nil {
				log.Printf("%s: setting chan qos data file path...", ssnOpts.UnitID)
				ssn.SetOnSuccess(func() {
					if err := SaveQoSFile(*cchan.GetQoSData(), qosPath, ssnOpts.UnitID+"_chan"); err != nil {
						log.Printf("%s: unable to save chan qos file: %s", ssnOpts.UnitID, err.Error())
					} else {
						log.Printf("%s: chan qos file is saved", ssnOpts.UnitID)
					}
				})
			}
		}()
	}

	switch cfg.Protocol {
	case config.WialonIPS:
		pb = &wialonips.WialonPacketBuilder{}
	default:
		log.Fatalf("unknown protocol: %s", cfg.Protocol)
	}

	ssn, err = session.CreateWialonSession(ssnOpts, tcp, pb)

	if err != nil {
		log.Fatalf("unable to create session: %s", err.Error())
	}

	return ssn
}

func CreateQoSFile(path, fileName string) (*os.File, error) {
	var (
		file *os.File
		err  error
	)

	if file, err = os.OpenFile(path+fileName+".csv", os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return nil, err
	}

	return file, nil
}

func InitQoS(path string) error {
	return os.MkdirAll(path, 0644)
}

func main() {
	var (
		wg sync.WaitGroup
	)

	cfgPath := flag.String("cfg", DefaultConfigPath, "path to config file")
	qosPath := flag.String("qos", DefaultQoSPath, "path to qos folder")

	flag.Parse()

	cfg := LoadConfig(*cfgPath)
	InitQoS(*qosPath)

	for _, sopts := range cfg.Simulations {
		wg.Add(1)

		log.Printf("creating session #%s", sopts.Creds.IMEI)
		ssn := CreateSession(*qosPath, sopts)

		logger := log.Default()

		if len(sopts.BoundingBoxes) < 1 {
			log.Fatalf("session #%s has not bounding boxes. please provide!", sopts.Creds.IMEI)
		}

		go func() {
			defer wg.Done()

			var err error

			ctx := context.Background()
			if err = ssn.Login(ctx); err != nil {
				logger.Fatalf("unable to login to session: %s", err.Error())
			}

			var startPos *mathutil.Vector2f = nil

			if !sopts.AutoStartPos {
				startPos = &sopts.StartPos
			}

			var endPos *mathutil.Vector2f = nil

			if !sopts.AutoStartPos {
				endPos = &sopts.EndPos
			}

			opts := session.SimulationOptions{
				StartPos: startPos,
				EndPos:   endPos,

				MaxPacketToSend: -1,
				Sattelites:      sopts.Sattelites,

				StepDistance: sopts.StepDistanceM,
				StepMillis:   sopts.StepMillis,

				MinDist: sopts.MinDist,
				MaxDist: sopts.MaxDist,
				Bounds:  sopts.BoundingBoxes,
			}

			if err = ssn.StartSimulation(ctx, opts); err != nil {
				logger.Fatalf("unable to start simulation: %s", err.Error())
			}

			logger.Printf("simulation of session #%s is over successfuly", sopts.Creds.IMEI)
		}()
	}

	wg.Wait()
}
