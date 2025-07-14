package main

import (
	"context"
	"os"
	"strconv"

	entdb "github.com/cydev/cgbot/internal/db"
	"github.com/cydev/cgbot/internal/ent"
	"github.com/go-faster/sdk/app"
	"go.uber.org/zap"

	"github.com/go-faster/errors"
	"github.com/go-faster/sdk/zctx"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"go.opentelemetry.io/otel/trace"
)

type Application struct {
	db     *ent.Client
	api    *tg.Client
	client *telegram.Client "go.uber.org/zap"

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
			if _, err := a.api.BotsSetBotCommands(ctx, &tg.BotsSetBotCommandsRequest{
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

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger, t *app.Telemetry) error {
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
			api:    tg.NewClient(client),
			client: client,
			waiter: waiter,
			trace:  t.TracerProvider().Tracer("recam.bot"),
		}
		dispatcher.OnChannelParticipant(a.onChannelParticipant)
		return a.Run(ctx)
	})
}
