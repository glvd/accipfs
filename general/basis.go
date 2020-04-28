package general

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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
	logD("rpcpost", "url", url, "method", method, "input", input)
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
	var size int64
	size, err = buf.ReadFrom(resp.Body)
	if size == 0 {
		return errors.New("no data response from remote node")
	}
	err = json2.DecodeClientResponse(buf, output)
	if err != nil {
		logE("decode failed", "data", buf.String())
		return err
	}
	return nil
}

// RPCAddress ...
func RPCAddress(addr core.NodeAddress) string {
	if addr.Prefix == "" {
		addr.Prefix = "rpc"
	}
	return fmt.Sprintf("http://%s:%d/%s", addr.Address, addr.Port, addr.Prefix)
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

// PipeDummy ...
func PipeDummy(ctx context.Context, module string, reader io.Reader) (e error) {
	line := make([]byte, 2048)
END:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, e = reader.Read(line)
			if e != nil || io.EOF == e {
				break END
			}
		}
	}

	return nil
}
