package sdk

import (
	"aletheia-go-sdk/sdk/config"
	"aletheia-go-sdk/sdk/models"
	"aletheia-go-sdk/sdk/services"
	"aletheia-go-sdk/sdk/transport"
	"time"
)

type Client struct {
	config      *config.AletheiaConfig
	eventSvc    *services.EventService
	resourceSvc *services.ResourceService
	userID      string
	userInfo    map[string]interface{}
	tags        map[string]string
	customField map[string]interface{}
}

func NewClient(cfg config.AletheiaConfig) (*Client, error) {
	if cfg.Logger == nil {
		cfg.Logger = config.NewDefaultLogger()
	}

	tp, err := transport.NewHTTPTransport(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Logger.Infof("SDK initialized with config: %+v", cfg)

	return &Client{
		config:      &cfg,
		eventSvc:    services.NewEventService(cfg, tp, cfg.Logger),
		resourceSvc: services.NewResourceService(cfg, tp, cfg.Logger),
		tags:        make(map[string]string),
		customField: make(map[string]interface{}),
	}, nil
}

func (c *Client) Close() error {
	c.resourceSvc.StopMonitoring()
	return nil
}

func (c *Client) MonitorResources(interval time.Duration, metricsCollector func() map[string]interface{}) {
	c.resourceSvc.StartMonitoring(interval, c.tags, c.userID, c.userInfo, metricsCollector)
}

func (c *Client) CaptureMessage(message string, level models.LogLevel, context map[string]interface{}) {
	c.eventSvc.CaptureMessage(message, level, context, c.tags, c.userID, c.userInfo)
}

func (c *Client) CaptureException(err error, eventMessage string, eventType string, context map[string]interface{}) {
	c.eventSvc.CaptureException(err, eventMessage, eventType, context, c.tags, c.userID, c.userInfo, c.customField)
}

func (c *Client) SetUser(userID string, info map[string]interface{}) {
	c.userID = userID
	c.userInfo = info
}

func (c *Client) SetTags(tags map[string]string) {
	for k, v := range tags {
		c.tags[k] = v
	}
}

func (c *Client) SetCustomField(key string, value interface{}) {
	c.customField["fields."+key] = value
}
