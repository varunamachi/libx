package mg

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// ReadAllAndClose - reads all data from the cursor and closes it
func ReadAllAndClose(
	gtx context.Context,
	cur *mongo.Cursor,
	out interface{}) error {

	defer func() {
		if err := cur.Close(gtx); err != nil {
			log.Error().Err(err).Msg("failed to close cursor")
		}
	}()
	return cur.All(gtx, out)
}

// ReadOneAndClose - reads one data item from the cursor and closes it
func ReadOneAndClose(
	gtx context.Context,
	cur *mongo.Cursor,
	out interface{}) error {

	defer func() {
		if err := cur.Close(gtx); err != nil {
			log.Error().Err(err).Msg("failed to close cursor")
		}
	}()
	return cur.Decode(out)
}
