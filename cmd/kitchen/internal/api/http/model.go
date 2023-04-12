package http

import "xsolla/cmd/kitchen/internal/app"

func toAppCookingStatus(in string) (app.CookingStatus, error) {
	switch in {
	case app.CookingStatusNew.String():
		return app.CookingStatusNew, nil
	case app.CookingStatusInProgress.String():
		return app.CookingStatusInProgress, nil
	case app.CookingStatusNeedToStart.String():
		return app.CookingStatusNeedToStart, nil
	case app.CookingStatusCompleted.String():
		return app.CookingStatusCompleted, nil
	default:
		return 0, app.ErrInvalidArgument
	}
}
