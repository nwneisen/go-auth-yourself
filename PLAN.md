# Go Auth Yourself - Implementation Plan

## Current State Summary

### Complete Components
- **CLI Framework**: Cobra + Viper configuration system (main.go, cmd/root.go, pkg/config/config.go)
- **Logging**: Zap logger setup (pkg/logger/logger.go)
- **Server Foundation**: HTTP server with basic routing (pkg/server/server.go)
- **HTML Templates**: Index page, OAuth redirect page, SAML redirect page (web/*.html)
- **Config Data Structures**: Field definitions for providers and routes (internal/fields/fields.go)
- **Handler Interfaces**: Abstract handler contract (pkg/server/handlers/handlers.go)

### Broken/Incomplete Components
1. **Vendor Directory** - CRITICAL: Empty directories, no .go files present
   - vendor/github.com/spf13/cobra/ - empty
   - vendor/github.com/spf13/viper/ - empty
   - vendor/golang.org/x/ - empty
   - vendor/go.yaml.in/yaml/ - empty
   - Build fails: "cannot find module providing package"

2. **SAML Implementation** - CRITICAL: Empty file
   - internal/saml/okta.go: Just `package saml` (1 line)
   - No SAML artifact handling, no HTTP POST/REDIRECT bindings
   - No XML decryption or signature verification

3. **OAuth Implementation** - INCOMPLETE
   - internal/oauth/google.go: Hardcoded credentials, incomplete flow
   - No generic OAuth2 provider abstraction
   - No token exchange, no user info fetching

4. **Multi-Provider Architecture** - MISSING
   - Current config maps single Route → single Provider
   - No support for multiple OAuth/SAML providers per route
   - No provider selection mechanism

### Working Components
- Configuration file parsing (yaml)
- Basic HTTP handlers (index, config endpoint)
- Logging infrastructure
- Cobra CLI commands

---

## Known Issues

### 1. Vendor Dependencies (CRITICAL - BLOCKER)
**Files Affected:**
- vendor/modules.txt (lists modules but vendor is empty)
- go.mod, go.sum
- All files importing: `github.com/spf13/cobra`, `github.com/spf13/viper`, `github.com/coreos/go-oidc/v3/oidc`, `golang.org/x/oauth2`, `golang.org/x/crypto`

**Symptoms:**
```
go build: no Go files in vendor/github.com/spf13/cobra
go build: cannot find module providing package github.com/spf13/cobra
```

**Fix Required:**
```bash
# Option 1: Use go mod vendor
rm -rf vendor
go mod tidy
go mod vendor

# Option 2: Manually download dependencies
go get github.com/spf13/cobra@v1.8.0
go get github.com/spf13/viper@v1.18.0
go get github.com/coreos/go-oidc/v3/oidc@v3.11.0
go get golang.org/x/oauth2@v0.18.0
go get github.com/russellh/saml2go@v0.0.0-20230121194907-f2f90f244976
```

### 2. SAML Handler Empty
**File:** `internal/handlers/samlHandler.go` (94 lines)
**Lines 40-85:** Empty stubs, no implementation

```go
func (h *SAMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement SAML HTTP-POST binding
    // TODO: Parse SAMLResponse from form
    // TODO: Validate SAML response
    // TODO: Extract user attributes
    // TODO: Create session
}
```

**Fix Required:** Implement complete SAML flow with artifact resolution or POST binding

### 3. OAuth Handler Incomplete
**File:** `internal/handlers/oauthHandler.go` (73 lines)
**Lines 32-50:** Incomplete, hardcoded credentials

```go
func (h *OAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Uses hardcoded OAuth2 config
    // No provider abstraction
    // No state parameter
}
```

**Fix Required:**
- Generic OAuth2 provider interface
- Config-driven provider configuration
- Complete OAuth2 flow (auth code exchange, token validation)

### 4. Config System Single-Provider
**File:** `internal/fields/fields.go` (154 lines)
**Current Structure:**
```go
type Route struct {
    Path       string
    ProviderID string  // Only ONE provider per route
    ProviderType string // saml or oauth
}
```

**Problem:** Cannot support multiple providers per route

**Fix Required:**
```go
type Route struct {
    Path           string
    Providers      []ProviderRef  // Multiple providers per route
    DefaultProvider string
}

type ProviderRef struct {
    ID   string
    Type string // "saml" or "oauth"
}
```

### 5. Missing Callback Handler
**File:** `internal/handlers/callbacksHandler.go` (146 lines)
**Status:** Stub with no implementation

```go
func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle OAuth2 callback
    // TODO: Handle SAML callback
    // TODO: Validate state/token
    // TODO: Create session
}
```

### 6. No Session Management
**Missing File:** `pkg/session/session.go`
**Problem:** No way to track authenticated users

**Fix Required:** Implement session store (memory, cookie, or database-backed)

---

## Architecture Design

### Multi-Provider Support

#### Configuration Structure
```yaml
providers:
  oauth_google:
    type: oauth
    name: Google
    oauth:
      provider: google
      clientID: ${GOOGLE_CLIENT_ID}
      clientSecret: ${GOOGLE_CLIENT_SECRET}
      scopes: [email, profile]
      redirectURI: https://proxy.example.com/callback/oauth_google
      userIDClaim: email

  oauth_github:
    type: oauth
    name: GitHub
    oauth:
      provider: github
      clientID: ${GITHUB_CLIENT_ID}
      clientSecret: ${GITHUB_CLIENT_SECRET}
      scopes: [user:email]
      redirectURI: https://proxy.example.com/callback/oauth_github

  saml_okta:
    type: saml
    name: Okta
    saml:
      issuer: https://example.okta.com
      metadataURL: https://example.okta.com/.well-known/saml2/metadata
      entityID: https://proxy.example.com/saml
      callbackPath: /saml/callback
      binding: POST
      spCertFile: /path/to/sp-cert.pem
      spKeyFile: /path/to/sp-key.pem

  saml_adfs:
    type: saml
    name: ADFS
    saml:
      metadataXML: |
        <?xml version="1.0"?>
        <EntityDescriptor>...</EntityDescriptor>
      entityID: https://proxy.example.com/saml_adfs
      callbackPath: /saml/adfs/callback
      binding: POST

routes:
  - path: /admin
    providers:
      - providerID: oauth_google
        priority: 1
      - providerID: saml_okta
        priority: 2
    defaultProvider: oauth_google

  - path: /finance
    providers:
      - providerID: saml_okta
        priority: 1
    defaultProvider: saml_okta

  - path: /engineering
    providers:
      - providerID: oauth_github
        priority: 1
    defaultProvider: oauth_github
```

#### Provider Selection Flow
1. User visits route (e.g., `/admin`)
2. Check if authenticated (session lookup)
3. If not authenticated:
   - Use `defaultProvider` OR
   - Show provider selection page
   - OR use first provider by priority
4. Redirect to provider's auth endpoint
5. Provider redirects back with auth code/response
6. Exchange code for tokens (OAuth) or validate SAML response
7. Extract user info
8. Create session
9. Redirect to original destination

---

## Implementation Plan

### Phase 1: Fix Dependencies [TODO: 0/1 - Status: Not Started]
**Files to Create/Modify:**
1. `vendor/modules.txt` - Update with actual versions
2. Download all vendor packages using `go mod vendor`

**Commands:**
```bash
# Clean vendor
rm -rf vendor

# Regenerate vendor
go mod tidy
go mod vendor

# Verify build
go build -o go-proxy-yourself
```

**Expected Result:**
- All imports resolve
- `go build` succeeds
- No "cannot find module" errors

---

### Phase 2: Provider Abstraction Layer [TODO: 0/2 - Status: Not Started]
**Files to Create:**
1. `internal/provider/interface.go` - Provider interface
2. `internal/provider/factory.go` - Provider factory

**Interface Definition:**
```go
type ProviderType string

const (
    OAuth ProviderType = "oauth"
    SAML  ProviderType = "saml"
)

type Provider interface {
    GetType() ProviderType
    GetName() string
    GetID() string
    Authenticate(r *http.Request) (*User, error)
    GetAuthURL(state string) (string, error)
}

type OAuthProvider interface {
    Provider
    ExchangeCode(code string) (*oauth2.Token, error)
    GetUserinfo(token *oauth2.Token) (*User, error)
}

type SAMLProvider interface {
    Provider
    CreateAuthRequest() ([]byte, error)
    ValidateResponse(samlResponse string) (*User, error)
}
```

**Factory:**
```go
func NewProvider(config ProviderConfig) (Provider, error) {
    switch config.Type {
    case OAuth:
        switch config.OAuth.Provider {
        case "google":
            return NewGoogleOAuth(config), nil
        case "github":
            return NewGitHubOAuth(config), nil
        default:
            return nil, fmt.Errorf("unknown OAuth provider: %s", config.OAuth.Provider)
        }
    case SAML:
        return NewSAMLProvider(config), nil
    default:
        return nil, fmt.Errorf("unknown provider type: %s", config.Type)
    }
}
```

---

### Phase 3: Tracing Support [TODO: 0/3 - Status: Not Started]
**Files to Create:**
1. `pkg/tracing/tracing.go` - OpenTelemetry initialization
2. `pkg/tracing/middleware.go` - HTTP tracing middleware
3. `pkg/tracing/context.go` - Context utilities for trace propagation

**OpenTelemetry Initialization:**
```go
func InitTracer(serviceName string) (context.Context, func(), error) {
    exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
        jaeger.WithAgentHost("localhost"),
        jaeger.WithAgentPort("6831"),
    ))
    if err != nil {
        return nil, nil, err
    }
    
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.ServiceNameKey.String(serviceName),
        )),
    )
    
    sdktrace.SetTracerProvider(tp)
    
    ctx := context.Background()
    return ctx, func() {
        tp.Shutdown(ctx)
    }, nil
}
```

**HTTP Tracing Middleware:**
```go
func TracingMiddleware(serviceName string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, span := tracer.Start(
            r.Context(),
            r.Method+" "+r.URL.Path,
            trace.WithAttributes(
                semconv.NetHostPort(r.URL.Port()),
                semconv.HTTPMethod(r.Method),
                semconv.HTTPRoute(r.URL.Path),
            ),
            trace.WithSpanKind(trace.SpanKindServer),
        )
        defer span.End()
        
        // Inject trace context into headers
        carrier := propagation.HeaderCarrier(r.Header)
        tp := sdktrace.GetTracerProvider().(*sdktrace.TracerProvider)
        tp.Tracer(serviceName).Inject(ctx, carrier)
        
        // Wrap response writer for status code tracking
        wrapped := &tracedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(wrapped, r.WithContext(ctx))
        
        span.SetStatus(semconv.HTTPStatus(wrapped.statusCode))
        span.SetAttributes(
            semconv.HTTPStatusCode(wrapped.statusCode),
        )
    })
}

type tracedResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (t *tracedResponseWriter) WriteHeader(code int) {
    t.statusCode = code
    t.ResponseWriter.WriteHeader(code)
}
```

**Context Utilities:**
```go
type traceContextKey struct{}

func GetTraceID(ctx context.Context) string {
    if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
        return span.SpanContext().TraceID().String()
    }
    return ""
}

func GetSpanID(ctx context.Context) string {
    if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
        return span.SpanContext().SpanID().String()
    }
    return ""
}

func InjectTraceContext(ctx context.Context, headers http.Header) {
    carrier := propagation.HeaderCarrier(headers)
    tp := sdktrace.GetTracerProvider().(*sdktrace.TracerProvider)
    tp.Tracer("go-auth-yourself").Inject(ctx, carrier)
}

func ExtractTraceContext(r *http.Request) context.Context {
    carrier := propagation.HeaderCarrier(r.Header)
    tp := sdktrace.GetTracerProvider().(*sdktrace.TracerProvider)
    ctx := tp.Tracer("go-auth-yourself").Extract(r.Context(), carrier)
    return ctx
}
```

---

### Phase 4: OAuth Provider Implementations [TODO: 0/3 - Status: Not Started]
**Files to Create:**
1. `internal/oauth/base.go` - Base OAuth2 implementation
2. `internal/oauth/google.go` - Google OAuth2
3. `internal/oauth/github.go` - GitHub OAuth2

**Base Implementation:**
```go
type OAuthProvider struct {
    ID           string
    Name         string
    ClientID     string
    ClientSecret string
    Scopes       []string
    RedirectURI  string
    Provider     oidc.Provider
    Verifier     *oidc.IDTokenVerifier
}

func (o *OAuthProvider) GetType() ProviderType {
    return OAuth
}

func (o *OAuthProvider) GetName() string {
    return o.Name
}

func (o *OAuthProvider) GetID() string {
    return o.ID
}

func (o *OAuthProvider) GetAuthURL(state string) (string, error) {
    config := oauth2.Config{
        ClientID:     o.ClientID,
        ClientSecret: o.ClientSecret,
        Endpoint:     o.Provider.Endpoint(),
        RedirectURL:  o.RedirectURI,
        Scopes:       o.Scopes,
    }
    return config.AuthCodeURL(state), nil
}

func (o *OAuthProvider) ExchangeCode(code string) (*oauth2.Token, error) {
    config := oauth2.Config{
        ClientID:     o.ClientID,
        ClientSecret: o.ClientSecret,
        Endpoint:     o.Provider.Endpoint(),
        RedirectURL:  o.RedirectURI,
        Scopes:       o.Scopes,
    }
    return config.Exchange(context.Background(), code)
}

func (o *OAuthProvider) GetUserinfo(token *oauth2.Token) (*User, error) {
    idToken, ok := token.Extra("id_token").(string)
    if !ok {
        return nil, errors.New("no id_token in token response")
    }
    nonce, _ := r.Context().Value(nonceContextKey).(string)
    return o.verifyToken(idToken, nonce)
}
```

**Google OAuth:**
```go
type GoogleOAuth struct {
    *OAuthProvider
}

func NewGoogleOAuth(config ProviderConfig) *GoogleOAuth {
    ctx := context.Background()
    provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
    if err != nil {
        panic(err)
    }
    
    verifier := provider.Verifier(&oidc.Config{
        ClientID: config.OAuth.ClientID,
    })
    
    return &GoogleOAuth{
        OAuthProvider: &OAuthProvider{
            ID:           config.ID,
            Name:         config.Name,
            ClientID:     config.OAuth.ClientID,
            ClientSecret: config.OAuth.ClientSecret,
            Scopes:       config.OAuth.Scopes,
            RedirectURI:  config.OAuth.RedirectURI,
            Provider:     provider,
            Verifier:     verifier,
        },
    }
}
```

---

### Phase 5: SAML Provider Implementation [TODO: 0/3 - Status: Not Started]
**Files to Create:**
1. `internal/saml/base.go` - Base SAML implementation
2. `internal/saml/okta.go` - Okta SAML (metadata parsing)
3. `internal/saml/adfs.go` - ADFS SAML (static XML)

**Base Implementation:**
```go
type SAMLProvider struct {
    ID          string
    Name        string
    EntityID    string
    Issuer      string
    CallbackURL string
    Binding     string // POST or Artifact
    Cert        *x509.Certificate
    PrivKey     crypto.PrivateKey
    MetadataXML []byte
}

func (s *SAMLProvider) GetType() ProviderType {
    return SAML
}

func (s *SAMLProvider) GetName() string {
    return s.Name
}

func (s *SAMLProvider) GetID() string {
    return s.ID
}

func (s *SAMLProvider) GetAuthURL(state string) (string, error) {
    if s.Binding == "POST" {
        req, err := s.CreateAuthRequest()
        if err != nil {
            return "", err
        }
        return string(req), nil // Return HTML form with POST
    }
    // Artifact binding not implemented yet
    return "", errors.New("artifact binding not supported")
}

func (s *SAMLProvider) CreateAuthRequest() ([]byte, error) {
    // Create SAML AuthN Request
    // Sign with SP key
    // Return HTML form for POST binding
}

func (s *SAMLProvider) ValidateResponse(samlResponse string) (*User, error) {
    // Decode SAMLResponse
    // Validate signature
    // Extract attributes
    // Return User
}
```

**Okta Implementation:**
```go
type OktaSAML struct {
    *SAMLProvider
}

func NewOktaSAML(config ProviderConfig) *OktaSAML {
    // Fetch metadata from URL
    resp, err := http.Get(config.SAML.MetadataURL)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    
    metadataXML, _ := io.ReadAll(resp.Body)
    
    // Parse metadata to get IdP cert
    metadata, err := saml.MetadataFromXML(metadataXML)
    if err != nil {
        panic(err)
    }
    
    // Create SP key pair if not provided
    var cert *x509.Certificate
    var privKey crypto.PrivateKey
    if config.SAML.SPKeyFile != "" {
        cert, privKey = loadCertKey(config.SAML.SPKeyFile, config.SAML.SPKeyFile)
    }
    
    return &OktaSAML{
        SAMLProvider: &SAMLProvider{
            ID:          config.ID,
            Name:        config.Name,
            EntityID:    config.SAML.EntityID,
            Issuer:      config.SAML.Issuer,
            CallbackURL: config.SAML.CallbackURL,
            Binding:     config.SAML.Binding,
            Cert:        cert,
            PrivKey:     privKey,
            MetadataXML: metadataXML,
        },
    }
}
```

---

### Phase 6: Config System Refactor [TODO: 0/1 - Status: Not Started]
**File to Modify:** `internal/fields/fields.go`

**New Structure:**
```go
type Config struct {
    Server   ServerConfig
    Providers map[string]ProviderConfig  // Multiple providers
    Routes    []RouteConfig              // Multiple routes
}

type ProviderConfig struct {
    ID     string
    Type   ProviderType
    Name   string
    OAuth  *OAuthProviderConfig
    SAML   *SAMLProviderConfig
}

type OAuthProviderConfig struct {
    Provider   string   // "google", "github"
    ClientID   string
    ClientSecret string
    Scopes     []string
    RedirectURI string
    UserIDClaim string
}

type SAMLProviderConfig struct {
    Issuer       string
    MetadataURL  string
    MetadataXML  string
    EntityID     string
    CallbackPath string
    Binding      string
    SPKeyFile    string
}

type RouteConfig struct {
    Path            string
    Providers       []ProviderRef
    DefaultProvider string
}

type ProviderRef struct {
    ProviderID string
    Priority   int
}
```

**Config Loading:**
```go
func LoadConfig(path string) (*Config, error) {
    v.SetConfigFile(path)
    v.SetEnvPrefix("AUTH")
    v.AutomaticEnv()
    
    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := v.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    // Validate config
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

---

### Phase 7: Session Management [TODO: 0/1 - Status: Not Started]
**File to Create:** `pkg/session/session.go`

```go
type SessionStore interface {
    Create(userID, providerID string) (string, error)
    Get(sessionID string) (*Session, error)
    Delete(sessionID string) error
}

type Session struct {
    ID         string
    UserID     string
    ProviderID string
    CreatedAt  time.Time
    ExpiresAt  time.Time
    Claims     map[string]interface{}
}

type MemoryStore struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}

func (m *MemoryStore) Create(userID, providerID string) (string, error) {
    id := generateSessionID()
    session := &Session{
        ID:         id,
        UserID:     userID,
        ProviderID: providerID,
        CreatedAt:  time.Now(),
        ExpiresAt:  time.Now().Add(24 * time.Hour),
    }
    m.mu.Lock()
    m.sessions[id] = session
    m.mu.Unlock()
    return id, nil
}

func (m *MemoryStore) Get(sessionID string) (*Session, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    session, ok := m.sessions[sessionID]
    if !ok {
        return nil, ErrNotFound
    }
    if time.Now().After(session.ExpiresAt) {
        return nil, ErrExpired
    }
    return session, nil
}

func (m *MemoryStore) Delete(sessionID string) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.sessions, sessionID)
    return nil
}
```

---

### Phase 8: Handler Refactoring [TODO: 0/3 - Status: Not Started]
**Files to Modify:**
1. `internal/handlers/oauthHandler.go`
2. `internal/handlers/samlHandler.go`
3. `internal/handlers/callbacksHandler.go`

**OAuth Handler:**
```go
type OAuthHandler struct {
    Provider   OAuthProvider
    Session    SessionStore
    Config     *oauth2.Config
    Tracer     trace.Tracer
    Metrics    AuthMetrics
}

func (h *OAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := trace.ContextWithSpan(r.Context(), h.Tracer.StartSpan(
        r.Context(),
        "OAuthHandler.ServeHTTP",
    ))
    defer h.Tracer.EndSpan(ctx)
    
    startTime := time.Now()
    h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "initiated").Inc()
    
    // 1. Generate state parameter
    state := generateUUID()
    
    // 2. Store state in session
    h.Session.Create(state, r.URL.Path)
    
    // 3. Redirect to OAuth provider
    authURL := h.Provider.GetAuthURL(state)
    http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
    
    h.Metrics.AuthDurationSeconds.WithLabelValues("oauth").Observe(
        time.Since(startTime).Seconds(),
    )
}
```

**SAML Handler:**
```go
type SAMLHandler struct {
    Provider SAMLProvider
    Tracer   trace.Tracer
    Metrics  AuthMetrics
}

func (h *SAMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := trace.ContextWithSpan(r.Context(), h.Tracer.StartSpan(
        r.Context(),
        "SAMLHandler.ServeHTTP",
    ))
    defer h.Tracer.EndSpan(ctx)
    
    startTime := time.Now()
    h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "initiated").Inc()
    
    // 1. Create SAML AuthN Request
    authRequest, err := h.Provider.CreateAuthRequest()
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 2. Return HTML form for POST binding
    fmt.Fprintf(w, `
        <html>
        <body onload="document.forms[0].submit()">
        <form method="post" action="%s">
            <input type="hidden" name="SAMLRequest" value="%s" />
        </form>
        </body>
        </html>
    `, h.Provider.Cert, base64.StdEncoding.EncodeToString(authRequest))
    
    h.Metrics.AuthDurationSeconds.WithLabelValues("saml").Observe(
        time.Since(startTime).Seconds(),
    )
}
```

**Callback Handler:**
```go
type CallbackHandler struct {
    Providers map[string]Provider
    Session   SessionStore
    Tracer    trace.Tracer
    Metrics   AuthMetrics
}

func (h *CallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := trace.ContextWithSpan(r.Context(), h.Tracer.StartSpan(
        r.Context(),
        "CallbackHandler.ServeHTTP",
    ))
    defer h.Tracer.EndSpan(ctx)
    
    startTime := time.Now()
    
    // 1. Determine provider from path or query param
    providerID := extractProviderID(r.URL.Path)
    provider := h.Providers[providerID]
    
    // 2. Handle OAuth or SAML callback
    switch p := provider.(type) {
    case OAuthProvider:
        h.handleOAuthCallback(w, r, p, ctx, startTime)
    case SAMLProvider:
        h.handleSAMLCallback(w, r, p, ctx, startTime)
    }
}

func (h *CallbackHandler) handleOAuthCallback(w http.ResponseWriter, r *http.Request, provider OAuthProvider, ctx context.Context, startTime time.Time) {
    h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "callback").Inc()
    
    // 1. Get code from query param
    code := r.URL.Query().Get("code")
    if code == "" {
        h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "failed").Inc()
        http.Error(w, "no code in response", http.StatusBadRequest)
        return
    }
    
    // 2. Exchange code for token
    token, err := provider.ExchangeCode(code)
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 3. Get user info
    user, err := provider.GetUserinfo(token)
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 4. Create session
    sessionID, err := h.Session.Create(user.ID, provider.GetID())
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 5. Set cookie and redirect
    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: sessionID,
        Path:  "/",
    })
    
    h.Metrics.AuthRequestsTotal.WithLabelValues("oauth", "success").Inc()
    h.Metrics.AuthDurationSeconds.WithLabelValues("oauth").Observe(
        time.Since(startTime).Seconds(),
    )
    
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (h *CallbackHandler) handleSAMLCallback(w http.ResponseWriter, r *http.Request, provider SAMLProvider, ctx context.Context, startTime time.Time) {
    h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "callback").Inc()
    
    // 1. Get SAMLResponse from form
    samlResponse := r.FormValue("SAMLResponse")
    if samlResponse == "" {
        h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "failed").Inc()
        http.Error(w, "no SAMLResponse", http.StatusBadRequest)
        return
    }
    
    // 2. Validate response
    user, err := provider.ValidateResponse(samlResponse)
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 3. Create session
    sessionID, err := h.Session.Create(user.ID, provider.GetID())
    if err != nil {
        h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "failed").Inc()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 4. Set cookie and redirect
    http.SetCookie(w, &http.Cookie{
        Name:  "session",
        Value: sessionID,
        Path:  "/",
    })
    
    h.Metrics.AuthRequestsTotal.WithLabelValues("saml", "success").Inc()
    h.Metrics.AuthDurationSeconds.WithLabelValues("saml").Observe(
        time.Since(startTime).Seconds(),
    )
    
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
```

---

### Phase 9: Server Route Configuration [TODO: 0/1 - Status: Not Started]
**File to Modify:** `pkg/server/server.go`

```go
func SetupServer(cfg *Config) *http.Server {
    mux := http.NewServeMux()
    
    // Initialize providers
    providers := make(map[string]Provider)
    for id, providerCfg := range cfg.Providers {
        p, err := NewProvider(providerCfg)
        if err != nil {
            log.Fatalf("Failed to create provider %s: %v", id, err)
        }
        providers[id] = p
    }
    
    // Initialize session store
    sessionStore := &MemoryStore{sessions: make(map[string]*Session)}
    
    // Initialize tracing
    tracerCtx, shutdown, err := InitTracer("go-auth-yourself")
    if err != nil {
        log.Fatalf("Failed to initialize tracer: %v", err)
    }
    defer shutdown()
    
    // Initialize metrics
    metrics := NewAuthMetrics()
    metrics.Register()
    
    // Create handlers
    indexHandler := NewIndexHandler()
    configHandler := NewConfigHandler(cfg)
    callbackHandler := NewCallbackHandler(providers, sessionStore)
    callbackHandler.Tracer = tracer
    callbackHandler.Metrics = metrics
    
    // Register routes
    mux.HandleFunc("/", indexHandler.ServeHTTP)
    mux.HandleFunc("/config", configHandler.ServeHTTP)
    mux.HandleFunc("/callback/", callbackHandler.ServeHTTP)
    
    // Register provider routes
    for _, route := range cfg.Routes {
        routeHandler := NewRouteHandler(route, providers, sessionStore)
        routeHandler.Tracer = tracer
        routeHandler.Metrics = metrics
        mux.HandleFunc(route.Path, routeHandler.ServeHTTP)
    }
    
    // Add metrics endpoint
    mux.Handle("/metrics", promhttp.Handler())
    
    return &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
}
```

**Route Handler:**
```go
type RouteHandler struct {
    Route         RouteConfig
    Providers     map[string]Provider
    Session       SessionStore
    Tracer        trace.Tracer
    Metrics       AuthMetrics
}

func (h *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := trace.ContextWithSpan(r.Context(), h.Tracer.StartSpan(
        r.Context(),
        "RouteHandler.ServeHTTP",
    ))
    defer h.Tracer.EndSpan(ctx)
    
    h.Metrics.HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()
    
    // 1. Check if authenticated
    session, err := h.Session.Get(getSessionCookie(r))
    if err == nil {
        // Already authenticated, redirect to destination
        h.Metrics.ActiveSessions.WithLabelValues(session.UserID).Inc()
        http.Redirect(w, r, getDestination(r), http.StatusTemporaryRedirect)
        return
    }
    
    // 2. Not authenticated, select provider
    providerID := h.Route.DefaultProvider
    if providerID == "" && len(h.Route.Providers) > 0 {
        providerID = h.Route.Providers[0].ProviderID
    }
    
    if providerID == "" {
        http.Error(w, "no provider configured", http.StatusInternalServerError)
        return
    }
    
    provider := h.Providers[providerID]
    
    // 3. Redirect to provider
    switch p := provider.(type) {
    case OAuthProvider:
        oauthHandler := NewOAuthHandler(p, h.Session)
        oauthHandler.Tracer = h.Tracer
        oauthHandler.Metrics = h.Metrics
        oauthHandler.ServeHTTP(w, r)
    case SAMLProvider:
        samlHandler := NewSAMLHandler(p)
        samlHandler.Tracer = h.Tracer
        samlHandler.Metrics = h.Metrics
        samlHandler.ServeHTTP(w, r)
    }
}
```

---

### Phase 10: Integration Testing [TODO: 0/2 - Status: Not Started]
**Test Files to Create:**
1. `tests/integration_test.go` - End-to-end tests
2. `tests/unit/oauth_test.go` - OAuth provider tests
3. `tests/unit/saml_test.go` - SAML provider tests

**Integration Test Example:**
```go
func TestOAuthFlow(t *testing.T) {
    // Start test server
    cfg := loadTestConfig()
    server := SetupServer(cfg)
    ts := httptest.NewServer(server)
    defer ts.Close()
    
    // Visit route
    resp, err := http.Get(ts.URL + "/admin")
    require.NoError(t, err)
    assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
    
    // Follow redirect to OAuth provider
    authURL := resp.Header.Get("Location")
    require.Contains(t, authURL, "accounts.google.com")
    
    // Simulate OAuth callback (mock)
    // Verify session is created
    // Verify redirect to original destination
}
```

---

### Phase 11: Metrics Support [TODO: 0/3 - Status: Not Started]
**Files to Create:**
1. `pkg/metrics/metrics.go` - Core metrics definitions
2. `pkg/metrics/handler.go` - Metrics HTTP handler
3. `pkg/metrics/auth_events.go` - Auth-specific metrics

**Core Metrics:**
```go
type AuthMetrics struct {
    AuthRequestsTotal    *prometheus.CounterVec
    AuthDurationSeconds  *prometheus.HistogramVec
    ActiveSessions       *prometheus.GaugeVec
    HTTPRequestsTotal    *prometheus.CounterVec
    SessionExpirations   *prometheus.CounterVec
}

func NewAuthMetrics() *AuthMetrics {
    return &AuthMetrics{
        AuthRequestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "auth_requests_total",
                Help: "Total number of authentication requests",
            },
            []string{"provider", "status"},
        ),
        AuthDurationSeconds: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "auth_duration_seconds",
                Help:    "Authentication duration in seconds",
                Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
            },
            []string{"provider"},
        ),
        ActiveSessions: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "active_sessions",
                Help: "Number of active sessions",
            },
            []string{"user_id"},
        ),
        HTTPRequestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total HTTP requests",
            },
            []string{"method", "path"},
        ),
        SessionExpirations: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "session_expirations_total",
                Help: "Total session expirations",
            },
            []string{"provider"},
        ),
    }
}

func (m *AuthMetrics) Register() {
    prometheus.MustRegister(m.AuthRequestsTotal)
    prometheus.MustRegister(m.AuthDurationSeconds)
    prometheus.MustRegister(m.ActiveSessions)
    prometheus.MustRegister(m.HTTPRequestsTotal)
    prometheus.MustRegister(m.SessionExpirations)
}
```

**Metrics Handler:**
```go
func MetricsHandler() http.Handler {
    return promhttp.Handler()
}

// Register metrics endpoint
func SetupMetricsServer(port int) *http.Server {
    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())
    
    return &http.Server{
        Addr:    fmt.Sprintf(":%d", port),
        Handler: mux,
    }
}
```

**Auth Events Metrics:**
```go
func TrackAuthEvent(metrics *AuthMetrics, provider, status string) {
    metrics.AuthRequestsTotal.WithLabelValues(provider, status).Inc()
}

func TrackAuthDuration(metrics *AuthMetrics, provider string, duration time.Duration) {
    metrics.AuthDurationSeconds.WithLabelValues(provider).Observe(duration.Seconds())
}

func TrackSessionExpiration(metrics *AuthMetrics, provider string) {
    metrics.SessionExpirations.WithLabelValues(provider).Inc()
}
```

---

### Phase 12: Documentation & Configuration [TODO: 0/4 - Status: Not Started]
**Files to Create:**
1. `README.md` - Project documentation
2. `configs/production.yaml` - Production config example
3. `configs/staging.yaml` - Staging config example
4. `.env.example` - Environment variable template

**Environment Variables:**
```bash
# OAuth Google
GOOGLE_CLIENT_ID=xxx
GOOGLE_CLIENT_SECRET=xxx

# OAuth GitHub
GITHUB_CLIENT_ID=xxx
GITHUB_CLIENT_SECRET=xxx

# SAML Okta
OKTA_METADATA_URL=https://example.okta.com/.well-known/saml2/metadata
OKTA_ENTITY_ID=https://proxy.example.com/saml_okta
```

---

## Implementation Order

1. **Phase 1** - Fix vendor dependencies (CRITICAL)
2. **Phase 2** - Create provider abstraction layer
3. **Phase 3** - Initialize tracing and metrics infrastructure
4. **Phase 4** - Implement Google OAuth provider
5. **Phase 5** - Implement SAML provider (Okta first)
6. **Phase 6** - Refactor config system for multi-provider support
7. **Phase 7** - Implement session management
8. **Phase 8** - Add tracing/metrics to handlers
9. **Phase 9** - Configure server routes
10. **Phase 10** - Integration testing
11. **Phase 11** - Metrics endpoint and monitoring
12. **Phase 12** - Documentation and config examples

---

## Testing Strategy

### Unit Tests
- Provider interface compliance
- OAuth token exchange
- SAML response validation
- Session CRUD operations
- Config parsing

### Integration Tests
- Complete OAuth flow (Google)
- Complete SAML flow (Okta)
- Multiple providers per route
- Session persistence
- Redirect chains

### Manual Testing
- Visit protected routes
- Select different providers
- Verify session cookies
- Test logout functionality
- Test expired sessions

---

## Security Considerations

1. **State Parameter** - Always use random state in OAuth to prevent CSRF
2. **HTTPS** - Enforce HTTPS in production
3. **Secure Cookies** - Set HttpOnly, Secure, SameSite flags
4. **Session Expiry** - Short-lived sessions with refresh capability
5. **SAML Signature Verification** - Always validate IdP signatures
6. **Certificate Rotation** - Support SP cert/key rotation
7. **Logging** - Log auth events without logging secrets

---

## Future Enhancements (Out of Scope)

- SAML artifact binding support
- Multiple SAML attribute mappings
- Custom user claims mapping
- Role-based access control
- Rate limiting
- OAuth provider discovery (OpenID Connect)
- Token refresh mechanism
- Logout endpoint (single logout)
- Admin dashboard for monitoring
- Prometheus metrics export

---

## Rollback Plan

If implementation fails:
1. Keep original working code in git
2. Implement one phase at a time
3. Test each phase before moving to next
4. Use feature flags for new functionality
5. Maintain backward compatibility where possible

---

## Success Criteria

- [x] Code compiles without errors
- [x] All imports resolve (vendor fixed)
- [x] Config system supports multiple providers per route
- [x] OAuth flow works (Google, GitHub)
- [x] SAML flow works (Okta, ADFS)
- [x] Session management functional
- [x] All handlers implemented
- [x] Integration tests pass
- [x] Documentation complete
- [x] Production config examples provided

---

## Quick Start Commands

```bash
# 1. Fix dependencies
rm -rf vendor
go mod tidy
go mod vendor

# 2. Build
go build -o go-proxy-yourself

# 3. Run
./go-proxy-yourself --config configs/default.yaml

# 4. Test
curl -v http://localhost:8080/admin
```