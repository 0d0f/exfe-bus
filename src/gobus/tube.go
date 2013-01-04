package gobus

type servicePair struct {
	url    string
	method string
}

type TubeClient struct {
	dispatcher *Dispatcher
	services   []servicePair
}

func NewTubeClient(dispatcher *Dispatcher) *TubeClient {
	return &TubeClient{
		dispatcher: dispatcher,
		services:   make([]servicePair, 0),
	}
}

func (c *TubeClient) AddService(url, method string) error {
	pair := servicePair{
		url:    url,
		method: method,
	}
	c.services = append(c.services, pair)
	return nil
}

func (c *TubeClient) Send(arg interface{}) error {
	for _, p := range c.services {
		var reply interface{}
		err := c.dispatcher.Do(p.url, p.method, arg, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TubeClient) SendWithIdentity(identity string, arg interface{}) error {
	for _, p := range c.services {
		var reply interface{}
		err := c.dispatcher.DoWithIdentity(identity, p.url, p.method, arg, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}
