package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	sport = ":8088"
	lmux  = "/rr"
)

var (
	//pages    = 50
	//keywords = ""
	dlog = &log.Logger{}
)

func main() {
	f, err := os.Create("rr.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	dlog = log.New(f, "[Debug]", log.LstdFlags)

	http.HandleFunc(lmux, rrHandler)
	http.ListenAndServe(sport, nil)
}

func rrHandler(w http.ResponseWriter, r *http.Request) {
	_req := r.URL.Query()
	keywords := _req.Get("k")
	pages, _ := strconv.Atoi(_req.Get("p"))
	if pages == 0 {
		pages = 50
	}

	buf := &bytes.Buffer{}
	buf.WriteString("<html>\n<body>")
	aurl := "https://www.douban.com/group/shanghaizufang/discussion?start="
	//aurl := "http://cn.bing.com/?intlF="
	burl := "https://www.douban.com/group/shzf/discussion?start="
	curl := "https://www.douban.com/group/zufan/discussion?start="
	urls := []string{aurl, burl, curl}
	for _, v := range urls {
		if err := clawContents(v, buf, pages, keywords); err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	buf.WriteString("</body>\n</html>")
	w.Write(buf.Bytes())
}

func clawContents(urlbody string, buf *bytes.Buffer, pages int, keywords string) error {
	for i := 0; i < pages; i += 25 {
		fullurl := urlbody + strconv.Itoa(i)
		/*
			resp, err := newDocumentBySocks5(fullurl)
			if err != nil {
				return err
			}
			if resp == nil {
				continue
			}

			b, _ := ioutil.ReadAll(resp.Body)
			doc, err := goquery.NewDocumentFromResponse(resp)
		*/
		doc, err := goquery.NewDocument(fullurl)
		if err != nil {
			dlog.Println(fullurl, err)
			continue
		}
		doc.Find("tr").Each(func(i int, s *goquery.Selection) {
			h, err := s.Html()
			if err != nil {
				dlog.Println(fullurl, err)
			} else {
				if strings.Contains(h, "title") && strings.Contains(h, keywords) {
					buf.WriteString(h)
					buf.WriteString("</br>")
				}
			}
		})
	}
	return nil
}

func newDocumentBySocks5(fullurl string) (*http.Response, error) {
	return callBySocks5(nil, fullurl)
}
