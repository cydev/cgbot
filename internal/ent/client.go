// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/cydev/cgbot/internal/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/cydev/cgbot/internal/ent/telegramchannel"
	"github.com/cydev/cgbot/internal/ent/telegramsession"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// TelegramChannel is the client for interacting with the TelegramChannel builders.
	TelegramChannel *TelegramChannelClient
	// TelegramSession is the client for interacting with the TelegramSession builders.
	TelegramSession *TelegramSessionClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.TelegramChannel = NewTelegramChannelClient(c.config)
	c.TelegramSession = NewTelegramSessionClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		TelegramChannel: NewTelegramChannelClient(cfg),
		TelegramSession: NewTelegramSessionClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:             ctx,
		config:          cfg,
		TelegramChannel: NewTelegramChannelClient(cfg),
		TelegramSession: NewTelegramSessionClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		TelegramChannel.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.TelegramChannel.Use(hooks...)
	c.TelegramSession.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.TelegramChannel.Intercept(interceptors...)
	c.TelegramSession.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *TelegramChannelMutation:
		return c.TelegramChannel.mutate(ctx, m)
	case *TelegramSessionMutation:
		return c.TelegramSession.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// TelegramChannelClient is a client for the TelegramChannel schema.
type TelegramChannelClient struct {
	config
}

// NewTelegramChannelClient returns a client for the TelegramChannel from the given config.
func NewTelegramChannelClient(c config) *TelegramChannelClient {
	return &TelegramChannelClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `telegramchannel.Hooks(f(g(h())))`.
func (c *TelegramChannelClient) Use(hooks ...Hook) {
	c.hooks.TelegramChannel = append(c.hooks.TelegramChannel, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `telegramchannel.Intercept(f(g(h())))`.
func (c *TelegramChannelClient) Intercept(interceptors ...Interceptor) {
	c.inters.TelegramChannel = append(c.inters.TelegramChannel, interceptors...)
}

// Create returns a builder for creating a TelegramChannel entity.
func (c *TelegramChannelClient) Create() *TelegramChannelCreate {
	mutation := newTelegramChannelMutation(c.config, OpCreate)
	return &TelegramChannelCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of TelegramChannel entities.
func (c *TelegramChannelClient) CreateBulk(builders ...*TelegramChannelCreate) *TelegramChannelCreateBulk {
	return &TelegramChannelCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TelegramChannelClient) MapCreateBulk(slice any, setFunc func(*TelegramChannelCreate, int)) *TelegramChannelCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TelegramChannelCreateBulk{err: fmt.Errorf("calling to TelegramChannelClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TelegramChannelCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TelegramChannelCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for TelegramChannel.
func (c *TelegramChannelClient) Update() *TelegramChannelUpdate {
	mutation := newTelegramChannelMutation(c.config, OpUpdate)
	return &TelegramChannelUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TelegramChannelClient) UpdateOne(tc *TelegramChannel) *TelegramChannelUpdateOne {
	mutation := newTelegramChannelMutation(c.config, OpUpdateOne, withTelegramChannel(tc))
	return &TelegramChannelUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TelegramChannelClient) UpdateOneID(id int64) *TelegramChannelUpdateOne {
	mutation := newTelegramChannelMutation(c.config, OpUpdateOne, withTelegramChannelID(id))
	return &TelegramChannelUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for TelegramChannel.
func (c *TelegramChannelClient) Delete() *TelegramChannelDelete {
	mutation := newTelegramChannelMutation(c.config, OpDelete)
	return &TelegramChannelDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TelegramChannelClient) DeleteOne(tc *TelegramChannel) *TelegramChannelDeleteOne {
	return c.DeleteOneID(tc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TelegramChannelClient) DeleteOneID(id int64) *TelegramChannelDeleteOne {
	builder := c.Delete().Where(telegramchannel.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TelegramChannelDeleteOne{builder}
}

// Query returns a query builder for TelegramChannel.
func (c *TelegramChannelClient) Query() *TelegramChannelQuery {
	return &TelegramChannelQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTelegramChannel},
		inters: c.Interceptors(),
	}
}

// Get returns a TelegramChannel entity by its id.
func (c *TelegramChannelClient) Get(ctx context.Context, id int64) (*TelegramChannel, error) {
	return c.Query().Where(telegramchannel.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TelegramChannelClient) GetX(ctx context.Context, id int64) *TelegramChannel {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *TelegramChannelClient) Hooks() []Hook {
	return c.hooks.TelegramChannel
}

// Interceptors returns the client interceptors.
func (c *TelegramChannelClient) Interceptors() []Interceptor {
	return c.inters.TelegramChannel
}

func (c *TelegramChannelClient) mutate(ctx context.Context, m *TelegramChannelMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TelegramChannelCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TelegramChannelUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TelegramChannelUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TelegramChannelDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown TelegramChannel mutation op: %q", m.Op())
	}
}

// TelegramSessionClient is a client for the TelegramSession schema.
type TelegramSessionClient struct {
	config
}

// NewTelegramSessionClient returns a client for the TelegramSession from the given config.
func NewTelegramSessionClient(c config) *TelegramSessionClient {
	return &TelegramSessionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `telegramsession.Hooks(f(g(h())))`.
func (c *TelegramSessionClient) Use(hooks ...Hook) {
	c.hooks.TelegramSession = append(c.hooks.TelegramSession, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `telegramsession.Intercept(f(g(h())))`.
func (c *TelegramSessionClient) Intercept(interceptors ...Interceptor) {
	c.inters.TelegramSession = append(c.inters.TelegramSession, interceptors...)
}

// Create returns a builder for creating a TelegramSession entity.
func (c *TelegramSessionClient) Create() *TelegramSessionCreate {
	mutation := newTelegramSessionMutation(c.config, OpCreate)
	return &TelegramSessionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of TelegramSession entities.
func (c *TelegramSessionClient) CreateBulk(builders ...*TelegramSessionCreate) *TelegramSessionCreateBulk {
	return &TelegramSessionCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TelegramSessionClient) MapCreateBulk(slice any, setFunc func(*TelegramSessionCreate, int)) *TelegramSessionCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TelegramSessionCreateBulk{err: fmt.Errorf("calling to TelegramSessionClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TelegramSessionCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TelegramSessionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for TelegramSession.
func (c *TelegramSessionClient) Update() *TelegramSessionUpdate {
	mutation := newTelegramSessionMutation(c.config, OpUpdate)
	return &TelegramSessionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TelegramSessionClient) UpdateOne(ts *TelegramSession) *TelegramSessionUpdateOne {
	mutation := newTelegramSessionMutation(c.config, OpUpdateOne, withTelegramSession(ts))
	return &TelegramSessionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TelegramSessionClient) UpdateOneID(id int) *TelegramSessionUpdateOne {
	mutation := newTelegramSessionMutation(c.config, OpUpdateOne, withTelegramSessionID(id))
	return &TelegramSessionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for TelegramSession.
func (c *TelegramSessionClient) Delete() *TelegramSessionDelete {
	mutation := newTelegramSessionMutation(c.config, OpDelete)
	return &TelegramSessionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TelegramSessionClient) DeleteOne(ts *TelegramSession) *TelegramSessionDeleteOne {
	return c.DeleteOneID(ts.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TelegramSessionClient) DeleteOneID(id int) *TelegramSessionDeleteOne {
	builder := c.Delete().Where(telegramsession.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TelegramSessionDeleteOne{builder}
}

// Query returns a query builder for TelegramSession.
func (c *TelegramSessionClient) Query() *TelegramSessionQuery {
	return &TelegramSessionQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTelegramSession},
		inters: c.Interceptors(),
	}
}

// Get returns a TelegramSession entity by its id.
func (c *TelegramSessionClient) Get(ctx context.Context, id int) (*TelegramSession, error) {
	return c.Query().Where(telegramsession.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TelegramSessionClient) GetX(ctx context.Context, id int) *TelegramSession {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *TelegramSessionClient) Hooks() []Hook {
	return c.hooks.TelegramSession
}

// Interceptors returns the client interceptors.
func (c *TelegramSessionClient) Interceptors() []Interceptor {
	return c.inters.TelegramSession
}

func (c *TelegramSessionClient) mutate(ctx context.Context, m *TelegramSessionMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TelegramSessionCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TelegramSessionUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TelegramSessionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TelegramSessionDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown TelegramSession mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		TelegramChannel, TelegramSession []ent.Hook
	}
	inters struct {
		TelegramChannel, TelegramSession []ent.Interceptor
	}
)
