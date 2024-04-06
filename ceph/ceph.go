package ceph

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type (
	Client struct {
		region               string
		endpoint             string
		virtualHostedBuckets bool
		insecure             bool
		signer               *v4.Signer
	}
	Request struct {
		method  string
		name    string
		headers map[string]string
		params  map[string]string
	}
)

const (
	defaultRegion   = "dc-sto1"
	defaultEndpoint = "objects.dc-sto1.glesys.net"
	defaultACL      = "private" // Can be one of: private, public-read, public-read-write, authenticated-read
)

func (i *Client) createS3Bucket(ctx context.Context, diags *diagnostics, name, acl string) error {
	var sc, tryCount = 0, 3
	var data []byte
	var err error
	for tryCount > 0 {
		tryCount--
		sc, data, err = send(ctx, i, Request{method: "PUT", name: name, headers: map[string]string{"x-amz-acl": acl}})
		if err != nil {
			return err
		}
		if sc == http.StatusOK {
			return nil
		}
		if sc == http.StatusConflict && strings.Contains(string(data), "<Code>BucketAlreadyExists</Code>") {
			return fmt.Errorf("createBucket: '%s' bucket already exists under different user’s ownership", name)
		}

		if tryCount > 0 {
			// Sleep before next try. If the S3 instance were just created we often get an HTTP 403 and an error message with
			// "<Code>InvalidAccessKeyId</Code>" the first time we try to create a bucket in the new S3 instance.
			time.Sleep(2 * time.Second)
		}
	}
	if sc > 299 {
		return fmt.Errorf("createBucket: failed to create bucket, status code: %d, response: %s", sc, data)
	}
	return nil
}

func (i *Client) setS3BucketPolicy(ctx context.Context, diags *diagnostics, name, acl string) error {
	var sc, tryCount = 0, 3
	var data []byte
	var err error
	for tryCount > 0 {
		tryCount--
		sc, data, err = send(ctx, i, Request{method: "PUT", name: name, headers: map[string]string{"x-amz-acl": acl}})
		if err != nil {
			return err
		}
		if sc == http.StatusOK {
			return nil
		}
		//if sc == http.StatusConflict && strings.Contains(string(data), "<Code>BucketAlreadyExists</Code>") {
		//	return fmt.Errorf("createBucket: '%s' bucket already exists under different user’s ownership", name)
		//}

		//if tryCount > 0 {
			//// Sleep before next try. If the S3 instance were just created we often get an HTTP 403 and an error message with
			//// "<Code>InvalidAccessKeyId</Code>" the first time we try to create a bucket in the new S3 instance.
			//time.Sleep(2 * time.Second)
		//}
	}
	if sc > 299 {
		return fmt.Errorf("createBucket: failed to create bucket, status code: %d, response: %s", sc, data)
	}
	return nil
}

func (i *Client) deleteS3Bucket(ctx context.Context, name string) error {
	sc, data, err := send(ctx, i, Request{method: "DELETE", name: name})
	if err != nil {
		return err
	}
	if sc == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("createBucket: failed to delete bucket, status code: %d, received data: %s", sc, data)
}

func (i *Client) s3BucketExist(ctx context.Context, name string) error {
	// Query bucket ACL to test if bucket exist
	sc, data, err := send(ctx, i, Request{method: "GET", name: name, params: map[string]string{"acl": ""}})
	if err != nil {
		return err
	}
	if sc == http.StatusOK {
		return nil
	}
	err = fmt.Errorf("bucketExist: failed to check bucket, status code: %d, received data: %s", sc, data)
	return err
}

func send(ctx context.Context, client *Client, request Request) (int, []byte, error) {
	var proto = "https"
	if client.insecure {
		proto = "http"
	}

	var uri string
	if client.virtualHostedBuckets {
		uri = fmt.Sprintf("%s://%s.%s", proto, request.name, client.endpoint)
	} else {
		uri = fmt.Sprintf("%s://%s/%s", proto, client.endpoint, request.name)
	}

	u, err := url.Parse(uri)
	if err != nil {
		return 0, nil, fmt.Errorf("send: failed to parse URL: %s, %w", uri, err)
	}

	qp := u.Query()
	for k, v := range request.params {
		qp.Set(k, v)
	}
	u.RawQuery = qp.Encode()
	req, err := http.NewRequestWithContext(ctx, request.method, u.String(), nil)
	if err != nil {
		return 0, nil, fmt.Errorf("send: failed to create request, %w", err)
	}

	for k, v := range request.headers {
		req.Header.Set(k, v)
	}

	_, err = client.signer.Sign(req, nil, "s3", client.region, time.Now())
	if err != nil {
		return 0, nil, fmt.Errorf("send: failed to sign request, %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("send: error while sending request, %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, fmt.Errorf("send: error while reading response, %w", err)
	}

	return res.StatusCode, data, nil
}
