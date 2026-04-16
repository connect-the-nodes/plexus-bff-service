package model

type ConnectorsInfo struct {
	Status  string
	Message string
}

type FeatureFlag struct {
	Name        string `json:"name" yaml:"name"`
	Parent      string `json:"parent" yaml:"parent"`
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Description string `json:"description" yaml:"description"`
}
