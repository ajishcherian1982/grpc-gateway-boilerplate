package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ajishcherian1982/grpc-gateway-boilerplate/db"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/store"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

func setUp(t *testing.T) (*Handler, func(t *testing.T)) {
	w := zerolog.ConsoleWriter{Out: ioutil.Discard}
	// w := zerolog.ConsoleWriter{Out: os.Stderr}
	l := zerolog.New(w).With().Timestamp().Logger()

	d, err := db.NewTestDB()
	if err != nil {
		t.Fatal(fmt.Errorf("failed to initialize database: %w", err))
	}

	us := store.NewUserStore(d)
	as := store.NewArticleStore(d)

	return New(&l, us, as), func(t *testing.T) {
		err := db.DropTestDB(d)
		if err != nil {
			t.Fatal(fmt.Errorf("failed to clean database: %w", err))
		}
	}
}

func ctxWithToken(ctx context.Context, token string) context.Context {
	scheme := "Token"
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %s", scheme, token))
	nCtx := metautils.NiceMD(md).ToIncoming(ctx)
	return nCtx
}
