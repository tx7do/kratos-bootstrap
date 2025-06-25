package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"

	elasticsearchV9 "github.com/elastic/go-elasticsearch/v9"
	esapiV9 "github.com/elastic/go-elasticsearch/v9/esapi"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	pagination "github.com/tx7do/kratos-bootstrap/api/gen/go/pagination/v1"
)

type Client struct {
	cli   *elasticsearchV9.Client
	log   *log.Helper
	codec encoding.Codec
}

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*Client, error) {
	c := &Client{
		log:   log.NewHelper(log.With(logger, "module", "elasticsearch-client")),
		codec: encoding.GetCodec("json"),
	}

	if err := c.createESClient(cfg); err != nil {
		return nil, err
	}

	return c, nil
}

// createESClient 创建Elasticsearch客户端
func (c *Client) createESClient(cfg *conf.Bootstrap) error {
	if cfg.Data == nil || cfg.Data.ElasticSearch == nil {
		return nil // No Elasticsearch configuration provided
	}

	cli, err := elasticsearchV9.NewClient(
		elasticsearchV9.Config{
			Addresses: cfg.Data.ElasticSearch.GetAddresses(),
			Username:  cfg.Data.ElasticSearch.GetUsername(),
			Password:  cfg.Data.ElasticSearch.GetPassword(),
		},
	)
	if err != nil {
		c.log.Errorf("failed to create elasticsearch client: %v", err)
		return err
	}

	c.cli = cli

	return nil
}

func (c *Client) Close() {

}

// CheckConnectStatus 检查Elasticsearch连接
func (c *Client) CheckConnectStatus() bool {
	if c.cli == nil {
		return false
	}

	resp, err := c.cli.Info()
	if err != nil {
		c.log.Errorf("failed to connect to elasticsearch: %v", err)
		return false
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			c.log.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.IsError() {
		c.log.Errorf("Error: %s", resp.String())
		return false
	}

	var r map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
		return false
	}

	c.log.Infof("Client Version: %s", elasticsearchV9.Version)
	c.log.Infof("Server Version: %s", r["version"].(map[string]interface{})["number"])

	return true
}

// IndexExists 检查索引是否存在
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	resp, err := c.cli.Indices.Exists(
		[]string{indexName},
		c.cli.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		c.log.Errorf("failed to check if index exists: %v", err)
		return false, err
	}

	return !resp.IsError(), nil
}

// CreateIndex 创建一条索引
//
//	如果mapping为空("")则表示不创建模型
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping, settings string) error {
	exist, err := c.IndexExists(ctx, indexName)
	if err != nil {
		return err
	}
	if exist {
		return ErrIndexAlreadyExists
	}

	body, err := MergeOptions(mapping, settings)
	if err != nil {
		c.log.Errorf("failed to merge options: %v", err)
		return err
	}

	resp, err := c.cli.Indices.Create(
		indexName,
		c.cli.Indices.Create.WithContext(ctx),
		c.cli.Indices.Create.WithBody(bytes.NewReader([]byte(body))),
	)
	if err != nil {
		c.log.Errorf("failed to create index: %v", err)
		return err
	}

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			c.log.Errorf("failed to parse error message: %v", err)
			return err
		}

		c.log.Errorf("create index failed: %s", errResp.Error)

		return ErrCreateIndex
	}

	return nil
}

// DeleteIndex 删除一条索引
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	exist, err := c.IndexExists(ctx, indexName)
	if err != nil {
		c.log.Errorf("failed to check if index exists: %v", err)
		return err
	}
	if !exist {
		return ErrIndexNotFound
	}

	resp, err := c.cli.Indices.Delete(
		[]string{indexName},
		c.cli.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		c.log.Errorf("failed to delete index: %v", err)
		return err
	}

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			c.log.Errorf("failed to parse error message: %v", err)
			return err
		}

		c.log.Errorf("delete index failed: %s", errResp.Error.Reason)

		return ErrDeleteIndex
	}

	return nil
}

// DeleteDocument 删除一条数据
func (c *Client) DeleteDocument(ctx context.Context, indexName, id string) error {
	_, err := c.cli.Delete(
		indexName, id,
		c.cli.Delete.WithContext(ctx),
	)
	if err != nil {
		c.log.Errorf("failed to delete document: %v", err)
		return err
	}
	return nil
}

// InsertDocument 插入一条数据
func (c *Client) InsertDocument(ctx context.Context, indexName, id string, data interface{}) error {
	var err error

	var dataBytes []byte
	dataBytes, err = json.Marshal(data)
	if err != nil {
		c.log.Errorf("failed to marshal data: %v", err)
		return err
	}

	var resp *esapiV9.Response

	if id == "" {
		resp, err = c.cli.Index(
			indexName,
			bytes.NewReader(dataBytes),
			c.cli.Index.WithContext(ctx),
		)
	} else {
		resp, err = c.cli.Create(
			indexName, id,
			bytes.NewReader(dataBytes),
			c.cli.Create.WithContext(ctx),
		)
	}
	if err != nil {
		c.log.Errorf("failed to insert document: %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			c.log.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			c.log.Errorf("failed to parse error message: %v", err)
			return err
		}

		c.log.Errorf("insert data failed: %s", errResp.Error.Reason)

		return ErrInsertDocument
	}

	return nil
}

// BatchInsertDocument 批量插入数据
func (c *Client) BatchInsertDocument(ctx context.Context, indexName string, dataSet []interface{}) error {
	var buf bytes.Buffer
	for _, data := range dataSet {
		meta := []byte(`{"index":{}}` + "\n")
		dataBytes, err := json.Marshal(data)
		if err != nil {
			c.log.Errorf("failed to marshal data: %v", err)
			continue
		}
		dataBytes = append(dataBytes, "\n"...)
		buf.Grow(len(meta) + len(dataBytes))
		buf.Write(meta)
		buf.Write(dataBytes)
	}

	resp, err := c.cli.Bulk(
		bytes.NewReader(buf.Bytes()),
		c.cli.Bulk.WithContext(ctx),
		c.cli.Bulk.WithIndex(indexName),
	)
	if err != nil {
		c.log.Errorf("failed to perform bulk insert: %v", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			c.log.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			c.log.Errorf("failed to parse error message: %v", err)
			return err
		}

		c.log.Errorf("batch insert data failed: %s", errResp.Error.Reason)

		return ErrBatchInsertDocument
	}

	return nil
}

func (c *Client) UpdateDocument(ctx context.Context, indexName string, pk string, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		c.log.Errorf("failed to marshal data: %v", err)
		return err
	}

	_, err = c.cli.Update(
		indexName, pk,
		bytes.NewReader(data),
		c.cli.Update.WithContext(ctx),
	)
	if err != nil {
		c.log.Errorf("failed to update document: %v", err)
		return err
	}

	return nil
}

// GetDocument 查询数据
func (c *Client) GetDocument(
	ctx context.Context,
	indexName string,
	id string,
	sourceFields []string,
	out interface{},
) error {
	resp, err := c.cli.Get(
		indexName, id,
		c.cli.Get.WithContext(ctx),
		c.cli.Get.WithSource(sourceFields...), // 指定返回的字段
	)
	if err != nil {
		c.log.Errorf("failed to get document: %v", err)
		return err
	}

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			c.log.Errorf("failed to parse error message: %v", err)
			return err
		}

		if resp.StatusCode == 404 {
			c.log.Warnf("document not found: %s", errResp.Error.Reason)
			return ErrDocumentNotFound
		}

		c.log.Errorf("get document failed: %s", errResp.Error.Reason)

		return ErrGetDocument
	}

	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		c.log.Errorf("failed to decode document: %v", err)
		return err
	}

	return nil
}

func (c *Client) Search(
	ctx context.Context,
	indexName string,
	req *pagination.PagingRequest,
) (*SearchResult, error) {
	var query string
	ParseQueryString(req.GetQuery())

	sortBy := make(map[string]bool)

	pageSize := req.GetPageSize()
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}

	return c.search(ctx, indexName, query, nil, sortBy, int(req.GetPage()), int(pageSize))
}

// search 查询数据
//
// @param ctx 上下文
// @param indexName 索引名
// @param query 查询条件，例如：field1:value1 AND field2:value2
// @param sourceFields 指定返回的字段，传入nil表示返回所有字段
// @param sortBy 排序
// @param from 分页的页码
// @param pageSize 分页每页的行数
func (c *Client) search(
	ctx context.Context,
	indexName string,
	query string,
	sourceFields []string,
	sortBy map[string]bool,
	from, pageSize int,
) (*SearchResult, error) {
	var sorts []string
	for k, v := range sortBy {
		if v {
			sorts = append(sorts, k+":asc")
		} else {
			sorts = append(sorts, k+":desc")
		}
	}

	resp, err := c.cli.Search(
		c.cli.Search.WithContext(ctx),
		c.cli.Search.WithIndex(indexName),
		c.cli.Search.WithFrom(from),
		c.cli.Search.WithSize(pageSize),
		c.cli.Search.WithSort(sorts...),
		c.cli.Search.WithQuery(query),
		c.cli.Search.WithSource(sourceFields...), // 指定返回的字段
	)
	if err != nil {
		c.log.Errorf("failed to search documents: %v", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			c.log.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.IsError() {
		var errResp *ErrorResponse
		if errResp, err = ParseErrorMessage(resp.Body); err != nil {
			return nil, err
		}

		c.log.Errorf("search document failed: %s", errResp.Error.Reason)

		return nil, ErrSearchDocument
	}

	var searchResult SearchResult
	if err = json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		c.log.Errorf("failed to decode search result: %v", err)
		return nil, err
	}

	return &searchResult, nil
}
