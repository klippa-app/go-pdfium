package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type PdfiumRPC struct{ client *rpc.Client }

func (g *PdfiumRPC) Ping() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Ping", new(any), &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (g *PdfiumRPC) Close() error {
	err := g.client.Call("Plugin.Close", new(any), new(any))
	if err != nil {
		return err
	}

	return nil
}

type PdfiumRPCServer struct {
	Impl Pdfium
}

func (s *PdfiumRPCServer) Ping(args any, resp *string) error {
	var err error
	*resp, err = s.Impl.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *PdfiumRPCServer) Close(args any, resp *any) error {
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

func (p *PdfiumPlugin) Server(*plugin.MuxBroker) (any, error) {
	return &PdfiumRPCServer{Impl: p.Impl}, nil
}

func (PdfiumPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &PdfiumRPC{client: c}, nil
}
