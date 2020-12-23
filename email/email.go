/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/5 0:20
 */
package email

import (
	"github.com/myxy99/component/config"
	"gopkg.in/gomail.v2"
)

type Email struct {
	o *options
}

func (i *emailInvoker) newEmail(options *options) *Email {
	return &Email{options}
}

func (e *Email) SendEmail(mailTo []string, subject string, body string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.o.Username, ""))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Cc", m.FormatAddress(e.o.Username, ""))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(
		e.o.Host,
		e.o.Port,
		e.o.Username,
		e.o.Password)
	err = d.DialAndSend(m)
	return err
}

func (i *emailInvoker) loadConfig() map[string]*options {
	conf := make(map[string]*options)

	prefix := i.key
	for name := range config.GetStringMap(prefix) {
		cfg := config.UnmarshalWithExpect(prefix+"."+name, newEmailOptions()).(*options)
		conf[name] = cfg
	}

	return conf
}
