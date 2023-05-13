# go-proxy-yourself

go-proxy-yourself is an authentication proxy for protecting web applications.

This project is being done for fun and likely has plenty of bugs.

## Setup

### Generate Certs

You will need to have certs for the server to run. Generate the certs using
`scripts/generate_certs`.

### Prepare Hosts

Add a section to your `etc/hosts` file similar to the following.

```
127.0.0.1       simple.app.local
```

The name used here will be the host used to determine the SAML settings to use
from the config. i.e. route key

### Developing

Start the development server in a docker container with hot reload by running

```bash
make dev
```

This will make the main `go-proxy-yourself` server available at `localhost:8080`
and `localhost:8443`. Requests to http will be automatically routed to https.

Along with the main `go-proxy-yourself` container, a simple `httpd` server will
be launched as an example app to protect. This server will be available at
`localhost:8081` directly.

Either a SAML or OAuth provider will need to be configured to actually protect
the applications. Their are multiple providers available.

When the dev server is starting, will copy
`configs/default.yml` to `configs/dev.yml` if it doesn't exist.

### Authentication Providers

#### SAML (not fully working)

Create your application in okta.com. Add your username to the allowed users and go to the setup page and copy the idpSsoUrl and idpIssuer into the `config/default.yaml` file. Download the okta cert into the config folder of the project.

#### OAuth

Create a google development project and copy the googleClientId and
googleClientSecret into the `config/default.yaml` file.
