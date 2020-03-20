package general

import (
	"bytes"
	"github.com/goextension/log"
	"github.com/gorilla/rpc/v2/json2"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// CurrentDir ...
func CurrentDir() string {
	dir, err := os.Getwd()
	if err == nil {
		return dir
	}
	return "."
}

// SplitIP ...
func SplitIP(addr string) (ip string, port int) {
	if addr == "" {
		return
	}
	s := strings.Split(addr, ":")
	if len(s) < 2 {
		return
	}
	ip = s[0]
	port, _ = strconv.Atoi(s[1])
	return
}

// RPCPost ...
func RPCPost(url string, method string, input, output interface{}) error {
	log.Infow("rpc post", "url", url, "method", method, "input", input)
	message, err := json2.EncodeClientRequest(method, input)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	buf.Grow(bytes.MinRead)
	_, err = buf.ReadFrom(resp.Body)
	log.Infow("rpc result", "response", string(buf.Bytes()))
	err = json2.DecodeClientResponse(buf, output)
	if err != nil {
		return err
	}
	return nil
}
