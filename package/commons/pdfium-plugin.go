package commons

import (
	"image"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type OpenDocumentRequest struct {
	File *[]byte
}

type RenderPageRequest struct {
	Page int
	DPI  int
}

type RenderPageResponse struct {
	Image *image.RGBA
}

type GetPageSizeRequest struct {
	Page int
	DPI  int
}

type GetPageSizeResponse struct {
	Width  int
	Height int
}

type Pdfium interface {
	Ping() (string, error)
	OpenDocument(*OpenDocumentRequest) error
	GetPageCount() (int, error)
	RenderPage(*RenderPageRequest) (RenderPageResponse, error)
	GetPageSize(*GetPageSizeRequest) (GetPageSizeResponse, error)
	Close() error
}

type PdfiumRPC struct{ client *rpc.Client }

func (g *PdfiumRPC) Ping() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Ping", new(interface{}), &resp)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (g *PdfiumRPC) OpenDocument(request *OpenDocumentRequest) error {
	err := g.client.Call("Plugin.OpenDocument", request, new(interface{}))
	if err != nil {
		return err
	}

	return nil
}

func (g *PdfiumRPC) GetPageCount() (int, error) {
	var resp int
	err := g.client.Call("Plugin.GetPageCount", new(interface{}), &resp)
	if err != nil {
		return 0, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPage(request *RenderPageRequest) (RenderPageResponse, error) {
	var resp RenderPageResponse
	err := g.client.Call("Plugin.RenderPage", request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageSize(request *GetPageSizeRequest) (GetPageSizeResponse, error) {
	var resp GetPageSizeResponse
	err := g.client.Call("Plugin.GetPageSize", request, &resp)
	if err != nil {
		return resp, err
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

func (s *PdfiumRPCServer) OpenDocument(request *OpenDocumentRequest, resp *interface{}) error {
	var err error
	err = s.Impl.OpenDocument(request)
	if err != nil {
		return err
	}
	return nil
}

func (s *PdfiumRPCServer) GetPageCount(args interface{}, resp *int) error {
	var err error
	*resp, err = s.Impl.GetPageCount()
	if err != nil {
		return err
	}
	return nil
}

func (s *PdfiumRPCServer) RenderPage(request *RenderPageRequest, resp *RenderPageResponse) error {
	var err error
	*resp, err = s.Impl.RenderPage(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *PdfiumRPCServer) GetPageSize(request *GetPageSizeRequest, resp *GetPageSizeResponse) error {
	var err error
	*resp, err = s.Impl.GetPageSize(request)
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
