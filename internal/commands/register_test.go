package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	product "github.com/neetozone/neeto-auth-cli"
	"github.com/neetozone/neeto-cli-commons/auth"
	"github.com/neetozone/neeto-cli-commons/cli"
	"github.com/neetozone/neeto-cli-commons/config"
)

type recordedRequest struct {
	Method string
	Path   string
	Body   map[string]any
}

func runCommand(t *testing.T, status int, response string, args ...string) (*recordedRequest, error) {
	t.Helper()

	recorded := &recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorded.Method = r.Method
		recorded.Path = r.URL.EscapedPath()
		if raw, err := io.ReadAll(r.Body); err == nil && len(raw) > 0 {
			_ = json.Unmarshal(raw, &recorded.Body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(server.Close)

	cfg, err := config.Parse(product.ConfigYAML)
	if err != nil {
		t.Fatalf("config.Parse: %v", err)
	}

	t.Setenv("HOME", t.TempDir())
	t.Setenv(cfg.BaseURLEnvVar(), server.URL)

	a := cli.New(*cfg)
	a.Printer.Out = io.Discard
	a.Printer.Err = io.Discard
	a.Root().SetOut(io.Discard)
	a.Root().SetErr(io.Discard)

	credentials := auth.Credentials{Subdomain: "acme", Email: "oliver.smith@example.com", SessionToken: "token"}
	if err := a.Auth.SaveStore(&auth.Store{Credentials: []auth.Credentials{credentials}}); err != nil {
		t.Fatalf("SaveStore: %v", err)
	}

	Register(a)
	a.Root().SetArgs(args)
	return recorded, a.Root().Execute()
}

func TestRegisterWiresProductCommandsOntoTheRoot(t *testing.T) {
	_, err := runCommand(t, http.StatusOK, `{"products":[],"pagination":{}}`, "products", "list")
	if err != nil {
		t.Fatalf("products list on the registered root: %v", err)
	}
}
