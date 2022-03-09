package visa

import (
	"context"
	"net/http"
	"os"

	"github.com/grokify/goauth"
	"github.com/grokify/mogo/crypto/tlsutil"
)

var (
	VisaAppKeyFileEnv     = "VISA_APP_KEY_FILE"
	VisaAppCertFileEnv    = "VISA_APP_CERT_FILE"
	VisaAppUserIdEnv      = "VISA_APP_USERID"
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

func NewClient(cfg Config) (*http.Client, error) {
	tlsConfig := tlsutil.NewTLSConfig()

	if err := tlsConfig.LoadX509KeyPair(cfg.AppCertFile, cfg.AppKeyFile); err != nil {
		return nil, err
	}

	if err := tlsConfig.LoadCACert(cfg.VDPCACertFile); err != nil {
		return nil, err
	}

	if err := tlsConfig.LoadCACert(cfg.CACertFile); err != nil {
		return nil, err
	}

	// tlsConfig.Inflate()

	if token, err := goauth.BasicAuthToken(cfg.Username, cfg.Password); err != nil {
		return nil, err
	} else {
		return goauth.NewClientTLSToken(
			context.Background(), tlsConfig.Config, token), nil
	}
}

func ConfigFromEnv() Config {
	return Config{
		AppKeyFile:    os.Getenv(VisaAppKeyFileEnv),
		AppCertFile:   os.Getenv(VisaAppCertFileEnv),
		VDPCACertFile: os.Getenv(VDPCACertFileEnv),
		CACertFile:    os.Getenv(GeoTrustCACertFileEnv),
		Username:      os.Getenv(VisaAppUserIdEnv),
		Password:      os.Getenv(VisaAppPasswordEnv),
	}
}

func NewClientEnv() (*http.Client, error) {
	return NewClient(ConfigFromEnv())
}
