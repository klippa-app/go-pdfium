package commons

import (
	"image"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type OpenDocumentRequest struct {
	File *[]byte
}

type RenderPageInDPIRequest struct {
	Page int
	DPI  int
}

type RenderPageInPixelsRequest struct {
	Page   int
	Width  int
	Height int
}

type RenderPageResponse struct {
	Image *image.RGBA
}

type GetPageSizeRequest struct {
	Page int
}

type GetPageSizeInPixelsRequest struct {
	Page int
	DPI  int
}

type GetPageSizeResponse struct {
	Width  float64
	Height float64
}

type GetPageSizeInPixelsResponse struct {
	Width  int
	Height int
}

type GetPageTextRequest struct {
	Page int
}

type GetPageTextResponse struct {
	Text *string
}

type GetPageTextStructuredRequest struct {
	Page int
	Mode GetPageTextStructuredRequestMode
}

type GetPageTextStructuredRequestMode string

const (
	GetPageTextStructuredRequestModeChars GetPageTextStructuredRequestMode = "char"
	GetPageTextStructuredRequestModeRects GetPageTextStructuredRequestMode = "rect"
	GetPageTextStructuredRequestModeBoth  GetPageTextStructuredRequestMode = "both"
)

type GetPageTextStructuredResponseChar struct {
	Text   string
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
	Angle  float64
}

type GetPageTextStructuredResponseRect struct {
	Text   string
	Left   float64
	Top    float64
	Right  float64
	Bottom float64
}

type GetPageTextStructuredResponse struct {
	Chars []*GetPageTextStructuredResponseChar
	Rects []*GetPageTextStructuredResponseRect
}

type Pdfium interface {
	Ping() (string, error)
	OpenDocument(*OpenDocumentRequest) error
	GetPageCount() (int, error)
	GetPageText(*GetPageTextRequest) (GetPageTextResponse, error)
	GetPageTextStructured(*GetPageTextStructuredRequest) (GetPageTextStructuredResponse, error)
	RenderPageInDPI(*RenderPageInDPIRequest) (RenderPageResponse, error)
	RenderPageInPixels(*RenderPageInPixelsRequest) (RenderPageResponse, error)
	GetPageSize(*GetPageSizeRequest) (GetPageSizeResponse, error)
	GetPageSizeInPixels(*GetPageSizeInPixelsRequest) (GetPageSizeInPixelsResponse, error)
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

func (g *PdfiumRPC) GetPageText(request *GetPageTextRequest) (GetPageTextResponse, error) {
	var resp GetPageTextResponse
	err := g.client.Call("Plugin.GetPageText", request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageTextStructured(request *GetPageTextStructuredRequest) (GetPageTextStructuredResponse, error) {
	var resp GetPageTextStructuredResponse
	err := g.client.Call("Plugin.GetPageTextStructured", request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPageInDPI(request *RenderPageInDPIRequest) (RenderPageResponse, error) {
	var resp RenderPageResponse
	err := g.client.Call("Plugin.RenderPageInDPI", request, &resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPageInPixels(request *RenderPageInPixelsRequest) (RenderPageResponse, error) {
	var resp RenderPageResponse
	err := g.client.Call("Plugin.RenderPageInPixels", request, &resp)
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

func (g *PdfiumRPC) GetPageSizeInPixels(request *GetPageSizeInPixelsRequest) (GetPageSizeInPixelsResponse, error) {
	var resp GetPageSizeInPixelsResponse
	err := g.client.Call("Plugin.GetPageSizeInPixels", request, &resp)
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

func (s *PdfiumRPCServer) GetPageText(request *GetPageTextRequest, resp *GetPageTextResponse) error {
	var err error
	*resp, err = s.Impl.GetPageText(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *PdfiumRPCServer) GetPageTextStructured(request *GetPageTextStructuredRequest, resp *GetPageTextStructuredResponse) error {
	var err error
	*resp, err = s.Impl.GetPageTextStructured(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *PdfiumRPCServer) RenderPageInDPI(request *RenderPageInDPIRequest, resp *RenderPageResponse) error {
	var err error
	*resp, err = s.Impl.RenderPageInDPI(request)
	if err != nil {
		return err
	}

	return nil
}

func (s *PdfiumRPCServer) RenderPageInPixels(request *RenderPageInPixelsRequest, resp *RenderPageResponse) error {
	var err error
	*resp, err = s.Impl.RenderPageInPixels(request)
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

func (s *PdfiumRPCServer) GetPageSizeInPixels(request *GetPageSizeInPixelsRequest, resp *GetPageSizeInPixelsResponse) error {
	var err error
	*resp, err = s.Impl.GetPageSizeInPixels(request)
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
