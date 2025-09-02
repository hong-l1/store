package objects

import (
	"awesomeProject1/store/apiserver"
	"awesomeProject1/store/objectstream"
	"fmt"
	"io"
)

func getStream(object string) (io.Reader, error) {
	server := apiserver.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
