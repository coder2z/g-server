/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 11:18
 */
package email

type options struct {
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Username string `json:"username,omitempty" yaml:"username"`
	Password string `json:"-" yaml:"password"`
}

func newEmailOptions() *options {
	return &options{
		Host:     "smtp.yeah.net",
		Port:     465,
		Username: "yangzzzzzz@yeah.net",
		Password: "your-password",
	}
}
