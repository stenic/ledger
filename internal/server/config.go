package server

type ServerOpts struct {
	StaticAssetPath string
	DisableAuth     bool
	EnableGZip      bool
	OidcIssuerURL   string
	OidcClientID    string
}
