package objects

import (
	"awesomeProject1/store/apiserver"
	"awesomeProject1/store/objectstream"
	"fmt"
)

func putStream(object string) (*objectstream.PutStream, error) {
	server := apiserver.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	return objectstream.NewPutStream(server, object), nil
}
