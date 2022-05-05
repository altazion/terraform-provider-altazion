package altazion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceModules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceModulesRead,
		Schema: map[string]*schema.Schema{
			"modules": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceModulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/GetModules", "http://localhost:19090"), nil)
	req.Header = map[string][]string{
		"Accept-Encoding": {"gzip, deflate"},
	}

	if err != nil {
		return diag.FromErr(err)
	}

	r, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer r.Body.Close()

	coffees := make([]map[string]interface{}, 0)
	err = json.NewDecoder(r.Body).Decode(&coffees)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("coffees", coffees); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
