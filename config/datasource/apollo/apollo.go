package apollo

import (
	"github.com/myxy99/component/config"
	"github.com/philchia/agollo/v4"
)

type apolloDataSource struct {
	client      agollo.Client
	namespace   string
	propertyKey string
	changed     chan struct{}
}

// NewDataSource creates an apolloDataSource
func NewDataSource(conf *agollo.Conf, namespace string, key string) config.DataSource {
	client := agollo.NewClient(conf, agollo.WithLogger(&agolloLogger{}))
	ap := &apolloDataSource{
		client:      client,
		namespace:   namespace,
		propertyKey: key,
		changed:     make(chan struct{}, 1),
	}
	_ = ap.client.Start()
	ap.client.OnUpdate(
		func(event *agollo.ChangeEvent) {
			ap.changed <- struct{}{}
		})
	return ap
}

// ReadConfig reads config content from apollo
func (ap *apolloDataSource) ReadConfig() ([]byte, error) {
	value := ap.client.GetString(ap.propertyKey, agollo.WithNamespace(ap.namespace))
	return []byte(value), nil
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ap *apolloDataSource) IsConfigChanged() <-chan struct{} {
	return ap.changed
}

// Close stops watching the config changed
func (ap *apolloDataSource) Close() error {
	_ = ap.client.Stop()
	close(ap.changed)
	return nil
}

type agolloLogger struct {
}

// Infof ...
func (l *agolloLogger) Infof(format string, args ...interface{}) {
}

// Errorf ...
func (l *agolloLogger) Errorf(format string, args ...interface{}) {
}
