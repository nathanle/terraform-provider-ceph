package ceph

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyAccessKey = "access_key"
	keySecretKey = "secret_key"
	keyRegion    = "region"
	keyEndpoint  = "endpoint"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyAccessKey: {
				Type:        schema.TypeString,
				Description: "The access key to authenticate to S3 instance.\nCan also be provided by setting the S3_ACCESS_KEY environment variable.",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("S3_ACCESS_KEY", nil),
			},
			keySecretKey: {
				Type:        schema.TypeString,
				Description: "The secret key to authenticate to S3 instance.\nCan also be provided by setting the S3_SECRET environment variable.",
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("S3_SECRET", nil),
			},
			keyRegion: {
				Type:        schema.TypeString,
				Description: "The region.\nCan also be provided by setting the S3_REGION environment variable, if not specified it will default to: " + defaultRegion,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("S3_REGION", defaultRegion),
			},
			keyEndpoint: {
				Type:        schema.TypeString,
				Description: "The endpoint to use for requests.\nCan also be provided by setting the S3_ENDPOINT environment variable, if not specified it will default to: " + defaultEndpoint,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("S3_ENDPOINT", defaultEndpoint),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ceph_s3_bucket": resourceS3Bucket(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	accessKey := d.Get(keyAccessKey).(string)
	secret := d.Get(keySecretKey).(string)
	region := d.Get(keyRegion).(string)
	endpoint := d.Get(keyEndpoint).(string)

	var diags diag.Diagnostics

	if strings.TrimSpace(accessKey) == "" || strings.TrimSpace(secret) == "" {
		return nil, diags
	}
	if strings.TrimSpace(region) == "" {
		region = defaultRegion
	}
	if strings.TrimSpace(endpoint) == "" {
		endpoint = defaultEndpoint
	}
	return &Client{region: region, endpoint: endpoint, signer: v4.NewSigner(credentials.NewStaticCredentials(accessKey, secret, ""))}, diags
}
