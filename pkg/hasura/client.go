package hasura

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type (
	GqlClient struct {
		query
		log *logrus.Entry
	}

	defaultQuery struct {
		endpoint string
		log      *logrus.Entry
	}

	query interface {
		doGraphql(ctx context.Context, query string, variables map[string]interface{}, token string) (io.ReadCloser, error)
	}
)

// doGraphql will fetch the query result from Hasura
func (c *defaultQuery) doGraphql(ctx context.Context, query string, variables map[string]interface{}, token string) (io.ReadCloser, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(`{"query":"`)
	buf.WriteString(jsonEscape(query))
	buf.WriteString(`"`)
	if variables != nil {
		buf.WriteString(`,"variables":`)
		b, err := json.Marshal(variables)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal variables: %w")
		}
		buf.Write(b)
	}
	buf.WriteString("}")
	c.log.Infof("%s", buf.String())
	req, err := http.NewRequest(http.MethodPost, c.endpoint, buf)
	if err != nil {
		return nil, err
	}
	// req.Header.Add("authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body, errors.New(fmt.Sprintf("%d %s", resp.StatusCode, "request failed with code"))
	}
	return resp.Body, nil
}

func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

// HasuraClient set the client for hasura
func HasuraClient(hasuraEP string, log *logrus.Entry) (*GqlClient, error) {
	return &GqlClient{
		query: &defaultQuery{endpoint: hasuraEP, log: log},
		log:   log,
	}, nil
}
