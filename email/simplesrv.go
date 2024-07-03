package email

import (
	"context"
	"net/url"

	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/httpx"
	"github.com/varunamachi/libx/rt"
)

type SimpleServiceClient struct {
	sendUrl *url.URL
	// Add auth info when required
}

func (ssc *SimpleServiceClient) Send(md *Message, html bool) error {

	res := httpx.NewClient(ssc.sendUrl.Host, "").
		Build().
		Path(ssc.sendUrl.Path).
		QBool("html", html).
		Post(context.TODO(), html)
	if err := res.Close(); err != nil {
		return err
	}

	return nil
}

func NewSimpleMailSrvClinetFromEnv(envPrefix string) (Provider, error) {
	urlStr := rt.EnvString(envPrefix+"_SIMPLE_SRV_SEND_URL", "")
	sendUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, errx.Errf(err, "failed to parse send URL from ")
	}
	return &SimpleServiceClient{
		sendUrl: sendUrl,
	}, nil
}
