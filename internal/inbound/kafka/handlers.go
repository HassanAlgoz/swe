package kafka

import (
	"log"

	"github.com/google/uuid"
)

func (c *consumer) handle(e Message) {
	switch e.Type {
	case "DoThis":
		c.doThis()
	case "DoThat":
		c.doThat()
	}
}

func (c *consumer) doThis() {
	from := uuid.New()
	to := uuid.New()
	amount := int64(100_00)
	err := c.app.MoneyTransfer(from, to, amount)
	if err != nil {
		log.Printf("%% Error: %v\n", err)
	}
}

func (c *consumer) doThat() {
}
