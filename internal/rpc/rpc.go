package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/readline.v1"
)

type Command struct {
	Addr        string
	Method      string
	Args        []string
	UserAgent   string
	HTTP        bool
	Interactive bool
	Input       io.Reader
	Output      io.Writer
	sock        net.Conn
	client      *rpc.Client
	id          int
}

func New() *Command {
	return new(Command)
}

func (c *Command) Run() error {
	if !c.HTTP {
		err := c.connect()
		if err != nil {
			return err
		}
	}

	if c.Method == "" {
		return c.interactive()
	}

	if args := c.Args; len(args) > 0 {
		req := request(args)
		return c.call(c.Method, req)
	}

	if c.Input == nil {
		return c.call(c.Method, nil)
	}

	dec := json.NewDecoder(c.Input)
	dec.UseNumber()

	for {
		var req interface{}
		err := dec.Decode(&req)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		err = c.call(c.Method, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) connect() (err error) {
	sock, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}

	c.sock = sock
	c.client = jsonrpc.NewClient(sock)
	return err
}

func (c *Command) interactive() error {
	rl, err := readline.New(c.Addr + "> ")
	if err != nil {
		return err
	}

	defer func() {
		cerr := rl.Close()
		if err == nil {
			err = cerr
		}
	}()

	for {
		l, err := rl.Readline()

		if err == readline.ErrInterrupt {
			break
		}

		if err != nil {
			return err
		}

		args, err := shellwords.Parse(l)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			continue
		}

		method := args[0]
		req := request(args[1:])
		err = c.call(method, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) call(method string, req interface{}) error {
	if c.HTTP {
		return c.post(method, req)
	}

	var reply *json.RawMessage
	err := c.client.Call(method, req, &reply)

	if err == rpc.ErrShutdown {
		err = c.connect()
		if err != nil {
			return fmt.Errorf("cannot re-connect %s", err)
		}

		err = c.client.Call(method, req, &reply)
	}

	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(reply, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Output, string(buf))
	return nil
}

func (c *Command) post(method string, req interface{}) error {
	r, err := c.request(method, req)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	if c.UserAgent != "" {
		r.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := http.DefaultClient.Do(r)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		out, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("received non json response\n%s", out)
	}

	res := struct {
		Result interface{} `json:"result"`
		Error  interface{} `json:"error"`
	}{}

	dec := json.NewDecoder(resp.Body)
	dec.UseNumber()
	err = dec.Decode(&res)
	if err != nil {
		return err
	}

	if res.Error != nil {
		fmt.Printf("%+v\n", res.Error)
		return nil
	}

	buf, err := json.MarshalIndent(res.Result, "", "  ")
	if err != nil {
		return err
	}

	fmt.Fprintln(c.Output, string(buf))
	return nil
}

func (c *Command) request(method string, req interface{}) (*http.Request, error) {
	c.id++

	buf, err := json.Marshal(map[string]interface{}{
		"id":     c.id,
		"method": method,
		"params": []interface{}{req},
	})
	if err != nil {
		return nil, err
	}

	return http.NewRequest("POST", c.Addr, bytes.NewReader(buf))
}

func request(args []string) interface{} {
	if len(args) == 1 {
		var ret interface{}
		err := json.Unmarshal([]byte(args[0]), &ret)
		if err == nil {
			return ret
		}
	}

	ret := make(map[string]interface{})

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		ret[parts[0]] = coerce(parts[1])
	}

	return ret
}

func coerce(s string) interface{} {
	var v interface{}

	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		v = s
	}

	return v
}
