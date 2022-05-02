# go-proxy-yourself

go-proxy-yourself is an authentication proxy for protecting web applications.

This project is being done for fun and likely has plenty of bugs.

## Setup

The project will launch a simple web server to protect while developing. The following will run the projects code but either SAML or OAuth will need to be configured to actually protect the applications.

1. Add a section to your `etc/hosts` file similar to the following.
```
127.0.0.1       simple.app.local
```

2. Start the development server
``` bash
make dev
```

This will start the development containers. The `go-proxy-yourself` container does hot reloading and will rebuild the code when changes are detected to the projects files.

3. Setup one of the following authentication procedures

### SAML (not fully working)

Create your application in okta.com. Add your username to the allowed users and go to the setup page and copy the idpSsoUrl and idpIssuer into the `config/default.yaml` file. Download the okta cert into the config folder of the project.

### OAuth

Create a google development project and copy the googleClientId and googleClientSecret into the `config/default.yaml` file.
