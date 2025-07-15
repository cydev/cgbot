package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/app"
	"github.com/go-faster/sdk/zctx"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	entdb "github.com/cydev/cgbot/internal/db"
	"github.com/cydev/cgbot/internal/ent"
	"github.com/cydev/cgbot/internal/oas"
)

type Application struct {
	db     *ent.Client
	tg     *tg.Client
	api    *oas.Client
	client *telegram.Client

	waiter *floodwait.Waiter
	trace  trace.Tracer
}

func (a *Application) Run(ctx context.Context) error {
	return a.waiter.Run(ctx, func(ctx context.Context) error {
		return a.client.Run(ctx, func(ctx context.Context) error {
			lg := zctx.From(ctx)
			if self, err := a.client.Self(ctx); err != nil || !self.Bot {
				if _, err := a.client.Auth().Bot(ctx, os.Getenv("BOT_TOKEN")); err != nil {
					return errors.Wrap(err, "auth bot")
				}
			} else {
				lg.Info("Already authenticated")
			}
			if self, err := a.client.Self(ctx); err == nil {
				lg.Info("Bot info",
					zap.Int64("id", self.ID),
					zap.String("username", self.Username),
					zap.String("first_name", self.FirstName),
					zap.String("last_name", self.LastName),
				)
			}
			if _, err := a.tg.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
				Scope:    &tg.BotCommandScopeDefault{},
				LangCode: "en",
				Commands: []tg.BotCommand{
					{
						Command:     "start",
						Description: "Start bot",
					},
				},
			}); err != nil {
				return errors.Wrap(err, "set commands")
			}
			<-ctx.Done()
			return ctx.Err()
		})
	})
}

func (a *Application) addChannel(ctx context.Context, channel *tg.Channel) error {
	return a.db.TelegramChannel.Create().
		SetID(channel.ID).
		SetAccessHash(channel.AccessHash).
		SetTitle(channel.Title).
		SetActive(true).
		Exec(ctx)
}

func (a *Application) removeChannel(ctx context.Context, channel *tg.Channel) error {
	if err := a.db.TelegramChannel.UpdateOneID(channel.ID).
		SetActive(false).
		Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return nil
		}
		return errors.Wrap(err, "update channel")
	}
	return nil
}

func (a *Application) onChannelParticipant(ctx context.Context, e tg.Entities, update *tg.UpdateChannelParticipant) error {
	switch update.NewParticipant.(type) {
	case *tg.ChannelParticipantBanned:
		// Bot was removed from channel.
		for _, c := range e.Channels {
			return a.removeChannel(ctx, c)
		}
	case *tg.ChannelParticipantAdmin:
		// Bot was added to channel.
		for _, c := range e.Channels {
			return a.addChannel(ctx, c)
		}
	default:
		if update.NewParticipant == nil {
			// Removed from channel.
			for _, c := range e.Channels {
				return a.removeChannel(ctx, c)
			}
		}
	}
	return nil
}

func (a *Application) onNewMessage(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
	ctx, span := a.trace.Start(ctx, "OnNewMessage")
	defer span.End()
	m, ok := u.Message.(*tg.Message)
	if !ok || m.Out {
		return nil
	}
	var (
		sender = message.NewSender(a.tg)
		reply  = sender.Reply(e, u)
		lg     = zctx.From(ctx).With(zap.Int("msg.id", m.ID))
	)
	peerUser, ok := m.PeerID.(*tg.PeerUser)
	if !ok {
		if _, err := reply.Text(ctx, "Invalid"); err != nil {
			return err
		}
		return nil
	}
	user := e.Users[peerUser.UserID]
	if user == nil {
		return nil
	}
	lg.Info("New message",
		zap.String("text", m.Message),
		zap.String("user", user.Username),
		zap.String("first_name", user.FirstName),
		zap.String("last_name", user.LastName),
		zap.Int64("user_id", user.ID),
	)
	switch {
	case m.Message == "/start":
		if _, err := reply.Text(ctx, "Hello, "+user.FirstName+"!"); err != nil {
			return errors.Wrap(err, "send message")
		}
	case strings.HasPrefix(m.Message, "/profile"):
		// /profile <gameName>#<tagLine>
		parts := strings.SplitN(m.Message, " ", 2)
		if len(parts) < 2 {
			if _, err := reply.Text(ctx, "Usage: /profile <gameName>#<tagLine>"); err != nil {
				return errors.Wrap(err, "send message")
			}
			return nil
		}
		gameName, tagLine := parts[1], ""
		if idx := strings.Index(gameName, "#"); idx != -1 {
			tagLine = gameName[idx+1:]
			gameName = gameName[:idx]
		}
		if tagLine == "" {
			if _, err := reply.Text(ctx, "Usage: /profile <gameName>#<tagLine>"); err != nil {
				return errors.Wrap(err, "send message")
			}
		}
		res, err := a.api.AccountV1GetByRiotId(ctx, oas.AccountV1GetByRiotIdParams{
			GameName: gameName,
			TagLine:  tagLine,
		})
		if err != nil {
			if _, err := reply.Text(ctx, "Error: "+err.Error()); err != nil {
				return errors.Wrap(err, "send message")
			}
			return nil
		}
		switch res := res.(type) {
		case *oas.AccountV1AccountDto:
			if _, err := reply.Text(ctx, "Account found: "+res.Puuid); err != nil {
				return errors.Wrap(err, "send message")
			}
		default:
			lg.Error("Unexpected response type",
				zap.Any("value", res),
				zap.String("type", fmt.Sprintf("%T", res)),
				zap.String("gameName", gameName),
				zap.String("tagLine", tagLine),
			)
			if _, err := reply.Text(ctx, "Unexpected response type"); err != nil {
				return errors.Wrap(err, "send message")
			}
		}
	}
	return nil
}

var _ oas.SecuritySource = (*securityProvider)(nil)

type securityProvider struct{}

// APIKey implements oas.SecuritySource.
func (s *securityProvider) APIKey(ctx context.Context, operationName oas.OperationName) (oas.APIKey, error) {
	return oas.APIKey{
		APIKey: os.Getenv("RIOT_API_KEY"),
	}, nil
}

// Rso implements oas.SecuritySource.
func (s *securityProvider) Rso(ctx context.Context, operationName oas.OperationName) (oas.Rso, error) {
	return oas.Rso{}, nil
}

// XRiotToken implements oas.SecuritySource.
func (s *securityProvider) XRiotToken(ctx context.Context, operationName oas.OperationName) (oas.XRiotToken, error) {
	return oas.XRiotToken{}, nil
}

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger, t *app.Telemetry) error {
		riotAPI, err := oas.NewClient("https://europe.api.riotgames.com",
			&securityProvider{},
			oas.WithMeterProvider(t.MeterProvider()),
			oas.WithTracerProvider(t.TracerProvider()),
		)
		if err != nil {
			return errors.Wrap(err, "create riot api client")
		}
		db, err := entdb.Open(ctx, os.Getenv("DATABASE_URL"), t)
		if err != nil {
			return errors.Wrap(err, "open database")
		}
		botToken := os.Getenv("BOT_TOKEN")
		if botToken == "" {
			return errors.New("BOT_TOKEN is empty")
		}
		appID, err := strconv.Atoi(os.Getenv("APP_ID"))
		if err != nil {
			return errors.Wrap(err, "parse APP_ID")
		}
		appHash := os.Getenv("APP_HASH")
		if appHash == "" {
			return errors.New("APP_HASH is empty")
		}
		waiter := floodwait.NewWaiter()
		dispatcher := tg.NewUpdateDispatcher()
		client := telegram.NewClient(appID, appHash, telegram.Options{
			Logger:         zctx.From(ctx).Named("tg"),
			TracerProvider: t.TracerProvider(),
			SessionStorage: entdb.NewSessionStorage(-1, db),
			UpdateHandler:  dispatcher,
			Middlewares: []telegram.Middleware{
				waiter,
			},
		})
		a := &Application{
			db:     db,
			tg:     tg.NewClient(client),
			api:    riotAPI,
			client: client,
			waiter: waiter,
			trace:  t.TracerProvider().Tracer("recam.bot"),
		}
		dispatcher.OnChannelParticipant(a.onChannelParticipant)
		dispatcher.OnNewMessage(a.onNewMessage)
		return a.Run(ctx)
	})
}
