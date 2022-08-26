package usecase

type Calculator interface {
	Run()
	SendMail(data templateData) error
}
