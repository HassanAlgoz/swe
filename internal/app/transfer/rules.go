package transfer

import (
	"fmt"
	"time"

	"github.com/hassanalgoz/swe/internal/common"
)

func isEligibleForDiscount(c *common.Customer) (bool, string) {
	// must has accounts
	if !(len(c.Accounts) > 0) {
		return false, "customer must have accounts"
	}
	// must be registered for at least 5 years
	if !(time.Since(c.CreatedAt) > 5*365*24*time.Hour) {
		return false, "customer must be registered for at least 5 years"
	}
	return true, ""
}

func isFreezed(a *common.Account) (bool, string) {
	if a.FreezedSince != nil {
		return true, fmt.Sprintf("freezed since %s", a.FreezedSince.Format("2006-01-02"))
	}
	return false, ""
}
