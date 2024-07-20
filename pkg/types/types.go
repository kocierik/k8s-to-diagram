package types

type K8sResource struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string            `yaml:"name"`
		Annotations map[string]string `yaml:"annotations"`
	} `yaml:"metadata"`
}

type Communication struct {
	Name     string                `json:"name"`
	Inbound  []CommunicationDetail `json:"inbound"`
	Outbound []CommunicationDetail `json:"outbound"`
}

type CommunicationDetail struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
}
