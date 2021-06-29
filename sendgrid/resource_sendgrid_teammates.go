package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sendgrid "github.com/trois-six/terraform-provider-sendgrid/sdk"
)

func resourceSendgridTeammates() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridTeammatesCreate,
		ReadContext:   resourceSendgridTeammatesRead,
		UpdateContext: resourceSendgridTeammatesUpdate,
		DeleteContext: resourceSendgridTeammatesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Description: "The email of the teammates.",
				Required:    true,
			},
			"is_admin": {
				Type:        schema.TypeBool,
				Description: "The email of the teammates.",
				Required:    true,
			},
			// "scopes": {
			// 	Type:        schema.TypeMap,
			// 	Description: "The IP addresses that should be assigned to this subuser.",
			// 	Elem:        &schema.Schema{Type: schema.TypeString},
			// 	Computed:    true,
			// },
		},
	}
}

func resourceSendgridTeammatesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	email := d.Get("email").(string)
	isAdmin := d.Get("is_admin").(bool)
	var scopes []string

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateTeammates(email, isAdmin, scopes)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: PendingID
	d.SetId(email)

	// if d.Get("disabled").(bool) {
	// 	return resourceSendgridSubuserUpdate(ctx, d, m)
	// }

	// 確認しようがないので、nilを返すべき
	// return resourceSendgridTeammatesRead(ctx, d, m)
	return nil
}

func resourceSendgridTeammatesRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	// hack to clear any on behalf of set in create sub user
	// to fix this properly I think we need to pass this down rather than setting global state
	c.OnBehalfOf = ""

	teammates, requestErr := c.ReadTeammates(d.Id())
	if requestErr.Err != nil {
		return diag.FromErr(requestErr.Err)
	}

	if len(teammates) == 0 {
		return diag.FromErr(teammatesNotFound(d.Id()))
	}
	//nolint:errcheck
	d.Set("email", teammates[0].Email)

	//nolint:errcheck
	d.Set("is_admin", teammates[0].IsAdmin)

	// TODO:
	// d.Set("scopes", teammates[0].Scopes)

	return nil
}

func resourceSendgridTeammatesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	if d.HasChange("disabled") {
		_, requestErr := c.UpdateTeammates(d.Id(), d.Get("disabled").(bool))
		if requestErr.Err != nil {
			return diag.FromErr(requestErr.Err)
		}
	}

	return resourceSendgridSubuserRead(ctx, d, m)
}

func resourceSendgridTeammatesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteTeammates(d.Id())
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
