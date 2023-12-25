package psg

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Psg struct {
	conn *pgxpool.Pool
}

func NewPsg(dburl string, login, pass string) (psg *Psg, err error) {
	defer func() { err = errors.Wrap(err, "postgres NewPsg()") }()

	psg = &Psg{}
	psg.conn, err = parseConnectionString(dburl, login, pass)
	if err != nil {
		return nil, err
	}

	err = psg.conn.Ping(context.Background())
	if err != nil {
		err = errors.Wrap(err, "psg.conn.Ping(context.Background())")
		return nil, err
	}

	return

}

func (psg *Psg) Close() {
	psg.conn.Close()
}

func parseConnectionString(dburl, user, password string) (db *pgxpool.Pool, err error) {
	var u *url.URL
	if u, err = url.Parse(dburl); err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга url строки")
	}
	u.User = url.UserPassword(user, password)
	db, err = pgxpool.New(context.Background(), u.String())
	if err != nil {
		return nil, errors.Wrap(err, "ошибка соединения с базой данных")
	}
	return
}
