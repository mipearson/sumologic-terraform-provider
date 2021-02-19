package sumologic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type Source struct {
	ID                         int                 `json:"id,omitempty"`
	Type                       string              `json:"sourceType"`
	Name                       string              `json:"name"`
	Description                string              `json:"description,omitempty"`
	Category                   string              `json:"category,omitempty"`
	HostName                   string              `json:"hostName,omitempty"`
	TimeZone                   string              `json:"timeZone,omitempty"`
	AutomaticDateParsing       bool                `json:"automaticDateParsing"`
	MultilineProcessingEnabled bool                `json:"multilineProcessingEnabled"`
	UseAutolineMatching        bool                `json:"useAutolineMatching"`
	ManualPrefixRegexp         string              `json:"manualPrefixRegexp,omitempty"`
	ForceTimeZone              bool                `json:"forceTimeZone"`
	DefaultDateFormats         []DefaultDateFormat `json:"defaultDateFormats,omitempty"`
	Filters                    []Filter            `json:"filters,omitempty"`
	CutoffTimestamp            int                 `json:"cutoffTimestamp,omitempty"`
	CutoffRelativeTime         string              `json:"cutoffRelativeTime,omitempty"`
}

type DefaultDateFormat struct {
	Format  string `json:"format"`
	Locator string `json:"locator"`
}

type Filter struct {
	Name       string `json:"name"`
	FilterType string `json:"filterType"`
	Regexp     string `json:"regexp"`
	Mask       string `json:"mask"`
}

type SourceList struct {
	Sources []Source `json:"sources"`
}

func resourceSumologicSource() *schema.Resource {
	return &schema.Resource{
		Delete: resourceSumologicSourceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "",
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "Etc/UTC",
			},
			"automatic_date_parsing": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"multiline_processing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"use_autoline_matching": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
			"manual_prefix_regexp": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  nil,
			},
			"force_timezone": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"default_date_formats": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:     schema.TypeString,
							Required: true,
						},
						"locator": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filter_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Exclude", "Include", "Hash", "Mask", "Forward"}, false),
						},
						"regexp": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mask": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cutoff_timestamp": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  0,
			},
			"cutoff_relative_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  nil,
			},
			"collector_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"lookup_by_name": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},
			"destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},
		},
	}
}

func resourceSumologicSourceDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Get("destroy").(bool) {
		id, _ := strconv.Atoi(d.Id())
		collectorID, _ := d.Get("collector_id").(int)

		return c.DestroySource(id, collectorID)
	}

	return nil
}

func resourceSumologicSourceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), "/")

	if len(ids) != 2 {
		return nil, fmt.Errorf("expected collector_id/source_id, got %s", d.Id())
	}

	d.SetId(ids[1])

	collectorID, _ := strconv.Atoi(ids[0])
	d.Set("collector_id", collectorID)

	return []*schema.ResourceData{d}, nil
}

func resourceToSource(d *schema.ResourceData) Source {
	id, _ := strconv.Atoi(d.Id())

	source := Source{}
	source.ID = id
	source.Name = d.Get("name").(string)
	source.Description = d.Get("description").(string)
	source.Category = d.Get("category").(string)
	source.HostName = d.Get("host_name").(string)
	source.TimeZone = d.Get("timezone").(string)
	source.AutomaticDateParsing = d.Get("automatic_date_parsing").(bool)
	source.MultilineProcessingEnabled = d.Get("multiline_processing_enabled").(bool)
	source.UseAutolineMatching = d.Get("use_autoline_matching").(bool)
	source.ManualPrefixRegexp = d.Get("manual_prefix_regexp").(string)
	source.ForceTimeZone = d.Get("force_timezone").(bool)
	source.DefaultDateFormats = getDefaultDateFormats(d)
	source.Filters = getFilters(d)
	source.CutoffTimestamp = d.Get("cutoff_timestamp").(int)
	source.CutoffRelativeTime = d.Get("cutoff_relative_time").(string)

	return source
}

func resourceSumologicSourceRead(d *schema.ResourceData, source Source) {
	d.Set("name", source.Name)
	d.Set("description", source.Description)
	d.Set("category", source.Category)
	d.Set("host_name", source.HostName)
	d.Set("timezone", source.TimeZone)
	d.Set("automatic_date_parsing", source.AutomaticDateParsing)
	d.Set("multiline_processing_enabled", source.MultilineProcessingEnabled)
	d.Set("use_autoline_matching", source.UseAutolineMatching)
	d.Set("manual_prefix_regexp", source.ManualPrefixRegexp)
	d.Set("force_timezone", source.ForceTimeZone)
	d.Set("default_date_formats", source.DefaultDateFormats)
	d.Set("filters", source.Filters)
	d.Set("cutoff_timestamp", source.CutoffTimestamp)
	d.Set("cutoff_relative_time", source.CutoffRelativeTime)
}

func getDefaultDateFormats(d *schema.ResourceData) []DefaultDateFormat {

	rawDefaultDateFormatsConfig := d.Get("default_date_formats").([]interface{})
	var defaultDateDormats []DefaultDateFormat

	for _, rawConfig := range rawDefaultDateFormatsConfig {
		config := rawConfig.(map[string]interface{})
		defaultDateFormat := DefaultDateFormat{}
		defaultDateFormat.Format = config["format"].(string)
		defaultDateFormat.Locator = config["locator"].(string)
		defaultDateDormats = append(defaultDateDormats, defaultDateFormat)
	}

	return defaultDateDormats
}

func getFilters(d *schema.ResourceData) []Filter {

	rawFilterConfig := d.Get("filters").([]interface{})
	var filters []Filter

	for _, rawConfig := range rawFilterConfig {
		config := rawConfig.(map[string]interface{})
		filter := Filter{}
		filter.Name = config["name"].(string)
		filter.FilterType = config["filter_type"].(string)
		filter.Regexp = config["regexp"].(string)
		filter.Mask = config["mask"].(string)
		filters = append(filters, filter)
	}

	return filters
}

func (s *Client) DestroySource(sourceID int, collectorID int) error {

	_, err := s.Delete(fmt.Sprintf("collectors/%d/sources/%d", collectorID, sourceID))

	return err
}

func (s *Client) GetSourceName(collectorID int, sourceName string) (*Source, error) {

	data, _, err := s.Get(
		fmt.Sprintf("collectors/%d/sources", collectorID),
	)

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var response SourceList
	err = json.Unmarshal(data, &response)

	if err != nil {
		return nil, err
	}

	for _, source := range response.Sources {
		if source.Name == sourceName {
			return &source, nil
		}
	}

	return nil, nil
}
