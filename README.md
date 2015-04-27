A simple server for [Imagr](https://github.com/grahamgilbert/imagr)  

The server will dynamically generate imagr_config.plist from a list of workflows in /imagr_repo/workflows and then serve the plist over HTTP.  

# Configuration
imagr-server expects a password to be set using the `IMAGR_PASSWORD` environment variable.


# Usage
`$ IMAGR_PASSWORD="password" imagr-server /path/to/imagr_repo`

# Docker usage
```
docker pull groob/imagr-server
docker run -it --rm \
  -p 80:3000 \
  -e IMAGR_PASSWORD="password" \
  --name imagr-server \
  -v /path/to/repo:/imagr_repo \
  groob/imagr-server
```

# API usage
imagr-server now supports GET/PUT/DELETE operations on workflows.
Example:

```
~ ❯❯❯ curl -H 'Content-Type: application/json' localhost:3000/v1/workflows/three
Workflow does not exist.%                                                                                                                                                                      ~ ❯❯❯ curl -H 'Content-Type: application/json' -X PUT -d @test.json localhost:3000/v1/workflows/three
~ ❯❯❯ curl -H 'Content-Type: application/json' localhost:3000/v1/workflows/three
{
        "name": "Yosemite - MunkiTools",
        "description": "Deploys the latest 10.10.x image. Munki tools and local admin account included.",
        "components": [
                {
                        "type": "image",
                        "url": "http://imagr/images/BaseImage-10.10.3-14D131.hfs.dmg"
                }
        ],
        "restart_action": "restart"
}%                                                                                                                                                                                             ~ ❯❯❯ curl -H 'Content-Type: application/json' -X DELETE localhost:3000/v1/workflows/three
~ ❯❯❯ curl -H 'Content-Type: application/json' localhost:3000/v1/workflows/three
Workflow does not exist.%
```
