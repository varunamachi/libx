package pg

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/netx"
)

var defDB *sqlx.DB
var conns map[string]*sqlx.DB

// ConnOpts - postgres connection options
type ConnOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	TimeZone string `json:"timeZone"`
}

// String - get usable connection string
func (c *ConnOpts) String() string {
	// postgres://username:password@url.com:5432/dbName
	// "postgres://%s:%s@%s:%d/%s?sslmode=disable&TimeZone=%s",
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone='%s'",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.TimeZone,
	)
}

// Url - get a postgres URL
func (c *ConnOpts) Url() (*url.URL, error) {
	urlStr := fmt.Sprintf(
		// postgres://username:password@url.com:5432/dbName
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&timezon=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.TimeZone,
	)
	return url.Parse(urlStr)
}

// connect - connects to DB based on connection string or URL
func Connect(
	gtx context.Context, url *url.URL, timeZone string) (*sqlx.DB, error) {
	if err := netx.WaitForPorts(gtx, url.Host, 60*time.Second); err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", url.String())
	if err != nil {
		// log.Error().Err(err).Msg("failed to open postgress connection")
		return nil, errx.Errf(err, "failed to open postgress connection")
	}

	retry := 3
	for retry >= 0 {
		if err = db.Ping(); err != nil {
			retry--
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	if retry < 0 {
		return nil, errx.Errf(err, "failed to ping the database")
	}

	if timeZone != "" {
		_, err = db.Exec(fmt.Sprintf("SET TIME ZONE '%s'", timeZone))
		if err != nil {
			return nil, errx.Errf(err, "failed to set postgres timezone")
		}
	}
	var dbNow time.Time
	if err := db.Get(&dbNow, "SELECT now()"); err != nil {
		return nil, errx.Errf(err, "failed get current time")

	}
	log.Info().Str("DB.CurrentTime", dbNow.Format(time.UnixDate)).Msg("")
	return db, nil
}

func PrintDbTime() {
	var dbNow time.Time
	if err := Conn().Get(&dbNow, "SELECT now()"); err != nil {
		err = errx.Errf(err, "failed get current time")
		errx.PrintSomeStack(err)

	}
	log.Info().Str("DB.CurrentTime", dbNow.Format(time.UnixDate)).
		Msg("")
}

// ConnectWithOpts - connect to postgresdb based on given options
func ConnectWithOpts(
	gtx context.Context, opts *ConnOpts) (*sqlx.DB, error) {
	hostPort := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	if err := netx.WaitForPorts(gtx, hostPort, 60*time.Second); err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", opts.String())
	if err != nil {
		// log.Error().Err(err).Msg("failed to open postgress connection")
		return nil, errx.Errf(err, "failed to open postgress connection")
	}

	retry := 3
	for retry >= 0 {
		if err = db.Ping(); err != nil {
			retry--
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	var dbNow time.Time
	if err := db.Get(&dbNow, "SELECT now()"); err != nil {
		return nil, errx.Errf(err, "failed get current time")

	}
	log.Info().Str("DB.CurrentTime", dbNow.Format(time.UnixDate)).Msg("")

	return db, nil
}

// NamedConn - gives connection to database associated with given name. If no
// connection exists with given name nil is returned. If name is empty default
// connection is returned
func NamedConn(name string) *sqlx.DB {
	if name == "" {
		return defDB
	}
	return conns[name]
}

// SetNamedConn - register a postgres connection against name
func SetNamedConn(name string, db *sqlx.DB) {
	conns[name] = db
}

// SetDefaultConn - sets the default postgres connection
func SetDefaultConn(db *sqlx.DB) {
	defDB = db
}

// Conn - Gives default connection
func Conn() *sqlx.DB {
	return defDB
}

func Rollback(op string, tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Error().Err(err).Str("op", op).Msg("failed to rollback transaction")
	}
}

// func Commit(op string, tx *sqlx.Tx) {
// 	if err := tx.Commit(); err != nil {
// 		log.Error().Err(err).Str("op", op).Msg("failed to commit transaction")
// 	}
// }
