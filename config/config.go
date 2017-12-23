package config

import (
	"io/ioutil"
	"net/url"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/keypair"
	"gopkg.in/yaml.v2"
)

type GateConfig struct {
	Port          string `yaml:"port"`
	RawHorizonURL string `yaml:"horizon_url"`
	Seed          string `yaml:"seed"`
	AccountID     string `yaml:"account_id"`
	LogLevelS     string `yaml:"log_level"`

	HorizonURL *url.URL
	KP         *keypair.Full
	LogLevel   logan.Level
}

func InitConfig(filePath string) (*GateConfig, error) {
	rawConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config = new(GateConfig)
	err = yaml.Unmarshal(rawConfig, config)
	if err != nil {
		return nil, err
	}

	config.HorizonURL, err = url.Parse(config.RawHorizonURL)
	if err != nil {
		return nil, err
	}

	err = config.parseKP()
	if err != nil {
		return nil, err
	}

	config.LogLevel = logLevel(config.LogLevelS)
	return config, nil
}

func (gc *GateConfig) parseKP() error {
	kp, err := keypair.Parse(gc.Seed)
	if err != nil {
		return errors.Wrap(err, "failed to parse seed")
	}

	var ok bool
	gc.KP, ok = kp.(*keypair.Full)
	if !ok {
		return errors.New("must be a seed")
	}

	if gc.AccountID == "" {
		gc.AccountID = gc.KP.Address()
	}
	_, err = keypair.Parse(gc.AccountID)
	if err != nil {
		return errors.Wrap(err, "failed to parse accountID")
	}
	return nil
}

func logLevel(ll string) logan.Level {
	switch ll {
	case "debug":
		return logan.DebugLevel
	case "error":
		return logan.ErrorLevel
	case "info":
		return logan.InfoLevel
	case "warn":
		return logan.WarnLevel
	default:
		return logan.WarnLevel
	}
}
