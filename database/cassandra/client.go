package cassandra

import (
	"crypto/tls"

	"github.com/gocql/gocql"

	"github.com/go-kratos/kratos/v2/log"

	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func NewCassandraClient(cfg *conf.Bootstrap, l *log.Helper) *gocql.Session {
	if cfg.Data == nil || cfg.Data.Cassandra == nil {
		l.Warn("cassandra config is nil")
		return nil
	}

	clusterConfig := gocql.NewCluster(cfg.Data.Cassandra.GetAddress())

	// 设置用户名密码
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Data.Cassandra.GetUsername(),
		Password: cfg.Data.Cassandra.GetPassword(),
	}

	clusterConfig.Keyspace = cfg.Data.Cassandra.GetKeyspace()

	// 设置ssl
	if cfg.Data.Cassandra.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = loadServerTlsConfig(cfg.Data.Cassandra.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			clusterConfig.SslOpts = &gocql.SslOptions{Config: tlsCfg}
		}
	}

	// 设置超时时间
	clusterConfig.ConnectTimeout = cfg.Data.Cassandra.ConnectTimeout.AsDuration()
	clusterConfig.Timeout = cfg.Data.Cassandra.Timeout.AsDuration()

	clusterConfig.Consistency = gocql.Consistency(cfg.Data.Cassandra.GetConsistency())

	// 禁止主机查找
	clusterConfig.DisableInitialHostLookup = cfg.Data.Cassandra.GetDisableInitialHostLookup()

	session, err := clusterConfig.CreateSession()
	if err != nil {
		l.Fatalf("failed opening connection to cassandra: %v", err)
		return nil
	}

	return session
}

func loadServerTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadServerTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
			cfg.InsecureSkipVerify,
		); err != nil {
			return nil, err
		}
	}

	return tlsCfg, err
}
