package app

type App struct {
	repo             Repo
	queue            Queue
	cron             Cron
	addressValidator AddressValidator
}

func New(r Repo, q Queue, c Cron, a AddressValidator) *App {
	return &App{
		repo:             r,
		queue:            q,
		cron:             c,
		addressValidator: a,
	}
}
