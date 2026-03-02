package helper

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GetGitHubOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		// ClientID didapat dari GitHub Developer Settings
		ClientID: os.Getenv("GITHUB_CLIENT_ID"),

		// ClientSecret didapat dari GitHub Developer Settings
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),

		// URL tujuan setelah user klik "Authorize" di GitHub
		// Harus SAMA PERSIS dengan yang didaftarkan di GitHub Console
		RedirectURL: "http://localhost:8080/api/v1/auth/github/callback",

		// Izin data apa saja yang ingin kita ambil dari GitHub
		Scopes: []string{"read:user", "user:email"},

		// Endpoint resmi GitHub untuk proses OAuth2
		Endpoint: github.Endpoint,
	}
}
