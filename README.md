# gwolf

A SaaS example application

#### Prerequisite
- Docker (keycloak image)
- `go1.23.5`

#### Setup guide
Make sure you have the most recent go version installed. Then, run keycloak docker image locally. This step is needed for the application to run. You will also need to make sure you've imported the keycloak realm configurations. See `realm-export.json` file in the root directory of this repo.

Export the variables in the `.env.axample` file. 

To run the application, you either us the `go run main.go` or the `make` commands. Reference the makefile in the root dir.