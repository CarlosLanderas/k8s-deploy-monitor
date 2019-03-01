package deployment

// ClusterConfig is the configuration used to connect a remote kubernetes cluster
type ClusterConfig struct {
	Host string
	VerifyTlsCertificate bool
	BearerToken string
}
