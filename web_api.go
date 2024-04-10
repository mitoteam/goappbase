package goappbase

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type (
	ApiRequest struct {
		inData  map[string]interface{}
		outData map[string]interface{}
		session sessions.Session

		context *gin.Context
	}

	ApiRequestHandler func(r *ApiRequest) error
)

func newApiRequest(c *gin.Context) (*ApiRequest, error) {
	r := &ApiRequest{
		inData:  make(map[string]interface{}),
		outData: make(map[string]interface{}),
		context: c,
	}

	//prepare session
	r.session = sessions.Default(c)
	r.session.Options(sessions.Options{
		MaxAge: 24 * 3600,
		Path:   "/",
	})

	//prepare input data
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1048576))
	if err != nil {
		return nil, err
	}
	//log.Println(string(body))

	if json.Valid(body) {
		if err := json.Unmarshal(body, &r.inData); err != nil {
			return nil, err
		}
	}
	//log.Println(r.inData)

	return r, nil
}

func (r *ApiRequest) GetInData(name string) string {
	if value, ok := r.inData[name]; ok {
		if _, ok := value.(string); ok {
			return value.(string)
		} else {
			return fmt.Sprintf("%v", value)
		}
	} else {
		return ""
	}
}

func (r *ApiRequest) GetInDataInt(name string, default_value int) int {
	if value, ok := r.inData[name]; ok {
		//log.Println(reflect.TypeOf(value))

		if _, ok := value.(int); ok {
			return value.(int)
		}

		if _, ok := value.(float64); ok {
			return int(value.(float64))
		}
	}

	return default_value
}

func (r *ApiRequest) GetOutData(name string) string {
	if value, ok := r.outData[name]; ok {
		return value.(string)
	} else {
		return ""
	}
}

func (r *ApiRequest) SetOutData(name string, value interface{}) {
	r.outData[name] = value
}

func (r *ApiRequest) setStatus(status, message string) {
	r.SetOutData("status", status)
	r.SetOutData("message", message)
}

func (r *ApiRequest) SetOkStatus(message string) {
	r.setStatus("ok", message)
}

func (r *ApiRequest) SetErrorStatus(message string) {
	r.setStatus("error", message)
}

func (r *ApiRequest) Session() sessions.Session {
	return r.session
}

func (r *ApiRequest) SessionClear() {
	r.session.Clear()

	r.session.Options(sessions.Options{
		MaxAge: -1, //remove immediately
		Path:   "/",
	})

	r.session.Save()
}

func (r *ApiRequest) SessionGet(key string) any {
	return r.session.Get(key)
}

func (r *ApiRequest) SessionSet(key string, value any) {
	r.session.Set(key, value)
}
