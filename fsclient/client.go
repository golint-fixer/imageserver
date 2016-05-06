// Package testdata provides test images.
package fsclient

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
	fs "github.com/weedfs"
	"encoding/json"
)

type MetaData struct {
	FileType string
	FileID   string
}

var (
	dbURLKeys   = "http://localhost:9988/v1/buckets/meta/keys/"
	clientDB    = &http.Client{}
	clientFS    = fs.NewClient("localhost:9333")
	
	// Server is an Image Server that uses filename as source.
	Server = imageserver.Server(imageserver.ServerFunc(func(params imageserver.Params) (*imageserver.Image, error) {
		queryString, err := params.GetString(imageserver.QueryKey)
		if err != nil {
			return nil, err
		}
		//log.Println(queryString)
		im, err := Get(queryString)
		if err != nil {
			return nil, &imageserver.ParamError{Param: imageserver.SourceParam, Message: err.Error()}
		}
		return im, nil
	}))
)

// Get returns an Image for a name.
func Get(name string) (*imageserver.Image, error) {
	//log.Println("Get image from FS Server!")
	log.Println("query: " + name)
	
	query, _ := url.Parse(name)
	//log.Println("path: " + query.Path)
	//log.Println(query.RawQuery)

	dbURL := dbURLKeys + query.Path // image path
	var metaDataRead MetaData

	req, err := http.NewRequest("GET", dbURL, nil)
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}
	
	log.Println("DB::start")
	resp, err := clientDB.Do(req)
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}
	log.Println("DB::end")
	
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, ImgError{url: name}
	}

	err = json.Unmarshal(data, &metaDataRead)
	log.Printf("DB::metadata: %v\n", metaDataRead)

	log.Println("FS::start")
	dlReader, err := clientFS.Download(metaDataRead.FileID)
	data, err = ioutil.ReadAll(dlReader)
	log.Println("FS::end")
	
	im := &imageserver.Image{
		Format: strings.TrimPrefix(metaDataRead.FileType, "image/"),
		Data: data,
	}

	log.Printf("FS::image: %s, %d\n", im.Format, len(im.Data))	

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
