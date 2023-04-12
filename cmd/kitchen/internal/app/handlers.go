package app

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
)

func (a *App) ListCooking(ctx context.Context, params CookingParams) ([]Cooking, int, error) {
	return a.repo.ListCooking(ctx, params)
}

func (a *App) ChangeCookingStatus(ctx context.Context, cookingID uuid.UUID, status CookingStatus) error {
	cooking, err := a.repo.GetCooking(ctx, cookingID)
	if err != nil {
		return fmt.Errorf("a.repo.GetOrder: %w", err)
	}

	if cooking.Status == status {
		return ErrSameStatus
	}

	err = a.repo.Tx(ctx, func(repo Repo) error {
		updatedCooking, err := a.repo.UpdateCooking(ctx, Cooking{
			Status: status,
		})
		if err != nil {
			return fmt.Errorf("a.repo.UpdateCooking: %w", err)
		}

		if status == CookingStatusInProgress || status == CookingStatusCompleted {
			_, err = a.repo.SaveTask(ctx, Task{
				Cooking:  *updatedCooking,
				TaskKind: TaskKindEventUpdateCooking,
			})
			if err != nil {
				return fmt.Errorf("a.repo.SaveTask: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
