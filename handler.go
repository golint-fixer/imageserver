package imageserver

import(
	"log"
)

// Handler handles an Image and returns an Image.
type Handler interface {
	Handle(*Image, Params) (*Image, error)
}

// HandlerFunc is a Handler func.
type HandlerFunc func(*Image, Params) (*Image, error)

// Handle implements Handler.
func (f HandlerFunc) Handle(im *Image, params Params) (*Image, error) {
	return f(im, params)
}

// HandlerServer is a Server implementation that calls a Handler.
type HandlerServer struct {
	Server
	Handler Handler
}

// Get implements Server.
func (srv *HandlerServer) Get(params Params) (*Image, error) {
	log.Println("Image::start")
	im, err := srv.Server.Get(params)
	log.Println("Image::end")
	
	if err != nil {
		return nil, err
	}
	
	log.Println("Process::start")
	im, err = srv.Handler.Handle(im, params)
	log.Println("Process::end")
	
	if err != nil {
		return nil, err
	}
	return im, nil
}
