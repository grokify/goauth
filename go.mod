module github.com/grokify/goauth

go 1.21

// toolchain go1.21.1 // incompatible with older Go versions.

require (
	github.com/SparkPost/gosparkpost v0.2.0
	github.com/caarlos0/env/v9 v9.0.0
	github.com/golang-jwt/jwt/v5 v5.1.0
	github.com/google/go-querystring v1.1.0
	github.com/grokify/go-salesforce v0.2.23
	github.com/grokify/gocharts/v2 v2.15.2
	github.com/grokify/mogo v0.61.0
	github.com/grokify/xgo v0.6.0
	github.com/jessevdk/go-flags v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.31.0
	golang.org/x/oauth2 v0.14.0
	google.golang.org/api v0.152.0
	gopkg.in/jeevatkm/go-model.v1 v1.1.0
)

require (
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/martinlindhe/base36 v1.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/quicktemplate v1.7.0 // indirect
	github.com/xuri/efp v0.0.0-20231025114914-d1ff6096ae53 // indirect
	github.com/xuri/excelize/v2 v2.8.0 // indirect
	github.com/xuri/nfp v0.0.0-20230919160717-d98342af3f05 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231120223509-83a465c0220f // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

// replace github.com/grokify/mogo => ../mogo

// replace github.com/grokify/gocharts/v2 => ../gocharts
