package app

import (
	"encoding/json"
	"fmt"
	"os"
)

type configShowOutput struct {
	ConfigPath      string                    `json:"configPath"`
	SecretStorePath string                    `json:"secretStorePath"`
	CurrentProfile  string                    `json:"currentProfile"`
	SelectedProfile string                    `json:"selectedProfile"`
	Profile         map[string]any            `json:"profile,omitempty"`
	Profiles        map[string]map[string]any `json:"profiles,omitempty"`
}

func runConfigShow(cfg Config) error {
	store, err := loadProfileStore(cfg.ConfigPath)
	if err != nil {
		return err
	}

	selectedProfile := cfg.Profile
	if selectedProfile == "" {
		if store.CurrentProfile != "" {
			selectedProfile = store.CurrentProfile
		} else {
			selectedProfile = "default"
		}
	}

	out := configShowOutput{
		ConfigPath:      cfg.ConfigPath,
		SecretStorePath: cfg.SecretStorePath,
		CurrentProfile:  store.CurrentProfile,
		SelectedProfile: selectedProfile,
		Profiles:        map[string]map[string]any{},
	}

	for name, profile := range store.Profiles {
		out.Profiles[name] = profileToMap(name, cfg.SecretStorePath, profile)
	}
	if profile, ok := store.Profiles[selectedProfile]; ok {
		out.Profile = profileToMap(selectedProfile, cfg.SecretStorePath, profile)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(out); err != nil {
		return fmt.Errorf("encoding config output: %w", err)
	}
	return nil
}

func profileToMap(profileName, secretStorePath string, profile SavedProfile) map[string]any {
	out := map[string]any{
		"source_base_url":       profile.SourceBaseURL,
		"target_base_url":       profile.TargetBaseURL,
		"identity_mapping_file": profile.IdentityMappingFile,
		"teams_file":            profile.TeamsFile,
		"persons_file":          profile.PersonsFile,
		"resources_file":        profile.ResourcesFile,
		"issues_csv":            profile.IssuesCSV,
		"output_dir":            profile.OutputDir,
		"report_format":         profile.ReportFormat,
		"saved_secrets":         secretStoreHasProfile(secretStorePath, profileName),
	}
	return out
}
