package commons

import (
	"net/rpc"

	"github.com/klippa-app/go-pdfium/pdfium/requests"
	"github.com/klippa-app/go-pdfium/pdfium/responses"

	"github.com/hashicorp/go-plugin"
)

type Pdfium interface {
	Ping() (string, error)
	OpenDocument(*requests.OpenDocument) error
	GetPageCount(*requests.GetPageCount) (*responses.GetPageCount, error)
	GetPageText(*requests.GetPageText) (*responses.GetPageText, error)
	GetPageTextStructured(*requests.GetPageTextStructured) (*responses.GetPageTextStructured, error)
	RenderPageInDPI(*requests.RenderPageInDPI) (*responses.RenderPage, error)
	RenderPagesInDPI(*requests.RenderPagesInDPI) (*responses.RenderPages, error)
	RenderPageInPixels(*requests.RenderPageInPixels) (*responses.RenderPage, error)
	RenderPagesInPixels(*requests.RenderPagesInPixels) (*responses.RenderPages, error)
	GetPageSize(*requests.GetPageSize) (*responses.GetPageSize, error)
	GetPageSizeInPixels(*requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error)
	RenderToFileRequest(*requests.RenderToFileRequest) (*responses.RenderToFileRequest, error)
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

func (g *PdfiumRPC) OpenDocument(request *requests.OpenDocument) error {
	err := g.client.Call("Plugin.OpenDocument", request, new(interface{}))
	if err != nil {
		return err
	}

	return nil
}

func (g *PdfiumRPC) GetPageCount(request *requests.GetPageCount) (*responses.GetPageCount, error) {
	resp := &responses.GetPageCount{}
	err := g.client.Call("Plugin.GetPageCount", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageText(request *requests.GetPageText) (*responses.GetPageText, error) {
	resp := &responses.GetPageText{}
	err := g.client.Call("Plugin.GetPageText", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageTextStructured(request *requests.GetPageTextStructured) (*responses.GetPageTextStructured, error) {
	resp := &responses.GetPageTextStructured{}
	err := g.client.Call("Plugin.GetPageTextStructured", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPageInDPI(request *requests.RenderPageInDPI) (*responses.RenderPage, error) {
	resp := &responses.RenderPage{}
	err := g.client.Call("Plugin.RenderPageInDPI", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPagesInDPI(request *requests.RenderPagesInDPI) (*responses.RenderPages, error) {
	resp := &responses.RenderPages{}
	err := g.client.Call("Plugin.RenderPagesInDPI", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPageInPixels(request *requests.RenderPageInPixels) (*responses.RenderPage, error) {
	resp := &responses.RenderPage{}
	err := g.client.Call("Plugin.RenderPageInPixels", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderPagesInPixels(request *requests.RenderPagesInPixels) (*responses.RenderPages, error) {
	resp := &responses.RenderPages{}
	err := g.client.Call("Plugin.RenderPagesInPixels", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageSize(request *requests.GetPageSize) (*responses.GetPageSize, error) {
	resp := &responses.GetPageSize{}
	err := g.client.Call("Plugin.GetPageSize", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) GetPageSizeInPixels(request *requests.GetPageSizeInPixels) (*responses.GetPageSizeInPixels, error) {
	resp := &responses.GetPageSizeInPixels{}
	err := g.client.Call("Plugin.GetPageSizeInPixels", request, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *PdfiumRPC) RenderToFileRequest(request *requests.RenderToFileRequest) (*responses.RenderToFileRequest, error) {
	resp := &responses.RenderToFileRequest{}
	err := g.client.Call("Plugin.RenderToFileRequest", request, resp)
	if err != nil {
		return nil, err
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

func (s *PdfiumRPCServer) OpenDocument(request *requests.OpenDocument, resp *interface{}) error {
	var err error
	err = s.Impl.OpenDocument(request)
	if err != nil {
		return err
	}
	return nil
}

func (s *PdfiumRPCServer) GetPageCount(request *requests.GetPageCount, resp *responses.GetPageCount) error {
	var err error
	implResp, err := s.Impl.GetPageCount(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) GetPageText(request *requests.GetPageText, resp *responses.GetPageText) error {
	var err error
	implResp, err := s.Impl.GetPageText(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) GetPageTextStructured(request *requests.GetPageTextStructured, resp *responses.GetPageTextStructured) error {
	var err error
	implResp, err := s.Impl.GetPageTextStructured(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) RenderPageInDPI(request *requests.RenderPageInDPI, resp *responses.RenderPage) error {
	var err error
	implResp, err := s.Impl.RenderPageInDPI(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) RenderPagesInDPI(request *requests.RenderPagesInDPI, resp *responses.RenderPages) error {
	var err error
	implResp, err := s.Impl.RenderPagesInDPI(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) RenderPageInPixels(request *requests.RenderPageInPixels, resp *responses.RenderPage) error {
	var err error
	implResp, err := s.Impl.RenderPageInPixels(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) RenderPagesInPixels(request *requests.RenderPagesInPixels, resp *responses.RenderPages) error {
	var err error
	implResp, err := s.Impl.RenderPagesInPixels(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) GetPageSize(request *requests.GetPageSize, resp *responses.GetPageSize) error {
	var err error
	implResp, err := s.Impl.GetPageSize(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) GetPageSizeInPixels(request *requests.GetPageSizeInPixels, resp *responses.GetPageSizeInPixels) error {
	var err error
	implResp, err := s.Impl.GetPageSizeInPixels(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

	return nil
}

func (s *PdfiumRPCServer) RenderToFileRequest(request *requests.RenderToFileRequest, resp *responses.RenderToFileRequest) error {
	var err error
	implResp, err := s.Impl.RenderToFileRequest(request)
	if err != nil {
		return err
	}

	// Overwrite the target address of resp to the target address of implResp.
	*resp = *implResp

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
