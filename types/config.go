package types

type DefaultConfig struct {
	AdminPod     string   `yaml:"adminpod"`     // Where to run privilege code
	AdminVolumes []string `yaml:"adminvolumes"` // Volume used by admin
	Pod          string   `yaml:"pod"`          // Default Pod
	Image        string   `yaml:"image"`        // Base Image
	Network      string   `yaml:"network"`      // Default SDN
	WorkingDir   string   `yaml:"working_dir"`  // Default CWD
	Cpu          float64  `yaml:"cpu"`          // Default CPU
	Memory       int64    `yaml:"memory"`       // Default Memory
	Timeout      int      `yaml:"timeout"`      // Default timeout
	OpenStdin    bool     `yaml:"OpenStdin"`    // Default Openstdin
}

// Config holds eru-lambda config
type Config struct {
	Servers     []string      `yaml:"servers"` // Cores' address and port
	Default     DefaultConfig `yaml:"default"`
	Concurrency int           `yaml:concurrency"` // Default concurrency
}

type RunParams struct {
	Name       string
	Command    string
	Network    string
	Workingdir string
	Image      string
	CPU        float64
	Mem        int64
	Count      int
	Timeout    int
	Envs       []string
	Volumes    []string
	OpenStdin  bool
	Pod        string
}
