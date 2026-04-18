package services

type MailQueue interface {
	PublishVerificationEmail(to, code string) error
	PublishLoginCodeEmail(to, code string) error
	StartConsumer()
}
