package login

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var ErrNoStoredToken = errors.New("no stored token")

type TokenStore interface {
	Load() (Token, error)
	Save(Token) error
	Delete() error
}

type DeviceConfig struct {
	Host         string
	Tenant       string
	ClientID     string
	ClientSecret string
	Store        TokenStore
	Open         func(url string) error
}

type Device struct {
	app   *oauthApp
	store TokenStore
	open  func(string) error
	log   *slog.Logger
	now   func() time.Time

	mu      sync.Mutex
	token   Token
	pending *pendingDevice
}

type pendingDevice struct {
	code     string
	interval time.Duration
	expires  time.Time
}

func NewDevice(log *slog.Logger, cfg DeviceConfig) *Device {
	app := newGitHubApp(cfg.ClientID, cfg.ClientSecret, "")
	if cfg.Host == "azure" {
		app = newAzureApp(cfg.Tenant, cfg.ClientID, cfg.ClientSecret, "")
	}
	return &Device{
		app:   app,
		store: cfg.Store,
		open:  cfg.Open,
		log:   log,
		now:   time.Now,
	}
}

func (d *Device) Routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/"+d.app.name+"/login", d.login)
	mux.HandleFunc("POST /auth/"+d.app.name+"/device/poll", d.poll)
	mux.HandleFunc("POST /auth/"+d.app.name+"/logout", d.logout)
}

func (d *Device) Token(_ http.ResponseWriter, r *http.Request) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.token.Access == "" {
		stored, err := d.store.Load()
		if err != nil {
			if !errors.Is(err, ErrNoStoredToken) {
				d.log.Warn("read stored "+d.app.name+" token", "error", err)
			}
			return "", ErrNotLoggedIn
		}
		d.token = stored
	}
	if d.token.Expires.IsZero() || d.token.Expires.After(d.now().Add(refreshMargin)) {
		return d.token.Access, nil
	}
	next, err := d.refresh(r.Context())
	if err != nil {
		d.log.Info(d.app.name+" login expired and could not be refreshed", "error", err)
		d.forget()
		return "", ErrNotLoggedIn
	}
	d.save(next)
	return next.Access, nil
}

func (d *Device) refresh(ctx context.Context) (Token, error) {
	if d.token.Refresh == "" {
		return Token{}, errors.New("no refresh token")
	}
	return d.app.refresh(ctx, d.token.Refresh)
}

func (d *Device) save(token Token) {
	d.token = token
	if err := d.store.Save(token); err != nil {
		d.log.Error("store "+d.app.name+" token", "error", err)
	}
}

func (d *Device) forget() {
	d.token = Token{}
	if err := d.store.Delete(); err != nil {
		d.log.Error("delete stored "+d.app.name+" token", "error", err)
	}
}

type loginPage struct {
	Label    string
	Host     string
	UserCode string
	URL      string
	Interval int
	Return   string
}

var loginTemplate = template.Must(template.New("login").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><title>Log in to fuda</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>body{font-family:system-ui,sans-serif;max-width:28rem;margin:4rem auto;padding:0 1rem;text-align:center}
code{display:block;font-size:2.25rem;letter-spacing:.2em;margin:1.5rem 0}</style></head>
<body><h1>Log in to fuda</h1>
<p>{{.Label}} opened in your browser. Type this code there:</p>
<code>{{.UserCode}}</code>
<p>Browser did not open? Go to <a href="{{.URL}}">{{.URL}}</a>.</p>
<p id="status">Waiting for you…</p>
<script>
const back = {{.Return}};
let wait = {{.Interval}} * 1000;
async function poll() {
  const res = await fetch('/auth/{{.Host}}/device/poll', { method: 'POST' });
  const body = await res.json();
  if (body.status === 'done') { location.replace(back); return; }
  if (body.status === 'failed') { document.getElementById('status').textContent = body.error; return; }
  if (body.interval) wait = body.interval * 1000;
  setTimeout(poll, wait);
}
setTimeout(poll, wait);
</script></body></html>
`))

func (d *Device) login(w http.ResponseWriter, r *http.Request) {
	code, err := d.app.startDevice(r.Context())
	if err != nil {
		d.log.Warn(d.app.name+" device login failed to start", "error", err)
		http.Error(w, d.app.label+" login failed. Try again.", http.StatusBadGateway)
		return
	}
	interval := max(code.Interval, 5)
	d.mu.Lock()
	d.pending = &pendingDevice{code: code.DeviceCode, interval: time.Duration(interval) * time.Second, expires: d.now().Add(15 * time.Minute)}
	d.mu.Unlock()

	if d.open != nil {
		if err := d.open(code.VerificationURI); err != nil {
			d.log.Warn("open browser for "+d.app.name+" login", "error", err)
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = loginTemplate.Execute(w, loginPage{
		Label:    d.app.label,
		Host:     d.app.name,
		UserCode: code.UserCode,
		URL:      code.VerificationURI,
		Interval: interval,
		Return:   safeReturn(r.URL.Query().Get("return")),
	})
}

func (d *Device) poll(w http.ResponseWriter, r *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.pending == nil || d.now().After(d.pending.expires) {
		d.pending = nil
		writeStatus(w, map[string]any{"status": "failed", "error": "The login ran out of time. Go back and try again."})
		return
	}
	token, err := d.app.pollDevice(r.Context(), d.pending.code)
	switch {
	case errors.Is(err, errAuthorizationPending):
		writeStatus(w, map[string]any{"status": "pending", "interval": d.pending.interval.Seconds()})
	case errors.Is(err, errSlowDown):
		d.pending.interval += 5 * time.Second
		writeStatus(w, map[string]any{"status": "pending", "interval": d.pending.interval.Seconds()})
	case err != nil:
		d.log.Warn(d.app.name+" device login failed", "error", err)
		d.pending = nil
		writeStatus(w, map[string]any{"status": "failed", "error": d.app.label + " did not let you in. Go back and try again."})
	default:
		d.pending = nil
		d.save(token)
		writeStatus(w, map[string]any{"status": "done"})
	}
}

func (d *Device) logout(w http.ResponseWriter, _ *http.Request) {
	d.Logout(w)
	w.WriteHeader(http.StatusNoContent)
}

func (d *Device) Logout(http.ResponseWriter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.forget()
}

func writeStatus(w http.ResponseWriter, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func (d *Device) Host() string { return d.app.name }
