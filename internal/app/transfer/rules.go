package transfer

import (
	"fmt"
	"time"

	"github.com/hassanalgoz/swe/internal/common"
	"github.com/spf13/viper"
)

var discountEligibilityDuration = viper.GetDuration("app.discount_eligibility_duration")

func isEligibleForDiscount(c *common.Customer) (bool, string) {
	// must has accounts
	if !(len(c.Accounts) > 0) {
		return false, "customer must have accounts"
	}
	// must be registered for at least discountEligibilityDuration
	if !(time.Since(c.CreatedAt) > discountEligibilityDuration) {
		return false, fmt.Sprintf("customer must be registered for at least %v", discountEligibilityDuration)
	}
	return true, ""
}

func isFreezed(a *common.Account) (bool, string) {
	if a.FreezedSince != nil {
		return true, fmt.Sprintf("freezed since %s", a.FreezedSince.Format("2006-01-02"))
	}
	return false, ""
}
