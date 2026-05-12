package helpers

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func GetGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Warning: git rev-parse failed, defaulting to 'dev'")
		return "dev"
	}
	branch := strings.TrimSpace(out.String())
	log.Println("ENV run in: ", branch)
	return branch
}

func GetEnvFileFromBranch(branch string) string {
	switch {
	case branch == "masters" || branch == "main":
		log.Print("Env prod")
		return "app/config/.env.production"
	case strings.HasPrefix(branch, "staging"):
		log.Print("Env Staging")
		return "app/config/.env.staging"
	case strings.HasPrefix(branch, "development"):
		log.Print("Env Dev")
		return "app/config/.env.development"
	default:
		log.Print("Env local")
		return "app/config/.env.local"
	}
}
