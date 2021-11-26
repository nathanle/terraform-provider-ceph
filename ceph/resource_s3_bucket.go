package ceph

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	diagnostics struct {
		diag.Diagnostics
	}
)

const (
	keyBucketName    = "name"
	keyBucketCreated = "created"
)

func (d *diagnostics) warn(msg string) {
	d.Diagnostics = append(d.Diagnostics, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  msg,
	})
}

func (d *diagnostics) warnf(format string, a ...interface{}) {
	d.Diagnostics = append(d.Diagnostics, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf(format, a...),
	})
}

func (d *diagnostics) error(err error) {
	d.Diagnostics = append(d.Diagnostics, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  err.Error(),
	})
}

func resourceS3Bucket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceS3BucketCreate,
		ReadContext:   resourceS3BucketRead,
		UpdateContext: resourceS3BucketUpdate,
		DeleteContext: resourceS3BucketDelete,
		Schema:        resourceS3BucketSchema(),
	}
}

func resourceS3BucketSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		keyBucketName: {
			Type:     schema.TypeString,
			Required: true,
		},
		// Do we need to set ACL on bucket, currently the buckets are created with "x-amz-acl" = "private", but if the
		// bucket already exist for "our" user, we don't check what ACL it have.
		//"acl": {
		//	Type:     schema.TypeString,
		//	Optional: true,
		//},
		keyBucketCreated: {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceS3BucketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diagnostics
	c := m.(*Client)
	bucketName := d.Get(keyBucketName).(string)

	err := c.createS3Bucket(ctx, &diags, bucketName, defaultACL)
	if err != nil {
		diags.error(err)
		return diags.Diagnostics
	}

	d.SetId(bucketName)
	d.Set(keyBucketCreated, time.Now().Format(time.RFC3339))
	return diags.Diagnostics
}

func resourceS3BucketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Client)
	bucketName := d.Get(keyBucketName).(string)

	// Just test if the bucket exist, for now.
	err := c.s3BucketExist(ctx, bucketName)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceS3BucketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Not really anything that we support currently since we only can/want to either create or destroy buckets right now.
	var diags diagnostics
	diags.warn("bucket update not supported")
	return diags.Diagnostics
}

func resourceS3BucketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diagnostics
	c := m.(*Client)
	bucketName := d.Get(keyBucketName).(string)

	err := c.deleteS3Bucket(ctx, bucketName)
	if err != nil {
		diags.error(err)
	}
	return diags.Diagnostics
}
