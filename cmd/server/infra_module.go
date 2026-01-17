package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaratBR/openlibrary/internal/app"
	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/MaratBR/openlibrary/internal/app/email"
	"github.com/MaratBR/openlibrary/internal/csrf"
	elasticstore "github.com/MaratBR/openlibrary/internal/elastic-store"
	i18n "github.com/MaratBR/openlibrary/internal/i18n"
	"github.com/MaratBR/openlibrary/internal/session"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/internal/upload"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"github.com/mailgun/errors"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"golang.org/x/text/language"
)

var infraModule = fx.Module("infra", fx.Provide(
	loadConfigOrPanic,
	connectToDatabase,
	newLocaleProvider,
	createCache,
	createMailService,
	createOpensearch,
	func() *csrf.Handler {
		return csrf.NewHandler("CSRF HANDLER HERE")
	},
	newRedisClient,
	session.NewRedisStore,
	upload.NewFileValidator,
))

func newLocaleProvider(cliParams cliParams) *i18n.LocaleProvider {
	return i18n.NewLocaleProvider(
		language.English,
		cliParams.Dev,
		map[language.Tag][]string{
			language.English: {
				"translations/en.toml",
			},
		},
	)
}

func newRedisClient(cfg *koanf.Koanf) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.String("redis.addr"),
	})

	return client
}

func connectToDatabase(config *koanf.Koanf) app.DB {
	connectionString := config.String("database.url")
	if connectionString == "" {
		slog.Error("database.url is empty")
		os.Exit(1)
	}
	var (
		db       *pgxpool.Pool
		err      error
		sleeping = time.Second
	)

	for {
		db, err = store.Connect(context.Background(), connectionString)
		if err != nil {
			slog.Error("failed to connect to DB", "err", err, "sleeping", sleeping)
			time.Sleep(sleeping)
			if sleeping < time.Second*60 {
				sleeping *= 2
			}
		} else {
			break
		}
	}

	return db
}

func createMailService(cfg *koanf.Koanf) (email.Service, error) {
	type_ := cfg.String("mail.type")

	switch type_ {
	case "mailgun":
		{
			domain := cfg.String("mailgun.domain")
			key := cfg.String("mailgun.key")
			senderName := cfg.String("mailgun.senderName")
			isEU := cfg.Bool("mailgun.isEU")

			if key == "" {
			}

			slog.Debug("mailgun is being used as email service of choice", "domain", domain, "key", "[REDACTED... DUH]")

			mg, err := email.NewMailgunEmailService(domain, key, senderName, isEU)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create mailgun service")
			}
			return mg, nil
		}
	case "console":
		return email.NewConsole(), nil
	case "":
		slog.Warn("empty mail.type value - falling back to 'console' email service")
		return email.NewConsole(), nil
	default:
		return nil, fmt.Errorf("unknown email type: %s", type_)
	}
}

func createCacheBackend(cfg *koanf.Koanf) (cache.CacheBackend, error) {
	backend := cfg.String("cache.type")
	switch backend {
	case "disabled":
		return cache.NewDisabledCacheBackend(), nil
	case "memory":
		return cache.NewMemoryCacheBackend(), nil
	case "redis":
		url := cfg.String("redis.url")
		if strings.Trim(url, " \n\t") == "" {
			return nil, fmt.Errorf("redis.url is empty")
		}

		return cache.NewRedisCacheBackend(
			url,
			cache.NewMemoryCacheBackend(),
		), nil
	default:
		return nil, fmt.Errorf("unknown cache backend: %s", backend)
	}

}

func createCache(config *koanf.Koanf) *cache.Cache {
	cacheBackend, err := createCacheBackend(config)
	if err != nil {
		panic(err)
	}
	cacheInstance := cache.New(cacheBackend)
	return cacheInstance
}

func createOpensearch(config *koanf.Koanf) *opensearchapi.Client {
	elasticsearchURL := config.String("elasticsearch.url")
	if elasticsearchURL == "" {
		slog.Error("elasticsearch.url is empty")
		os.Exit(1)
	}

	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Addresses: []string{elasticsearchURL},
			Username:  config.String("opensearch.user"),
			Password:  config.String("opensearch.password"),
		},
	})
	if err != nil {
		panic(err)
	}

	go func() {
		err = elasticstore.Setup(context.Background(), client)
		if err != nil {
			slog.Error("FAILED TO SETUP ELASTIC", "err", err)
		}
	}()

	return client
}
