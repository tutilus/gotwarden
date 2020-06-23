package util

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func MarshalArray(array []interface{}) []byte {
	if len(array) > 0 {
		raw, err := json.Marshal(array)
		if err != nil {
			log.Printf("Marshal of struct failed")
		}
		return raw
	}
	return nil
}

func MarshalObject(v interface{}) []byte {
	if v == nil {
		return nil
	} else {
		raw, err := json.Marshal(v)
		if err != nil {
			log.Printf("Marshal of struct failed")
		}
		return raw
	}
}

func UnmarshalArray(object []byte) []interface{} {
	if object == nil {
		return nil
	} else {
		var raw []interface{}
		err := json.Unmarshal(object, &raw)
		if err != nil {
			log.Printf("Unmarshal of object failed")
		}
		return raw
	}
}

func UnmarshalObject(object []byte) *interface{} {
	if object == nil {
		return nil
	} else {
		var raw interface{}
		err := json.Unmarshal(object, &raw)
		if err != nil {
			log.Printf("Unmarshal of object failed")
		}
		return &raw
	}
}

func Trace(c *gin.Context) {
	log.Printf("Trace de la requête %s", c.Request.RequestURI)
	body, err := ioutil.ReadAll(c.Request.Body)

	log.Printf("Header Request %s", c.Request.Header)
	if err != nil {
		log.Printf("Read body error %s", err)
	} else {
		log.Printf("Body Request %s", body)
	}
}
