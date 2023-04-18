package http

import (
	"fmt"
	"xsolla/cmd/kitchen/internal/app"
)

type requestListOfCooking struct {
	CookingStatus app.CookingStatus `json:"cooking_status" validate:"required"`
	Limit         uint16            `json:"limit" validate:"gt=0"`
	Offset        uint16            `json:"offset"`
}

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
		return 0, fmt.Errorf("unknown cooking status: %s: %w", in, app.ErrInvalidArgument)
	}
}
