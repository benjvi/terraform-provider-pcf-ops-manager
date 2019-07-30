package pcf_ops_manager

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func resourcePcfDirector() *schema.Resource {
	return &schema.Resource{
		Create: resourcePcfDirectorCreate,
		Read:   resourcePcfDirectorRead,
		Update: resourcePcfDirectorUpdate,
		Delete: resourcePcfDirectorDelete,
		Importer:  &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"director_config": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
				Description: "JSON config file for director, IaaS, and security properties",
			},

			"force_delete": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
				Description: "(not implemented) - Continue with deletion even when there products deployed (this will destroy everything!)",
			},
		},
	}
}

func resourcePcfDirectorCreate(d *schema.ResourceData, m interface{}) error {

	return resourcePcfDirectorUpdate(d, m)
}

func resourcePcfDirectorRead(d *schema.ResourceData, m interface{}) error {
	opsmanClient := m.(*OpsmanClient)
	c := cleanhttp.DefaultClient()
	if opsmanClient.skipSslValidation {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Transport = logging.NewTransport("pcf_ops_manager", tr)
	} else {
		c.Transport = logging.NewTransport("pcf_ops_manager", c.Transport)
	}
	req, _ := http.NewRequest("GET","https://"+opsmanClient.target+"/api/v0/staged/director/properties", nil)
	req.Header["Authorization"] = []string{"Bearer "+opsmanClient.token}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Error getting staged director properties from opsman %q: %q", opsmanClient.target, err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if body == nil || string(body) == "" {
		return fmt.Errorf("Error getting staged director properties from opsman %q: %+v", opsmanClient.target, resp)
	}
	d.SetId(opsmanClient.target)
	d.Set("director_config",string(body))
	return nil
}

func resourcePcfDirectorUpdate(d *schema.ResourceData, m interface{}) error {
	opsmanClient := m.(*OpsmanClient)
	c := cleanhttp.DefaultClient()
	if opsmanClient.skipSslValidation {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.Transport = logging.NewTransport("pcf_ops_manager", tr)
	} else {
		c.Transport = logging.NewTransport("pcf_ops_manager", c.Transport)
	}
	directorConfigReader := strings.NewReader(d.Get("director_config").(string))
	req, _ := http.NewRequest("PUT","https://"+opsmanClient.target+"/api/v0/staged/director/properties", directorConfigReader)
	req.Header["Authorization"] = []string{"Bearer "+opsmanClient.token}
	req.Header["Content-Type"] = []string{"application/json"}
	_, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Error setting staged director properties on opsman %q: %q", opsmanClient.target, err.Error())
	}

	applyDirectorChangesConfig := "{\"deploy_products\": \"none\"}"
	applyChangesConfigReader := strings.NewReader(applyDirectorChangesConfig)
	req, _ = http.NewRequest("POST","https://"+opsmanClient.target+"/api/v0/installations", applyChangesConfigReader)
	req.Header["Authorization"] = []string{"Bearer "+opsmanClient.token}
	req.Header["Content-Type"] = []string{"application/json"}
	_, err = c.Do(req)
	if err != nil {
		return fmt.Errorf("Error applying director changes %q: %q", opsmanClient.target, err.Error())
	}

	//TODO wait for apply changes to finish
	return resourcePcfDirectorRead(d, m)
}

func resourcePcfDirectorDelete(d *schema.ResourceData, m interface{}) error {
	//TODO
	// delete-installation is dangerous but we should at least provide the option for people to do it
	return nil
}

func validateConfigJson(configI interface{}, k string) ([]string, []error) {
	dataJSON := configI.(string)
	dataMap := map[string]interface{}{}
	err := json.Unmarshal([]byte(dataJSON), &dataMap)
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
}


func suppressEquivalentJsonDiffs(k, old, new string, d *schema.ResourceData) bool {
	ob := bytes.NewBufferString("")
	if err := json.Compact(ob, []byte(old)); err != nil {
		return false
	}

	nb := bytes.NewBufferString("")
	if err := json.Compact(nb, []byte(new)); err != nil {
		return false
	}

	return jsonBytesEqual(ob.Bytes(), nb.Bytes())
}

func jsonBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}