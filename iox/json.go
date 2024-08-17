package iox

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/varunamachi/libx/errx"
)

// PrintJSON - dumps JSON representation of given data to stdout
func PrintJSON(o interface{}) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		log.Error().Err(err).Msg("failed to write formatted JSON to console")
		return
	}

	fmt.Println(string(b))
}

// WriteJSON - writes JSON representation of given data to given writer
func WriteJSON(writer io.Writer, o interface{}) error {
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return errx.Wrap(err)
	}

	fmt.Fprintln(writer, string(b))
	return nil
}

func WriteJSONFile(
	path string, conflictPolicy FileConflictPolicy, data interface{}) error {
	file, err := CreateFile(path, conflictPolicy)
	if err != nil {
		log.Error().Err(err).Msg("failed to create JSON file")
		return errx.Wrap(err)
	}
	defer file.Close()

	if err := WriteJSON(file, data); err != nil {
		log.Error().Err(err).Msg("failed to write JSON content to file")
		return errx.Wrap(err)
	}

	return nil
}

func LoadJsonFile(path string, out interface{}) error {
	reader, err := os.Open(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).
			Msg("failed to open JSON file")
		return errx.Errf(err, "failed to open JSON file at %s", path)
	}
	if err = LoadJson(reader, out); err != nil {
		return errx.Wrap(err)
	}
	return nil
}

func LoadJson(reader io.Reader, out interface{}) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		const msg = "Failed to read from reader"
		log.Error().Err(err).Msg(msg)
		return errx.Errf(err, msg)
	}

	if err = json.Unmarshal(data, out); err != nil {
		const msg = "Failed to decode JSON data"
		log.Error().Err(err).Msg(msg)
		return errx.Errf(err, "Failed to decode JSON data")
	}
	return nil
}

// FormattedJSON - converts given data to JSON and returns as pretty printed
func FormattedJSON(o interface{}) (string, error) {
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return "", errx.Errf(err, "failed to generate formatted JSON")
	}
	return string(b), nil
}
