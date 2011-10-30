
package main

import (
	"fmt"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-web/replican/web/clnt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}
	
	store, err := clnt.Connect(os.Args[1])
	if err != nil { die(err) }
	
	fs.Walk(store.Root(), func(node fs.Node) bool {
		switch node.(type) {
		case *fs.Dir:
			d := node.(*fs.Dir)
			fmt.Printf("d\t%s\t%s\n", d.Strong(), d.Name())
		case *fs.File:
			f := node.(*fs.File)
			fmt.Printf("f\t%s\t%s\n", f.Strong(), f.Name())
		case *fs.Block:
			b := node.(*fs.Block)
			f := b.Parent().(*fs.File)
			fmt.Printf("b\t%s\t%s\t%d\n", node.Strong(), f.Name(), b.Position())
		}
		return true
	})
}

func die(err os.Error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.String())
	os.Exit(1)
}

