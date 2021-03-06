// Package testdata provides test images.
package client

import (
	"fmt"
	"io/ioutil"
	//"path/filepath"
	//"runtime"
	
	"github.com/pierrre/imageserver"

	"log"
	"net/http"
	"net/url"
	"strings"
)

type ImgReq struct {
	path  string
	query url.Values	
}

var (
	client = &http.Client{}
	
	// Server is an Image Server that uses filename as source.
	Server = imageserver.Server(imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		queryString, err := params.GetString(imageserver.QueryKey)
		if err != nil {
			return nil, err
		}
		
		im, err := Get(queryString)
		if err != nil {
			return nil, &imageserver.ParamError{Param: imageserver.SourceParam, Message: err.Error()}
		}
		return im, nil
	}))
)

// Get returns an Image for a name.
func Get(name string) (*imageserver.Image, error) {
	log.Println("url: " + name)
	query, _ := url.Parse(name)
	log.Println("path: " + query.Path)
	log.Println(query.RawQuery)
	
	u := &url.URL{
		Scheme:   "http",
		Host:     "localhost:8081",
		Path:     "/" + query.Path,
		RawQuery: query.RawQuery,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}	

	format:= resp.Header.Get("Content-Type")
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}

	im := &imageserver.Image{
		Format: strings.TrimPrefix(format, "image/"),
		Data: data,
	}
	log.Println("Get Img from file: " + im.Format)
	log.Println(len(im.Data))	

	//im, ok := Images[name]
	//if !ok {
	//	return nil, fmt.Errorf("unknown image \"%s\"", name)
	//}
	return im, nil
}

type ImgError struct {  
	url  string 
}  

func (e ImgError) Error() string {  
	return fmt.Sprintf("unknown image \"%s\"", e.url)  
} 
