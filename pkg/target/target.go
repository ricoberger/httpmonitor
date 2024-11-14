package target

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Method   string        `yaml:"method"`
	Body     string        `yaml:"body"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Token    string        `yaml:"token"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Client interface {
	Start()
	Name() string
	Results() []Result
	LastResult() Result
}

type client struct {
	config       Config
	httpClient   *http.Client
	resultsMutex sync.RWMutex
	results      []Result
}

func (c *client) Start() {
	c.check()

	for {
		<-time.After(c.config.Interval)
		c.check()
	}
}

func (c *client) Name() string {
	return c.config.Name
}

func (c *client) Results() []Result {
	c.resultsMutex.RLock()
	defer c.resultsMutex.RUnlock()

	return c.results
}

func (c *client) LastResult() Result {
	c.resultsMutex.RLock()
	defer c.resultsMutex.RUnlock()

	if len(c.results) == 0 {
		return Result{}
	}

	return c.results[len(c.results)-1]
}

func (c *client) check() {
	var statusCode int
	var result Result

	defer func() {
		result.End(time.Now(), statusCode)

		c.resultsMutex.Lock()
		defer c.resultsMutex.Unlock()

		c.results = append(c.results, result)

		if len(c.results) > 3600 {
			c.results = c.results[1:]
		}
	}()

	var body io.Reader
	if c.config.Body != "" {
		body = strings.NewReader(c.config.Body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, c.config.Method, c.config.URL, body)
	if err != nil {
		return
	}

	req = req.WithContext(withTrace(ctx, &result))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		return
	}

	statusCode = resp.StatusCode
}

func NewClient(config Config) Client {
	roundTripper := DefaultRoundTripper

	if config.Username != "" && config.Password != "" {
		roundTripper = BasicAuthTransport{
			Transport: roundTripper,
			Username:  config.Username,
			Password:  config.Password,
		}
	}

	if config.Token != "" {
		roundTripper = TokenAuthTransporter{
			Transport: roundTripper,
			Token:     config.Token,
		}
	}

	return &client{
		config: config,
		httpClient: &http.Client{
			Transport: roundTripper,
		},
	}
}
