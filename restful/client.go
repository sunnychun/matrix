package restful

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/ironzhang/matrix/codes"
	"github.com/ironzhang/matrix/context-value"
	"github.com/ironzhang/matrix/restful/codec"
	"github.com/ironzhang/matrix/tlog"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	Client HTTPClient
	Codec  codec.Codec
	Writer io.Writer
}

func (c *Client) DoContext(ctx context.Context, method, url string, args, reply interface{}) (err error) {
	log := tlog.WithContext(ctx).Sugar().With("method", method, "url", url)
	verbose := context_value.ParseVerbose(ctx)

	var b bytes.Buffer

	// Encode
	if args != nil {
		if err = c.codec().Encode(&b, args); err != nil {
			log.Errorw("encode", "error", err)
			return err
		}
	}

	// New http request
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		log.Errorw("new request", "error", err)
		return err
	}
	setHeaderWithContext(ctx, req.Header)

	// Print request
	if verbose {
		c.printRequest(ctx, req)
	}

	// Do
	resp, err := c.client().Do(req)
	if err != nil {
		log.Errorw("client do", "error", err)
		return err
	}
	defer resp.Body.Close()

	// Print response
	if verbose {
		c.printResponse(ctx, resp)
	}

	// Handle error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e codec.Error
		if err = c.codec().DecodeError(resp.Body, &e); err != nil {
			log.Errorw("decode error", "error", err, "status", http.StatusText(resp.StatusCode))
			return err
		}
		return Errorf(resp.StatusCode, codes.Code(e.Code), e.Cause)
	}

	// Decode
	if reply != nil {
		if err = c.codec().Decode(resp.Body, reply); err != nil {
			log.Errorw("decode", "error", err, "status", http.StatusText(resp.StatusCode))
			return err
		}
	}

	return nil
}

func (c *Client) client() HTTPClient {
	if c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

func (c *Client) codec() codec.Codec {
	if c.Codec == nil {
		return codec.DefaultCodec
	}
	return c.Codec
}

func (c *Client) writer() io.Writer {
	if c.Writer == nil {
		return os.Stdout
	}
	return c.Writer
}

func (c *Client) printRequest(ctx context.Context, r *http.Request) {
	b, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		tlog.WithContext(ctx).Sugar().Errorw("dump request out", "error", err)
		return
	}
	traceId := context_value.ParseTraceId(ctx)
	fmt.Fprintf(c.writer(), "traceId(%s) client request:\n%s\n", traceId, b)
}

func (c *Client) printResponse(ctx context.Context, r *http.Response) {
	b, err := httputil.DumpResponse(r, true)
	if err != nil {
		tlog.WithContext(ctx).Sugar().Errorw("dump response", "error", err)
		return
	}
	traceId := context_value.ParseTraceId(ctx)
	fmt.Fprintf(c.writer(), "traceId(%s) client response:\n%s\n", traceId, b)
}

func setHeaderWithContext(ctx context.Context, h http.Header) {
	if v := context_value.ParseTraceId(ctx); v != "" {
		h.Set(xTraceId, v)
	}
	if v := context_value.ParseVerbose(ctx); v {
		h.Set(xVerbose, "1")
	}
}