package transfer

import "github.com/hassanalgoz/swe/internal/entities"

func isFreezed(a *entities.Account) bool {
	return a.FreezedSince != nil
}
