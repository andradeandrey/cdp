package rpcc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

type testServer struct {
	srv    *httptest.Server
	wsConn *websocket.Conn // Set after Dial.
	conn   *Conn
}

func (ts *testServer) Close() error {
	defer ts.srv.Close()
	return ts.conn.Close()
}

func newTestServer(t *testing.T, respond func(*websocket.Conn, *rpcRequest) error) *testServer {
	// Timeouts to prevent tests from running forever.
	timeout := 5 * time.Second

	var err error
	ts := &testServer{}
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: timeout,
	}

	setupDone := make(chan struct{})
	ts.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, r.Header)
		if err != nil {
			t.Fatal(err)
		}
		ts.wsConn = conn
		close(setupDone)
		defer conn.Close()

		for {
			conn.SetReadDeadline(time.Now().Add(timeout))
			conn.SetWriteDeadline(time.Now().Add(timeout))

			var req rpcRequest
			if err := conn.ReadJSON(&req); err != nil {
				break
			}
			if err := respond(conn, &req); err != nil {
				break
			}
		}
	}))

	ts.conn, err = Dial("ws" + strings.TrimPrefix(ts.srv.URL, "http"))
	if err != nil {
		t.Fatal(err)
	}

	<-setupDone
	return ts
}

func TestConn_Invoke(t *testing.T) {
	responses := []string{
		"hello",
		"world",
	}
	srv := newTestServer(t, func(conn *websocket.Conn, req *rpcRequest) error {
		resp := rpcResponse{
			ID:     req.ID,
			Result: []byte(fmt.Sprintf("%q", responses[int(req.ID)-1])),
		}
		return conn.WriteJSON(&resp)
	})
	defer srv.Close()

	var reply string
	err := Invoke(nil, "test.Hello", nil, &reply, srv.conn)
	if err != nil {
		t.Error(err)
	}

	if reply != "hello" {
		t.Errorf("test.Hello: got reply %q, want %q", reply, "hello")
	}

	err = Invoke(nil, "test.World", nil, &reply, srv.conn)
	if err != nil {
		t.Error(err)
	}

	if reply != "world" {
		t.Errorf("test.World: got reply %q, want %q", reply, "world")
	}
}

func TestConn_InvokeError(t *testing.T) {
	want := "bad request"

	srv := newTestServer(t, func(conn *websocket.Conn, req *rpcRequest) error {
		resp := rpcResponse{
			ID:    req.ID,
			Error: &rpcError{Message: want},
		}
		return conn.WriteJSON(&resp)
	})
	defer srv.Close()

	switch err := Invoke(nil, "test.Hello", nil, nil, srv.conn).(type) {
	case *rpcError:
		if err.Message != want {
			t.Errorf("Invoke err.Message: got %q, want %q", err.Message, want)
		}
	default:
		t.Errorf("Invoke: want *rpcError, got %#v", err)
	}
}

func TestConn_DecodeError(t *testing.T) {
	srv := newTestServer(t, func(conn *websocket.Conn, req *rpcRequest) error {
		msg := fmt.Sprintf(`{"id": %d, "result": {}}`, req.ID)
		w, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			t.Fatal(err)
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			t.Fatal(err)
		}
		w.Close()
		return nil
	})
	defer srv.Close()

	var reply string
	err := Invoke(nil, "test.DecodeError", nil, &reply, srv.conn)
	if err == nil || !strings.HasPrefix(err.Error(), "rpcc: decoding") {
		t.Errorf("test.DecodeError: got %v, want error with %v", err, "rpcc: decoding")
	}
}

type badEncoder struct {
	ch  chan struct{}
	err error
}

func (enc *badEncoder) Decode(v interface{}) error { return nil }
func (enc *badEncoder) Encode(v interface{}) error { return enc.err }

func TestConn_EncodeFailed(t *testing.T) {
	enc := &badEncoder{err: errors.New("fail"), ch: make(chan struct{})}
	conn := &Conn{
		ctx:     context.Background(),
		pending: make(map[uint64]*rpcCall),
		codec:   enc,
	}

	err := Invoke(nil, "test.Hello", nil, nil, conn)
	if err != enc.err {
		t.Errorf("Encode: got %v, want %v", err, enc.err)
	}
}

func TestConn_Notify(t *testing.T) {
	srv := newTestServer(t, nil)
	defer srv.Close()

	s, err := NewStream(nil, "test.Notify", srv.conn)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	go func() {
		resp := rpcResponse{
			Method: "test.Notify",
			Params: []byte(`"hello"`),
		}
		if err := srv.wsConn.WriteJSON(&resp); err != nil {
			t.Fatal(err)
		}
	}()

	var reply string
	if err = s.RecvMsg(&reply); err != nil {
		t.Error(err)
	}

	want := "hello"
	if reply != want {
		t.Errorf("test.Notify reply: got %q, want %q", reply, want)
	}

	s.Close()
	if err = s.RecvMsg(nil); err == nil {
		t.Error("test.Notify read after closed: want error, got nil")
	}
}

func TestConn_RemoteDisconnected(t *testing.T) {
	srv := newTestServer(t, nil)
	defer srv.Close()

	srv.wsConn.Close()
	err := Invoke(nil, "test.Hello", nil, nil, srv.conn)
	wsErr, ok := err.(*websocket.CloseError)
	if !(ok && wsErr.Code == websocket.CloseAbnormalClosure) && err != ErrConnClosing {
		t.Errorf("Invoke error: got %v, want websocket.CloseError or ErrConnClosing", err)
	}
}

func TestDialContext_DeadlineExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	_, err := DialContext(ctx, "")

	// Should return deadline even when dial address is bad.
	if err != context.DeadlineExceeded {
		t.Errorf("DialContext: got %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestMain(m *testing.M) {
	enableDebug = true
	os.Exit(m.Run())
}