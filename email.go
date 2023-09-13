package email

import (
	"github.com/go-mail/mail"
	"github.com/go-tron/config"
	"github.com/go-tron/logger"
)

type Config struct {
	Host   string
	Port   int
	User   string
	Pass   string
	Logger logger.Logger
}

func New(c *Config) *Email {
	if c.Host == "" || c.Port == 0 || c.User == "" || c.Pass == "" {
		panic("invalid config")
	}
	return &Email{c}
}

func NewWithConfig(c *config.Config) *Email {
	return New(&Config{
		Host:   c.GetString("email.host"),
		Port:   c.GetInt("email.port"),
		User:   c.GetString("email.user"),
		Pass:   c.GetString("email.pass"),
		Logger: logger.NewZapWithConfig(c, "email", "error"),
	})
}

type Email struct {
	*Config
}

func (e *Email) Send(to string, title string, body string, opts ...Option) error {
	return e.SendToMany([]string{to}, title, body, opts...)
}

func (e *Email) SendToMany(to []string, title string, body string, opts ...Option) error {
	m := mail.NewMessage()
	m.SetHeader("From", e.User)   //这种方式可以添加别名
	m.SetHeader("To", to...)      //发送给多个用户
	m.SetHeader("Subject", title) //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文
	_ = e.MailBuilder(m, opts...)
	d := mail.NewDialer(e.Host, e.Port, e.User, e.Pass)
	err := d.DialAndSend(m)

	if err != nil {
		e.Logger.Error("",
			e.Logger.Field("to", to),
			e.Logger.Field("title", title),
			e.Logger.Field("body", body),
			e.Logger.Field("err", err),
		)
	}
	return err
}
