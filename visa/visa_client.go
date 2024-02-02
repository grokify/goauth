package visa

import (
	"context"
	"net/http"
	"os"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/crypto/tlsutil"
)

var (
	VisaAppKeyFileEnv     = "VISA_APP_KEY_FILE"
	VisaAppCertFileEnv    = "VISA_APP_CERT_FILE"
	VisaAppUserIDEnv      = "VISA_APP_USERID"
	VisaAppPasswordEnv    = "VISA_APP_PASSWORD" // #nosec G101
	VDPCACertFileEnv      = "VISA_VDP_CA_CERT_FILE"
	GeoTrustCACertFileEnv = "VISA_GEOTRUST_CA_CERT_FILE"
)

type Config struct {
	AppKeyFile    string
	AppCertFile   string
	VDPCACertFile string
	CACertFile    string
	Username      string
	Password      string
}

func (cfg Config) NewClient() (*http.Client, error) {
	tlsConfig, err := tlsutil.NewTLSConfig("", "", []string{}, []string{}, false)
	if err != nil {
		return nil, err
	}

	if len(cfg.AppCertFile) > 0 || len(cfg.AppKeyFile) > 0 {
		if err := tlsConfig.LoadX509KeyPair(cfg.AppCertFile, cfg.AppKeyFile); err != nil {
			return nil, err
		}
	}

	if len(cfg.VDPCACertFile) > 0 {
		if err := tlsConfig.LoadRootCACert(cfg.VDPCACertFile); err != nil {
			return nil, err
		}
	}

	if len(cfg.CACertFile) > 0 {
		if err := tlsConfig.LoadRootCACert(cfg.CACertFile); err != nil {
			return nil, err
		}
	}

	// tlsConfig.Inflate()

	if token, err := authutil.BasicAuthToken(cfg.Username, cfg.Password); err != nil {
		return nil, err
	} else {
		return authutil.NewClientTLSToken(
			context.Background(), tlsConfig.Config, token), nil
	}
}

func ConfigFromEnv() Config {
	return Config{
		AppKeyFile:    os.Getenv(VisaAppKeyFileEnv),
		AppCertFile:   os.Getenv(VisaAppCertFileEnv),
		VDPCACertFile: os.Getenv(VDPCACertFileEnv),
		CACertFile:    os.Getenv(GeoTrustCACertFileEnv),
		Username:      os.Getenv(VisaAppUserIDEnv),
		Password:      os.Getenv(VisaAppPasswordEnv),
	}
}

func NewClientEnv() (*http.Client, error) {
	return ConfigFromEnv().NewClient()
}
