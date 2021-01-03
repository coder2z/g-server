module github.com/myxy99/component

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/aliyun/aliyun-oss-go-sdk v2.1.5+incompatible
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.4.4
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.3
	github.com/json-iterator/go v1.1.10
	github.com/mitchellh/mapstructure v1.4.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/philchia/agollo/v4 v4.1.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d
	google.golang.org/grpc v1.34.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.0.3
	gorm.io/driver/postgres v1.0.6
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.9
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.19.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
