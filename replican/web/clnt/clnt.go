
package clnt

import (
	"bytes"
	"fmt"
	"gob"
	"http"
	"io"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-web/replican/web"
	
)

type RemoteStore struct {
	http.Client
	url string
	root *fs.Dir
	index *fs.BlockIndex
}

func Connect(url string) (*RemoteStore, os.Error) {
	store := &RemoteStore{ url: url }
	
	err := store.Pull()
	return store, err
}

func (store *RemoteStore) Pull() os.Error {
	req, err := http.NewRequest("GET", 
		fmt.Sprintf("%s/%s", store.url, web.NODEROOT), nil)
	if err != nil { return err }
	
	resp, err := store.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	
	decoder := gob.NewDecoder(resp.Body)
	
	store.root = &fs.Dir{}
	err = decoder.Decode(store.root)
	if err != nil {
		return err
	}
	
	store.index = fs.IndexBlocks(store.root)
	return nil
}

func (store *RemoteStore) Root() *fs.Dir {
	return store.root
}

func (store *RemoteStore) Index() *fs.BlockIndex {
	return store.index
}

func (store *RemoteStore) ReadBlock(strong string) ([]byte, os.Error) {
	req, err := http.NewRequest("GET", 
		fmt.Sprintf("%s/%s/%s", store.url, web.BLOCKS, strong), nil)
	if err != nil { return nil, err }
	
	resp, err := store.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	
	buffer := bytes.NewBuffer([]byte{})
	n, err := io.Copy(buffer, resp.Body)
	if err != nil { return nil, err }
	if n < int64(fs.BLOCKSIZE) {
		return nil, io.ErrShortWrite
	}
	
	return buffer.Bytes(), nil
}
	
func (store *RemoteStore) ReadInto(
		strong string, from int64, length int64, writer io.Writer) (int64, os.Error) {
	req, err := http.NewRequest("GET", 
		fmt.Sprintf("%s/%s/%s", store.url, web.FILES, strong), nil)
	if err != nil { return -1, err }
	
	resp, err := store.Do(req)
	if err != nil { return -1, err }
	defer resp.Body.Close()
	
	n, err := io.Copy(writer, resp.Body)
	if err != nil { return n, err }
	if n < length {
		return n, io.ErrUnexpectedEOF
	}
	
	return n, nil
}



