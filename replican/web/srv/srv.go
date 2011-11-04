package srv

import (
	"bytes"
	"fmt"
	"gob"
	"http"
	"io"
	"strconv"

	"github.com/cmars/replican-sync/replican/fs"

	"gorilla.googlecode.com/hg/gorilla/mux"
)

func writeResponseError(resp http.ResponseWriter, statusCode int, msg string) {
	resp.WriteHeader(statusCode)
	resp.Write([]byte(msg))
}

func setBinaryResp(resp http.ResponseWriter) {
	resp.Header().Add("Content-Type", "application/octet-stream")
	resp.Header().Add("Content-Transfer-Encoding", "binary")
}

func TreeHandler(store *fs.LocalStore) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		setBinaryResp(resp)

		// gob the root
		buffer := bytes.NewBuffer([]byte{})
		encoder := gob.NewEncoder(buffer)
		err := encoder.Encode(store.Root())
		if err != nil {
			writeResponseError(resp, http.StatusInternalServerError, err.String())
			return
		}

		rootGob := buffer.Bytes()

		resp.Write(rootGob)
	}
}

func BlockHandler(store *fs.LocalStore) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		setBinaryResp(resp)

		strong, hasVar := mux.Vars(req)["strong"]
		if !hasVar {
			writeResponseError(resp, http.StatusInternalServerError,
				"Missing parameter: strong")
			return
		}

		if !hasVar {
			resp.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := store.ReadBlock(strong)
		if err != nil {
			writeResponseError(resp, http.StatusInternalServerError, err.String())
			return
		}

		resp.Write(buf)
	}
}

func FileHandler(store *fs.LocalStore) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		setBinaryResp(resp)

		strong, hasStrong := mux.Vars(req)["strong"]
		if !hasStrong {
			writeResponseError(resp, http.StatusInternalServerError,
				"Missing parameter: strong")
			return
		}

		offsetStr, hasOffset := mux.Vars(req)["offset"]
		if !hasOffset {
			writeResponseError(resp, http.StatusInternalServerError,
				"Missing parameter: offset")
			return
		}
		offset, err := strconv.Atoi64(offsetStr)
		if err != nil {
			writeResponseError(resp, http.StatusInternalServerError,
				fmt.Sprintf("Invalid format for length: %s", offsetStr))
			return
		}

		lengthStr, hasLength := mux.Vars(req)["length"]
		if !hasLength {
			writeResponseError(resp, http.StatusInternalServerError,
				"Missing parameter: length")
			return
		}
		length, err := strconv.Atoi64(lengthStr)
		if err != nil {
			writeResponseError(resp, http.StatusInternalServerError,
				fmt.Sprintf("Invalid format for length: %s", lengthStr))
			return
		}

		buffer := bytes.NewBuffer([]byte{})
		n, err := store.ReadInto(strong, offset, length, buffer)
		if err != nil {
			writeResponseError(resp, http.StatusInternalServerError, err.String())
			return
		}
		if n < length {
			writeResponseError(resp, http.StatusInternalServerError, io.ErrShortWrite.String())
		}

		resp.Write(buffer.Bytes())
	}
}
