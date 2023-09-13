package email

import (
	"github.com/go-mail/mail"
	"io"
)

const (
	AttachmentReaderType = 0
	AttachmentFileType   = 1
)

type Option func(*MailOption)

type MailOption struct {
	AttachMents []AttachMent
}

type AttachMent struct {
	Name     string
	Type     int
	Content  io.Reader
	Settings []mail.FileSetting
}

func (e *Email) WithAttachReader(name string, reader io.Reader, settings ...mail.FileSetting) Option {
	return func(opt *MailOption) {
		opt.AttachMents = append(opt.AttachMents,
			AttachMent{
				Name:     name,
				Type:     AttachmentReaderType,
				Content:  reader,
				Settings: settings,
			},
		)
	}
}

func (e *Email) WithAttachFile(name string, settings ...mail.FileSetting) Option {
	return func(opt *MailOption) {
		opt.AttachMents = append(opt.AttachMents,
			AttachMent{
				Name:     name,
				Type:     AttachmentFileType,
				Settings: settings,
			},
		)
	}
}

func (e *Email) MailBuilder(m *mail.Message, opts ...Option) *MailOption {
	mailOption := &MailOption{}
	for _, apply := range opts {
		if apply != nil {
			apply(mailOption)
		}
	}

	for _, attachMent := range mailOption.AttachMents {
		if attachMent.Type == AttachmentReaderType {
			m.AttachReader(attachMent.Name, attachMent.Content, attachMent.Settings...)
		} else {
			m.Attach(attachMent.Name, attachMent.Settings...)
		}
	}
	return mailOption
}
