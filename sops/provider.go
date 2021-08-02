package sops

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/lokkersp/terraform-provider-sops/sops/internal/sops"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kms": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: providerDescriptions["kms"],
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"age": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: providerDescriptions["aws"],
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sops_file":     dataSourceFile(),
			"sops_external": dataSourceExternal(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"sops_file": resourceSourceFile(),
		},
		ConfigureContextFunc: ConfigureProvider,
	}
}

var providerDescriptions = map[string]string{
	"kms": "Configuration for encrypt files with AWS KMS.",
	"age": "Configuration for encrypt files with Age.",
}

type EncryptConfig struct {
	Kms sops.KmsConf
	Age string
}

func ConfigureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	encConf := &EncryptConfig{}
	kms, err := sops.GetKmsConf(d)
	if err != nil {
		fmt.Println("failed to init kms")
	} else {
		encConf.Kms = kms
	}
	age, err := sops.GetAgeConf(d)
	if err != nil {
		fmt.Println("failed to init age")
	} else {
		encConf.Age = age
	}

	return encConf, diags
}
