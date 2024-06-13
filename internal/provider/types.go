package provider

type RunQueue struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	EntityName            string        `json:"entityName"`
	PrioritizationMode    string        `json:"prioritizationMode"`
	ExternalLinks         ExternalLinks `json:"externalLinks"`
	CreatedAt             string        `json:"createdAt"`
	UpdatedAt             string        `json:"updatedAt"`
	DefaultResourceConfig struct {
		ID                string                 `json:"id"`
		Resource          string                 `json:"resource"`
		Config            map[string]interface{} `json:"config"`
		TemplateVariables []TemplateVariableWithName
	} `json:"defaultResourceConfig"`
}

type ExternalLink struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type ExternalLinks struct {
	Links []ExternalLink `json:"links"`
}

type TemplateVariableWithName struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Schema      string  `json:"schema"`
}

type TemplateVariable struct {
	Description *string `json:"description,omitempty"`
	Schema      string  `json:"schema"`
}

type UpsertRunQueueInput struct {
	QueueName          string  `json:"queueName"`
	EntityName         string  `json:"entityName"`
	ProjectName        string  `json:"projectName"`
	ResourceType       string  `json:"resourceType"`
	ResourceConfig     string  `json:"resourceConfig"`
	TemplateVariables  *string `json:"templateVariables"`
	PrioritizationMode *string `json:"prioritizationMode"`
	ExternalLinks      *string `json:"externalLinks"`
}
