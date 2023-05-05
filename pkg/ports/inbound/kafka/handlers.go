package kafka

func (c *consumer) handle(e Message) {
	switch e.Type {
	case "DoThis":
		c.doThis()
	case "DoThat":
		c.doThat()
	}
}

func (c *consumer) doThis() {

}

func (c *consumer) doThat() {
}
