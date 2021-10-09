package storage

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"gopkg.in/olivere/elastic.v5"
	"strings"
	"time"
)

var (
	esClient     *elastic.Client
	ESLogHandler *Elastic
)

type Elastic struct {
	esClient *elastic.Client
	esIndex  string
	esType   string
}

func InitElastic(cfg config.Configer) (err error) {
	esHosts := cfg.String("elasticsearch::hosts")
	esUser := cfg.String("elasticsearch::user")
	esSecret := cfg.String("elasticsearch::secret")
	if esHosts == "" {
		panic("elasticsearch config error")
	}
	logs.Info("elasticsearch,  hosts = %s, user=%s, secret=%s", esHosts, esUser, esSecret)

	hosts := strings.Split(esHosts, ",")

	esClient, err = elastic.NewClient(
		elastic.SetURL(hosts...),
		elastic.SetBasicAuth(esUser, esSecret),
		elastic.SetSniff(false),
		elastic.SetMaxRetries(3),
	)
	if err != nil {
		logs.Error("err = %v", err)
	}
	return err
}

//func SendES(esIndex string, msg string) {
//	newESClient(esIndex).sendMsg(msg)
//}
//
func newESClient(esIndex string) *Elastic {
	esType := esIndex
	esIndex = fmt.Sprintf("%s-%s", esIndex, time.Now().Format("200601"))

	ESLogHandler := &Elastic{
		esClient: esClient,
		esIndex:  esIndex,
		esType:   esType,
	}
	return ESLogHandler
}

func (this *Elastic) SendMsg(msg interface{}) {
	ctx := context.Background()

	logs.Info("Elastic SendMsg:%s", msg)
	doc, err := this.esClient.Index().
		Index(this.esIndex).
		Type(this.esType).
		BodyJson(msg).
		Do(ctx)
	if err != nil {
		logs.Info("add document error = %v", err)
	} else {
		logs.Info("index = %s, type = %s, id = %s", this.esIndex, this.esType, doc.Id)
	}
}
