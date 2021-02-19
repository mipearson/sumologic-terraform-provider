package sumologic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSumologicCloudsyslogSource() *schema.Resource {
	cloudSyslogSource := resourceSumologicSource()
	cloudSyslogSource.Create = resourceSumologicCloudSyslogSourceCreate
	cloudSyslogSource.Read = resourceSumologicCloudSyslogSourceRead
	cloudSyslogSource.Update = resourceSumologicCloudSyslogSourceUpdate

	cloudSyslogSource.Schema["token"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return cloudSyslogSource
}

func resourceSumologicCloudSyslogSourceCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("lookup_by_name").(bool) {
		source, err := c.GetSourceName(d.Get("collector_id").(int), d.Get("name").(string))

		if err != nil {
			return err
		}

		if source != nil {
			d.SetId(strconv.Itoa(source.ID))
		}
	}

	if d.Id() == "" {
		source := resourceToCloudSyslogSource(d)

		id, err := c.CreateCloudsyslogSource(source, d.Get("collector_id").(int))

		if err != nil {
			return err
		}

		d.SetId(strconv.Itoa(id))
	}

	return resourceSumologicCloudSyslogSourceRead(d, meta)
}

func resourceSumologicCloudSyslogSourceUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	source := resourceToCloudSyslogSource(d)

	err := c.UpdateCloudSyslogSource(source, d.Get("collector_id").(int))

	if err != nil {
		return err
	}

	return resourceSumologicCloudSyslogSourceRead(d, meta)
}

func resourceToCloudSyslogSource(d *schema.ResourceData) CloudSyslogSource {
	source := resourceToSource(d)
	source.Type = "Cloudsyslog"

	cloudsyslogSource := CloudSyslogSource{
		Source: source,
	}

	return cloudsyslogSource
}

func resourceSumologicCloudSyslogSourceRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id, _ := strconv.Atoi(d.Id())
	source, err := c.GetCloudSyslogSource(d.Get("collector_id").(int), id)
	if err != nil {
		return err
	}

	if source == nil {
		log.Printf("[WARN] Cloud Syslog source not found, removing from state: %v - %v", id, err)
		d.SetId("")

		return nil
	}

	resourceSumologicSourceRead(d, source.Source)
	d.Set("token", source.Token)

	return nil
}
