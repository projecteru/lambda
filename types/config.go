package types

// Config holds eru-lambda config
type Config struct {
	Servers    []string `yaml:"servers"`    // Cores' address and port
	AdminPod   string   `yaml:"adminpod"`   // Where to run privilege code
	DefaultSDN string   `yaml:"defaultsdn"` // Default SDN
}
