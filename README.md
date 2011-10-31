
# replican-web - HTTP filesystem synchronization support #

replican-web adds remote filesystem synchronization support to [replican-sync](https://github.com/cmars/replican-sync)
over HTTP connections.

This is done by exposing fs.BlockStore functionality as a simple web API. 
On the client side, we proxy the web API with a fs.BlockStore implementation.  

[Gorilla mux](http://gorilla-web.appspot.com/pkg/gorilla/mux/) is used for the web API. Pretty sweet.

### Implemented ###

* Toy example client and server

### Planned/In development ###

* Contribute protocol to a multi-protocol synchronization peer.
* Performance benchmarking

Server:

* Periodic background local store scanning, inotify
* Respond appropriate HTTP status when not ready
  * Startup
  * Multiple clients
  * Updating index

Client:

* Bi-directional sync (currently only pull supported)



