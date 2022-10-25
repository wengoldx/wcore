// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/10/20   jidi           New version
// -------------------------------------------------------------------

package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/wengoldx/wing/invar"
)

func respError(res *esapi.Response) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read resp body err:%v", err)
	}
	resp := &ErrorResp{}
	if err := json.Unmarshal(body, resp); err != nil {
		return fmt.Errorf("json unmarshal resp body err:%v", err)
	}
	reason := errors.New(resp.ErrorReason.Reason)
	return reason
}

// -------------------------------------------------------------

/*
create the new index and setting index mapping, if index is exist, the method will return error
mapping = `
{
	"mappings": {
		"properties": {
			"title": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_smart"
			},
		}
	}
}`
*/
func (e *ESClient) CreateIndexMapping(index, mapping string) error {
	res, err := e.Conn.Indices.Create(
		index,
		e.Conn.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("create index err:%v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return respError(res)
	}
	return nil
}

/*
update the index mapping, for exmple: add new filed or update filed setting
mapping := `
{
	"properties": {
		"title": {
		"type": "text",
		"analyzer": "ik_max_word",
		"search_analyzer": "ik_smart"
		}
	}
}`
*/
func (e *ESClient) UpdateIndexMapping(index []string, mapping string) error {
	res, err := e.Conn.Indices.PutMapping(
		index,
		strings.NewReader(mapping),
	)
	if err != nil {
		return fmt.Errorf("update index mapping err:%v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return respError(res)
	}
	return nil
}

/*
create new doc, if the index not exist, will auto create index mapping
*/
func (e *ESClient) CreateIndexDoc(index string, doc interface{}, docID ...string) error {
	id := ""
	if len(docID) > 0 {
		id = docID[0]
	}
	byteDoc, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("json index doc err:%v", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(byteDoc),
		Refresh:    "wait_for",
	}

	res, err := req.Do(context.Background(), e.Conn)
	if err != nil {
		return fmt.Errorf("create index doc err:%v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return respError(res)
	}
	return nil
}

/*
update the specified fields in the index
doc := `
{
	"doc": {
		"fields":"value"
	}
}
`
*/
func (e *ESClient) UpdateIndexDoc(index, docID, doc string) error {
	res, err := e.Conn.Update(
		index,
		docID,
		strings.NewReader(doc),
		e.Conn.Update.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("update index doc err:%v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return respError(res)
	}
	return nil
}

/*
search doc by query in the index, and set page, limit
page default 0 and limit default 10
*/
func (e *ESClient) SearchIndex(index, query string, page int, limit ...int) (*Response, error) {
	size := 10
	if len(limit) > 0 {
		size = limit[0]
	}
	res, err := e.Conn.Search(
		e.Conn.Search.WithIndex(index),
		e.Conn.Search.WithSize(size),
		e.Conn.Search.WithFrom(page),
		e.Conn.Search.WithBody(strings.NewReader(query)),
	)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("search index err:%v", err)
	}

	if res.IsError() {
		return nil, respError(res)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read search resp err:%v", err)
	}
	resp := &Response{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, fmt.Errorf("json unmarshal resp err:%v", err)
	}
	return resp, nil
}

// batch check indexs whether exist
func (e *ESClient) IsExistIndex(index []string) bool {
	exist, _ := e.Conn.Indices.Get(index)
	return exist.StatusCode == invar.StatusOK
}
