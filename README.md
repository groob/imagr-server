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
