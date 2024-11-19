package cassandra

import (
	"crypto/tls"

	"github.com/gocql/gocql"

	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

func NewCassandraClient(cfg *conf.Bootstrap, l *log.Helper) *gocql.Session {
	if cfg.Data == nil || cfg.Data.Cassandra == nil {
		l.Warn("cassandra config is nil")
		return nil
	}

	clusterConfig := gocql.NewCluster(cfg.Data.Cassandra.Address)

	// 设置用户名密码
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Data.Cassandra.Username,
		Password: cfg.Data.Cassandra.Password,
	}

	clusterConfig.Keyspace = cfg.Data.Cassandra.Keyspace

	// 设置ssl
	if cfg.Data.Cassandra.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadServerTlsConfig(cfg.Data.Cassandra.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			clusterConfig.SslOpts = &gocql.SslOptions{Config: tlsCfg}
		}
	}

	// 设置超时时间
	clusterConfig.ConnectTimeout = cfg.Data.Cassandra.ConnectTimeout.AsDuration()
	clusterConfig.Timeout = cfg.Data.Cassandra.Timeout.AsDuration()

	clusterConfig.Consistency = gocql.Consistency(cfg.Data.Cassandra.Consistency)

	// 禁止主机查找
	clusterConfig.DisableInitialHostLookup = cfg.Data.Cassandra.DisableInitialHostLookup

	session, err := clusterConfig.CreateSession()
	if err != nil {
		l.Fatalf("failed opening connection to cassandra: %v", err)
		return nil
	}

	return session
}
