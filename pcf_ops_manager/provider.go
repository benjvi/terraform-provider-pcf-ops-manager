package pcf_ops_manager

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type OpsmanClient struct {
	target, token, clientId, clientSecret, username, password string
	skipSslValidation bool
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"target_hostname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"username", "password", "client_id", "client_secret"},
				Description: "Use generated token from UAA in lieu of normal auth, see: https://pcf.pcf-aws.bjv.me/docs#authentication",
			},

			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"username", "password", "token"},
			},

			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"username", "password", "token"},
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"client_id", "client_secret", "token"},
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{"client_id", "client_secret", "token"},
			},

			"skip_ssl_validation": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:	 false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"pcfom_director": resourcePcfDirector(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	om := &OpsmanClient{
		target: d.Get("target_hostname").(string),
		token: d.Get("token").(string),
		skipSslValidation: d.Get("skip_ssl_validation").(bool),
		//TODO
	}
	return om, nil
}

