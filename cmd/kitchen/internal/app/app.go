package app

type App struct {
	repo  Repo
	queue Queue
	cron  Cron
}

func New(r Repo, q Queue, c Cron) *App {
	return &App{
		repo:  r,
		queue: q,
		cron: c,
	}
}
