package types

type DefaultConfig struct {
	AdminPod     string   `yaml:"adminpod"`     // Where to run privilege code
	AdminVolumes []string `yaml:"adminvolumes"` // Volume used by admin
	Pod          string   `yaml:"pod"`          // Default Pod
	Image        string   `yaml:"image"`        // Base Image
	Network      string   `yaml:"network"`      // Default SDN
	Cpu          float64  `yaml:"cpu"`          // Default CPU
	Memory       int64    `yaml:"memory"`       // Default Memory
	Timeout      int      `yaml:"timeout"`      // Default timeout
}

// Config holds eru-lambda config
type Config struct {
	Servers     []string      `yaml:"servers"` // Cores' address and port
	Default     DefaultConfig `yaml:"default"`
	Concurrency int           `yaml:concurrency"` // Default concurrency
}
