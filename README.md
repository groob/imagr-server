A simple server for [Imagr](https://github.com/grahamgilbert/imagr)  

The server will dynamically generate imagr_config.plist from a list of workflows in /imagr_repo/workflows and then serve the plist over HTTP.  

# Configuration
imagr-server expects a password to be set using the `IMAGR_PASSWORD` environment variable.


# Usage
`$ IMAGR_PASSWORD="password" imagr-server -repo /path/to/imagr_repo`  
Use `-serve` to serve the repo over HTTP.  
`$ IMAGR_PASSWORD="password" imagr-server -repo /path/to/imagr_repo -serve`

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
I use `uuidgen` to name files based on UUID  
When the file is saved as a plist, the UUID becomes the filename.  
When returning JSON, the UUID will become an ID.  


Example API usage to add/remove workflows:

```
$ curl -H 'Content-Type: application/json' -X PUT -d @test.json "http://imagr/v1/workflows/$(uuidgen)"
E77EBEF4-D55B-4743-BEB8-B26E4C87E73F
$ curl -X DELETE http://imagr/v1/workflows/E77EBEF4-D55B-4743-BEB8-B26E4C87E73F
E77EBEF4-D55B-4743-BEB8-B26E4C87E73F
```

# Building

Install the [latest Go distribution](https://golang.org/dl).
Ensure you have a GOPATH defined. For example:

`mkdir ~/go; export GOPATH="$HOME/go"`

Go get imagr-server and its dependencies:

`go get -u github.com/groob/imagr-server`

Build it:

`cd "$GOPATH/src/github.com/groob/imagr-server"; go build`

You should now find the output binary `imagr-server` in the current working directory.
