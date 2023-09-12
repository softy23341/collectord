package services

// SystemMailFrom TBD
const SystemMailFrom = "support@collectionsapp.ru" // help@clr.su

type (
	// Mail TBD
	Mail struct {
		To      []string
		From    string
		Subject string
		Body    string
	}

	// MailSender TBD
	MailSender interface {
		Send(mail *Mail) error
	}
)
