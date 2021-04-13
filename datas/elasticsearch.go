package datas

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	AggsAttrTypeValue   = "value"
	AggsAttrTypeBuckets = "buckets"
)

// var (
// 	indexMappingMap map[string]map[string]string
// )

// EsConfig elasticsearch连接配置
type EsConfig struct {
	Address []string
}

type esRespMapper map[string]interface{}

var NoDataError = errors.New("no data found")

// EsClient elasticsearch客户端的封装
type EsClient struct {
	*elasticsearch.Client
}

var esCli *EsClient

// InitESCli 初始化elasticsearch客户端
func InitESCli(conf *EsConfig) {
	var err error
	cfg := elasticsearch.Config{
		Addresses: conf.Address,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("create elasticsearch client with err: %s", err)
	}

	res, err := client.Info()
	if err != nil {
		log.Fatalf("get elasticsearch client info with err: %s", err)
	}

	if res.IsError() {
		log.Fatalf("get elasticsearch client info with err:: %s", res.String())
	}

	res.Body.Close()

	esCli = &EsClient{client}
}

// GetESClient 获取elasticsearch客户端的指针
func GetESClient() (*EsClient, error) {
	return esCli, nil
}

// GetRespMap 从elasticsearch的Response中获取esRespMap对象
func GetRespMap(r *esapi.Response) (esRespMapper, error) {
	respMap := make(esRespMapper)
	err := json.NewDecoder(r.Body).Decode(&respMap)
	r.Body.Close()
	return respMap, errors.Wrap(err, "json decode")
}

// GetDatas 从esRespMap对象中获取命中文档数和文档信息
func (respMap esRespMapper) GetDatas() (int64, []map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	hit := respMap["hits"].(map[string]interface{})
	hitCount := int64(hit["total"].(map[string]interface{})["value"].(float64))
	if hitCount < 1 {
		return hitCount, result, nil
	}
	hits := hit["hits"].([]interface{})
	for _, source := range hits {
		result = append(result, source.(map[string]interface{})["_source"].(map[string]interface{}))
	}
	return hitCount, result, nil
}

// GetAggs 从esRespMap对象中获取aggregations中指定key的值
func (respMap esRespMapper) GetAggs(attrType string, aggsKeys ...string) ([]interface{}, error) {
	results := make([]interface{}, len(aggsKeys))
	for index, aggsKey := range aggsKeys {
		aggsMap := respMap["aggregations"].(map[string]interface{})
		if aggsMap == nil {
			return nil, fmt.Errorf("no option of 'aggregations' found")
		}
		if aggsMap[aggsKey] == nil {
			return nil, fmt.Errorf("no key of '%s' found", aggsKey)
		}
		aggsVal := aggsMap[aggsKey].(map[string]interface{})
		if aggsVal == nil {
			return nil, fmt.Errorf("no key of '%s' found", aggsKey)
		}
		results[index] = aggsVal[attrType]
	}

	return results, nil
}

// GetDataByID 从esRespMap对象中获取根据_id查询的数据
func (respMap esRespMapper) GetDataByID() (bool, map[string]interface{}) {
	found := respMap["found"].(bool)
	if !found {
		return found, nil
	}
	return found, respMap["_source"].(map[string]interface{})
}
