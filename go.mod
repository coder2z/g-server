module github.com/coder2z/component

go 1.15

require (
	github.com/HdrHistogram/hdrhistogram-go v1.0.1 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.974
	github.com/aliyun/aliyun-oss-go-sdk v2.1.6+incompatible
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/coder2z/g-saber v1.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/go-redis/redis/v8 v8.7.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/philchia/agollo/v4 v4.1.3
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200425165423-262c93980547
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb
	google.golang.org/grpc v1.31.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.0.4
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.2
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.19.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.29.0
