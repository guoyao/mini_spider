package media

import (
	"io/ioutil"
	"net/http"

	"mini_spider/util"
)

type Webpage struct {
	url      string
	request  *http.Request
	response *http.Response
	content  []byte
	depth    uint
	children []*Webpage
}

func NewWebpage(req *http.Request, resp *http.Response) *Webpage {
	return &Webpage{request: req, response: resp}
}

func (w *Webpage) GetName() string {
	return util.URLEncode(w.request.URL.String())
}

func (w *Webpage) GetContent() ([]byte, error) {
	if w.content != nil {
		return w.content, nil
	}

	if w.response != nil {
		content, err := ioutil.ReadAll(w.response.Body)
		if err != nil {
			return nil, err
		}

		w.content = content

		return w.content, nil
	}

	return []byte(""), nil
}
