package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type PdfiumRPC struct{ client *rpc.Client }

func (g *PdfiumRPC) Ping() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Ping", new(interface{}), &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (g *PdfiumRPC) Close() error {
	err := g.client.Call("Plugin.Close", new(interface{}), new(interface{}))
	if err != nil {
		return err
	}

	return nil
}

type PdfiumRPCServer struct {
	Impl Pdfium
}

func (s *PdfiumRPCServer) Ping(args interface{}, resp *string) error {
	var err error
	*resp, err = s.Impl.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *PdfiumRPCServer) Close(args interface{}, resp *interface{}) error {
	var err error
	err = s.Impl.Close()
	if err != nil {
		return err
	}
	return nil
}

type PdfiumPlugin struct {
	Impl Pdfium
}

func (p *PdfiumPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PdfiumRPCServer{Impl: p.Impl}, nil
}

func (PdfiumPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PdfiumRPC{client: c}, nil
}
