package general

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/gorilla/rpc/v2/json2"
	"io"
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
	err = json2.DecodeClientResponse(buf, output)
	if err != nil {
		return err
	}
	return nil
}

// RPCAddress ...
func RPCAddress(addr core.NodeAddress) string {
	return fmt.Sprintf("http://%s:%d", addr.Address, addr.Port)
}

// PipeScreen ...
func PipeScreen(ctx context.Context, module string, reader io.Reader) (e error) {
	r := bufio.NewReader(reader)
	var lines []byte
END:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			lines, _, e = r.ReadLine()
			if e != nil || io.EOF == e {
				break END
			}
			if strings.TrimSpace(string(lines)) != "" {
				fmt.Printf("[%s]:%+v\n", module, string(lines))
			}
		}
	}

	return nil
}
