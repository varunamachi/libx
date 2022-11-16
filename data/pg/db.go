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

//ConnOpts - postgres connection options
type ConnOpts struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbName"`
	TimeZone string `json:"timeZone"`
}

//String - get usable connection string
func (c *ConnOpts) String() string {
	return fmt.Sprintf(
		// postgres://username:password@url.com:5432/dbName
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&timezone=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.TimeZone,
	)
}

//Url - get a postgres URL
func (c *ConnOpts) Url() (*url.URL, error) {
	urlStr := fmt.Sprintf(
		// postgres://username:password@url.com:5432/dbName
		"postgres://%s:%s@%s:%d/%s?sslmode=disable&timezone=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.TimeZone,
	)
	return url.Parse(urlStr)
}

//connect - connects to DB based on connection string or URL
func Connect(gtx context.Context, url *url.URL) (*sqlx.DB, error) {
	if err := netx.WaitForPorts(gtx, url.Host, 60*time.Second); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	db, err := sqlx.Open("postgres", url.String())
	if err != nil {
		log.Error().Err(err).Msg("failed to open postgress connection")
	}
	return db, err
}

//ConnectWithOpts - connect to postgresdb based on given options
func ConnectWithOpts(
	gtx context.Context, opts *ConnOpts) (db *sqlx.DB, err error) {
	u, err := opts.Url()
	if err != nil {
		return nil, errx.Errf(err, "failed to create pg URL")
	}
	return Connect(gtx, u)
}

//NamedConn - gives connection to database associated with given name. If no
//connection exists with given name nil is returned. If name is empty default
//connection is returned
func NamedConn(name string) *sqlx.DB {
	if name == "" {
		return defDB
	}
	return conns[name]
}

//SetNamedConn - register a postgres connection against name
func SetNamedConn(name string, db *sqlx.DB) {
	conns[name] = db
}

//SetDefaultConn - sets the default postgres connection
func SetDefaultConn(db *sqlx.DB) {
	defDB = db
}

//Conn - Gives default connection
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
