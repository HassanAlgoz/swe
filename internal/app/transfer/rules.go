package transfer

import (
	"fmt"
	"time"

	"github.com/hassanalgoz/swe/internal/ent"
	"github.com/looplab/fsm"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

type RTPEvent string

const (
	RTPEventReject RTPEvent = "Reject"
	RTPEventAccept RTPEvent = "Accept"
	RTPEventCancel RTPEvent = "Cancel"
	RTPEventExpire RTPEvent = "Expire"
)

var RTPEvents = fsm.Events{
	{Src: []string{string(ent.RTPStateInitiated)}, Name: string(RTPEventReject), Dst: string(ent.RTPStateRejected)},
	{Src: []string{string(ent.RTPStateInitiated)}, Name: string(RTPEventAccept), Dst: string(ent.RTPStateAccepted)},
	{Src: []string{string(ent.RTPStateInitiated)}, Name: string(RTPEventCancel), Dst: string(ent.RTPStateCancelled)},
	{Src: []string{string(ent.RTPStateInitiated)}, Name: string(RTPEventExpire), Dst: string(ent.RTPStateExpired)},
}

var discountEligibilityDuration = viper.GetDuration("app.discount_eligibility_duration")

func isEligibleForDiscount(c *ent.Customer) (bool, string) {
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

func isFreezed(a *ent.Account) (bool, string) {
	if a.FreezedSince != nil {
		return true, fmt.Sprintf("freezed since %s", a.FreezedSince.Format("2006-01-02"))
	}
	return false, ""
}

func isValidTransition(s1 ent.RTPState, s2 ent.RTPState) (bool, string) {
	fsm := fsm.NewFSM(
		string(s1),
		RTPEvents,
		fsm.Callbacks{},
	)
	index := slices.IndexFunc(fsm.AvailableTransitions(), func(v string) bool {
		return v == string(s2)
	})
	if index == -1 {
		return false, fmt.Sprintf("transition from %s to %s is not allowed", string(s1), string(s2))
	}
	return true, ""
}
