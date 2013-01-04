package gobus

type servicePair struct {
	url    string
	method string
}

type TubeClient struct {
	name     string
	client   *Client
	services []servicePair
}

func NewTubeClient(name string) *TubeClient {
	return &TubeClient{
		name:     name,
		client:   NewClient(new(JSON)),
		services: make([]servicePair, 0),
	}
}

func (c *TubeClient) Name() string {
	return c.name
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
		err := c.client.Do(p.url, p.method, arg, &reply)
		if err != nil {
			return err
		}
	}
	return nil
}
