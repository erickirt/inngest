package state

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/inngest/inngestctl/inngest/client"
	"github.com/inngest/inngestctl/inngest/log"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	ErrNoState = fmt.Errorf("no Inngest state found")
)

const (
	SettingRanInit = "ranInit"
)

func SaveSetting(ctx context.Context, key string, value interface{}) error {
	s, _ := GetState(ctx)
	if s == nil {
		s = &State{Settings: make(map[string]interface{})}
	}
	s.Settings[key] = value
	return s.Persist(ctx)
}

func GetSetting(ctx context.Context, key string) interface{} {
	s, _ := GetState(ctx)
	if s == nil {
		return nil
	}
	setting, ok := s.Settings[key]
	if !ok {
		return nil
	}
	return setting
}

// State persists across each cli invokation, allowing functionality such as workspace
// switching, etc.
type State struct {
	client.Client `json:"-"`

	SelectedWorkspace *Workspace             `json:"workspace,omitempty"`
	Credentials       []byte                 `json:"credentials"`
	Account           client.Account         `json:"account"`
	Settings          map[string]interface{} `json:"settings"`
}

func (s State) Persist(ctx context.Context) error {
	path, err := homedir.Expand("~/.config/inngest")
	if err != nil {
		return fmt.Errorf("error reading ~/.config/inngest")
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("error creating ~/.config/inngest")
	}

	byt, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling state: %w", err)
	}

	path, err = homedir.Expand("~/.config/inngest/state")
	if err != nil {
		return fmt.Errorf("error reading ~/.config/inngest")
	}

	return ioutil.WriteFile(path, byt, 0600)
}

func (s *State) SetWorkspace(ctx context.Context, w client.Workspace) error {
	s.SelectedWorkspace = &Workspace{Workspace: w}
	return s.Persist(ctx)
}

// Workspace represents a single workspace within an Inngest account. The pertinent
// fields for the active workspace are marshalled into State.
type Workspace struct {
	client.Workspace

	IsOverridden bool `json:"-"`
}

// Client returns an API client, attempting to use authentication from
// state if found.
func Client(ctx context.Context) client.Client {
	state, _ := GetState(ctx)
	if state != nil {
		return state.Client
	}
	return client.New(client.WithAPI(viper.GetString("api")))
}

func AccountIdentifier(ctx context.Context) (string, error) {
	state, err := GetState(ctx)
	if err != nil {
		return "", err
	}

	// Add your account identifier locally, before finding action versions.
	if state.Account.Identifier.Domain == nil {
		return state.Account.Identifier.DSNPrefix, nil
	}

	return *state.Account.Identifier.Domain, nil
}

func GetState(ctx context.Context) (*State, error) {
	path, err := homedir.Expand("~/.config/inngest")
	if err != nil {
		return nil, fmt.Errorf("error reading ~/.config/inngest")
	}

	dir := os.DirFS(path)
	byt, err := fs.ReadFile(dir, "state")
	if errors.Is(err, fs.ErrNotExist) {
		return nil, ErrNoState
	}

	state := &State{}
	if err := json.Unmarshal(byt, state); err != nil {
		return nil, fmt.Errorf("invalid state file: %w", err)
	}

	// add the client using our stored credentials.
	state.Client = client.New(
		client.WithCredentials(state.Credentials),
		client.WithAPI(viper.GetString("api")), // "INNGEST_API", set up by commands/root
	)

	wid := viper.GetString("workspace.id")
	if wid == "" {
		return state, nil
	}

	id, err := uuid.Parse(wid)
	if err != nil {
		log.From(ctx).Warn().Err(err).Msg("invalid WORKSPACE_ID uuid")
		return state, nil
	}

	state.SelectedWorkspace = &Workspace{Workspace: client.Workspace{ID: id}, IsOverridden: true}
	return state, nil
}

func RequireState(ctx context.Context) *State {
	state, err := GetState(ctx)
	if err == ErrNoState {
		fmt.Println("\nRun `inngestctl login` and log in before running this command.")
		os.Exit(1)
	}

	if err != nil {
		log.From(ctx).Fatal().Msgf("error reading state: %s", err.Error())
	}

	return state
}