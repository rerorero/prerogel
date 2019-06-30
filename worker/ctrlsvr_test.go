package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rerorero/prerogel/command"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/sirupsen/logrus/hooks/test"
)

func testListener(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	return ln
}

func Test_newCtrlServer(t *testing.T) {
	var coordinatorFunc actor.ActorFunc
	coordinator := actor.EmptyRootContext.Spawn(actor.PropsFromFunc(func(c actor.Context) {
		if coordinatorFunc != nil {
			coordinatorFunc(c)
		}
	}))
	l, _ := test.NewNullLogger()
	sut := newCtrlServer(coordinator, l)
	ln := testListener(t)
	go func() {
		if err := sut.serve(ln); err != nil {
			t.Fatal(err)
		}
	}()
	defer sut.shutdown(context.TODO())
	u, err := url.Parse("http://" + ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	type mock struct {
		coordinator actor.ActorFunc
	}
	type args struct {
		method string
		path   string
		req    interface{}
	}
	tests := []struct {
		name    string
		mock    mock
		args    args
		wantRes func(r *http.Response)
	}{
		{
			name: "stats ok",
			mock: mock{
				coordinator: func(c actor.Context) {
					if _, ok := c.Message().(*command.CoordinatorStats); ok {
						c.Respond(&command.CoordinatorStatsAck{
							SuperStep:        1,
							NrOfActiveVertex: 2,
							NrOfSentMessages: 3,
							State:            "foo",
						})
					}
				},
			},
			args: args{
				method: http.MethodGet,
				path:   APIPathStats,
				req:    nil,
			},
			wantRes: func(r *http.Response) {
				var ack command.CoordinatorStatsAck
				if err := json.NewDecoder(r.Body).Decode(&ack); err != nil {
					t.Fatal(err)
				}
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
				if diff := cmp.Diff(ack, command.CoordinatorStatsAck{
					SuperStep:        1,
					NrOfActiveVertex: 2,
					NrOfSentMessages: 3,
					State:            "foo",
				}); diff != "" {
					t.Fatalf("unmatch %s", diff)
				}
			},
		},
		{
			name: "load vartex ok",
			mock: mock{
				coordinator: func(c actor.Context) {
					if cmd, ok := c.Message().(*command.LoadVertex); ok {
						if cmd.VertexId != "123" {
							t.Fatal("unexpected id")
						}
						c.Respond(&command.LoadVertexAck{
							VertexId: "123",
							Error:    "",
						})
					}
				},
			},
			args: args{
				method: http.MethodPost,
				path:   APIPathLoadVertex,
				req: &command.LoadVertex{
					VertexId: "123",
				},
			},
			wantRes: func(r *http.Response) {
				var ack command.LoadVertexAck
				if err := json.NewDecoder(r.Body).Decode(&ack); err != nil {
					t.Fatal(err)
				}
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
				if ack.VertexId != "123" {
					t.Fatal("not match")
				}
			},
		},
		{
			name: "load partition vertices ok",
			mock: mock{
				coordinator: func(c actor.Context) {},
			},
			args: args{
				method: http.MethodPost,
				path:   APIPathLoadPartitionVertices,
				req:    nil,
			},
			wantRes: func(r *http.Response) {
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
			},
		},
		{
			name: "start super step ok",
			mock: mock{
				coordinator: func(c actor.Context) {},
			},
			args: args{
				method: http.MethodPost,
				path:   APIPathStartSuperStep,
				req:    nil,
			},
			wantRes: func(r *http.Response) {
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
			},
		},
		{
			name: "agg ok",
			mock: mock{
				coordinator: func(c actor.Context) {
					if _, ok := c.Message().(*command.ShowAggregatedValue); ok {
						c.Respond(&command.ShowAggregatedValueAck{
							AggregatedValues: map[string]string{
								"foo": "bar",
							},
						})
					}
				},
			},
			args: args{
				method: http.MethodPost,
				path:   APIPathShowAggregatedValue,
				req:    nil,
			},
			wantRes: func(r *http.Response) {
				var ack command.ShowAggregatedValueAck
				if err := json.NewDecoder(r.Body).Decode(&ack); err != nil {
					t.Fatal(err)
				}
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
				if ack.AggregatedValues["foo"] != "bar" {
					t.Fatal("not match")
				}
			},
		},
		{
			name: "agg ok",
			mock: mock{
				coordinator: func(c actor.Context) {
					if _, ok := c.Message().(*command.Shutdown); ok {
						c.Respond(&command.ShutdownAck{})
					}
				},
			},
			args: args{
				method: http.MethodPost,
				path:   APIPathShutdown,
				req:    &command.Shutdown{},
			},
			wantRes: func(r *http.Response) {
				if r.StatusCode != http.StatusOK {
					t.Fatal("not ok")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coordinatorFunc = tt.mock.coordinator
			url2 := *u
			url2.Path = path.Join(tt.args.path)

			v, err := json.Marshal(tt.args.req)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(tt.args.method, url2.String(), bytes.NewBuffer(v))
			if err != nil {
				t.Fatal(err)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			tt.wantRes(res)
		})
	}
}
