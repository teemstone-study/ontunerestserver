package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"ontunerestserver/data"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPages_CRUD(t *testing.T) {

	assert := assert.New(t)
	ah := MakeHandler(CreateDBInfo("../server.env").GetDBConnString(), nil, nil)
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	//Delete 테스트
	req, _ := http.NewRequest("DELETE", ts.URL+"/pages", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	//post isnert 2 Pages
	resp, err = http.Post(ts.URL+"/pages", "application/json",
		strings.NewReader(`[{"index": 0,"contents": "first,second,first"},{"index": 2,"contents": "first,second,first,second"}]`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	//get pages 2 Pages
	resp, err = http.Get(ts.URL + "/pages")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	pages := []*data.Page{}
	err = json.NewDecoder(resp.Body).Decode(&pages)
	assert.NoError(err)
	assert.Equal(len(pages), 2)
	for _, p := range pages { //첫번째 인자가 index, 두번째가 value
		if p.Index == 0 {
			assert.Equal("first,second,first", p.Contents)
		} else if p.Index == 2 {
			assert.Equal("first,second,first,second", p.Contents)
		} else {
			assert.Error(fmt.Errorf("testIndex should be 0 or 2"))
		}
	}

	//put upate 2 Pages
	req, _ = http.NewRequest("PUT", ts.URL+"/pages",
		strings.NewReader(`[{"index": 0,"contents": "first"},{"index": 2,"contents": "first,second"}]`)) //PUT 는 많이 사용하지 않아 따로 요청을 만들어 줘야한다. NewRequest
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	su := new(Success)
	err = json.NewDecoder(resp.Body).Decode(su)
	assert.NoError(err)
	assert.Equal(su.Success, true)

	//get pages 2 Pages
	// resp, err = http.Get(ts.URL + "/pages")
	// assert.NoError(err)
	// assert.Equal(http.StatusOK, resp.StatusCode)
	// pages = []*data.Page{}
	// err = json.NewDecoder(resp.Body).Decode(&pages)
	// assert.NoError(err)
	// assert.Equal(len(pages), 2)
	// for _, p := range pages {
	// 	if p.Index == 0 {
	// 		assert.Equal("first", p.Contents)
	// 	} else if p.Index == 2 {
	// 		assert.Equal("first,second", p.Contents)
	// 	} else {
	// 		assert.Error(fmt.Errorf("testIndex should be 0 or 2"))
	// 	}
	// }

	//get pages 2 Pages
	resp, err = http.Get(ts.URL + "/pages/0")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	page := new(data.Page)
	err = json.NewDecoder(resp.Body).Decode(page)
	assert.NoError(err)
	assert.Equal(page.Contents, "first")

	//get pages 2 Pages
	resp, err = http.Get(ts.URL + "/pages/2")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	page = new(data.Page)
	err = json.NewDecoder(resp.Body).Decode(page)
	assert.NoError(err)
	assert.Equal(page.Contents, "first,second")

	//get pages 2 Pages
	resp, err = http.Get(ts.URL + "/pages/1")
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
}
