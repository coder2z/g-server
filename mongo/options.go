/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package mongo

type options struct {
	Debug    bool   `mapStructure:"debug"`
	URL      string `mapStructure:"url"`
	Source   string `mapStructure:"source"`
	User     string `mapStructure:"user"`
	Password string `mapStructure:"password"`
}

func newMongoOptions() *options {
	return &options{
		Debug:    true,
		URL:      "",
		Source:   "",
		User:     "",
		Password: "",
	}
}
