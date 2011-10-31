
package main

import (
	"fmt"
	"http"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-web/replican/web"
	"github.com/cmars/replican-web/replican/web/srv"
	
	"gorilla.googlecode.com/hg/gorilla/mux"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path>\n", os.Args[0])
		os.Exit(1)
	}
	
	local, err := fs.NewLocalStore(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.String())
		os.Exit(1)
	}
	
	var router = new(mux.Router)
	router.HandleFunc(fmt.Sprintf("/%s", web.NODEROOT), srv.TreeHandler(local))
	router.HandleFunc(fmt.Sprintf("/%s/{strong}", web.BLOCKS), srv.BlockHandler(local))
	router.HandleFunc(fmt.Sprintf("/%s/{strong}/{offset}/{length}", web.FILES),
		srv.FileHandler(local))
	
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

