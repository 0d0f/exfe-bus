package gobus

type servicePair struct {
	service *Client
	method  string
}

type TubeClient struct {
	name     string
	services []servicePair
}

func NewTubeClient(name string) *TubeClient {
	return &TubeClient{
		name:     name,
		services: make([]servicePair, 0),
	}
}

func (c *TubeClient) Name() string {
	return c.name
}

func (c *TubeClient) AddService(service, method string) error {
	client, err := NewClient(service)
	if err != nil {
		return err
	}
	pair := servicePair{
		service: client,
		method:  method,
	}
	c.services = append(c.services, pair)
	return nil
}

func (c *TubeClient) Send(arg interface{}) error {
	for _, p := range c.services {
		var reply interface{}
		err := p.service.Do(p.method, arg, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}
