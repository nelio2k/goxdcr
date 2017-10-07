// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package base

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/log"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MAX_PAYLOAD_SIZE  uint32        = 1000
	ReadWriteDeadline time.Duration = 1 * time.Second
)

type ConnType int

const (
	MemConn      ConnType = iota
	SSLOverMem   ConnType = iota
)

var (
	dialer *net.Dialer = &net.Dialer{Timeout: ShortHttpTimeout}
)

func (connType ConnType) String() string {
	if connType == MemConn {
		return "MemConn"
	} else if connType == SSLOverMem {
		return "SSLOverMem"
	} else {
		return "InvalidConnType"
	}
}

type NewConnFunc func() (*mcc.Client, error)

type ConnPool interface {
	Get() (*mcc.Client, error)
	GetNew() (*mcc.Client, error)
	GetCAS() uint32
	Release(client *mcc.Client)
	ReleaseConnections(cas uint32)
	NewConnFunc() NewConnFunc
	Name() string
	Size() int
	MaxConn() int
	ConnType() ConnType
	Hostname() string
	Password() string
	Close()
	Stale() bool
	SetStale(stale bool)
}

type SSLConnPool interface {
	ConnPool
	Certificate() []byte
}

type connPool struct {
	name     string
	clients  chan *mcc.Client
	hostName string
	// username and password used in setting up target connection
	userName    string
	password    string
	bucketName  string
	maxConn     int
	newConnFunc NewConnFunc
	logger      *log.CommonLogger
	lock        *sync.RWMutex
	cas         uint32
	plainAuth   bool
	stale       bool
	state_lock  *sync.RWMutex
}

type sslOverMemConnPool struct {
	connPool
	remote_memcached_port int
	certificate           []byte
	// whether target cluster supports SANs in certificates
	san_in_certificate bool
}

type connPoolMgr struct {
	conn_pools_map map[string]ConnPool
	map_lock       sync.RWMutex
	once           sync.Once
	logger         *log.CommonLogger
}

var _connPoolMgr connPoolMgr

// ensure that _connPoolMgr is initialized
func init() {
	ConnPoolMgr()
	mcc.DefaultDialTimeout = ShortHttpTimeout
}

var WrongConnTypeError = errors.New("There is an exiting pool with the same name but with different connection type")

/******************************************************************
 *
 *  Connection management
 *
 ******************************************************************/

func parseUsernamePassword(u string) (username string, password string, err error) {
	username = ""
	password = ""

	var url *url.URL
	url, err = url.Parse(u)
	if err != nil {
		return "", "", err
	}

	user := url.User
	if user != nil {
		username = user.Username()
		var isSet bool
		password, isSet = user.Password()
		if !isSet {
			password = ""
		}
	}

	return username, password, nil
}

func (p *connPool) init() {
	p.newConnFunc = p.newConn
}

func (p *connPool) Name() string {
	return p.name
}

func (p *connPool) Hostname() string {
	return p.hostName
}

func (p *connPool) Password() string {
	return p.password
}

func (p *connPool) Size() int {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.clients != nil {
		return len(p.clients)
	} else {
		return 0
	}
}

func (p *connPool) MaxConn() int {
	return p.maxConn
}

func (p *connPool) Get() (*mcc.Client, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.clients != nil {
		p.logger.Debugf("There are %d connections in the pool\n", len(p.clients))
		select {
		case client, ok := <-p.clients:
			if ok {
				return client, nil
			}
		default:
			//no more connection, create more
			mcClient, err := p.newConnFunc()
			return mcClient, err
		}
	}
	return nil, errors.New("connection pool is closed")
}

func (p *connPool) GetNew() (*mcc.Client, error) {
	return p.newConnFunc()
}

func (p *connPool) newConn() (*mcc.Client, error) {
	return NewConn(p.hostName, p.userName, p.password, p.bucketName, p.plainAuth, KeepAlivePeriod, p.logger)
}
func (p *connPool) NewConnFunc() NewConnFunc {
	return p.newConnFunc
}

func (p *connPool) ConnType() ConnType {
	return MemConn
}

func (p *connPool) Stale() bool {
	p.state_lock.RLock()
	defer p.state_lock.RUnlock()
	return p.stale
}

func (p *connPool) SetStale(stale bool) {
	p.state_lock.Lock()
	defer p.state_lock.Unlock()
	p.stale = stale
}

func authClient(client *mcc.Client, userName, password, bucketName string, plainAuth bool, logger *log.CommonLogger) error {
	var err error
	// authenticate using user/pass
	if userName != "" {
		if plainAuth {
			// use PLAIN authentication
			_, err = client.Auth(userName, password)
		} else {
			// use SCRAM-SHA authentication mechanisms
			_, err = client.AuthScramSha(userName, password)
		}

		if err != nil {
			logger.Errorf("err from authentication for user %v = %v\n", userName, err)
			client.Close()
			return err
		}
	}

	// if bucketName is different from userName, need to do selectBucket
	if bucketName != "" && bucketName != userName {
		_, err = client.SelectBucket(bucketName)
		if err != nil {
			logger.Errorf("err from select bucket for %v = %v\n", bucketName, err)
			client.Close()
			return err
		}
	}

	return nil
}

func (p *sslOverMemConnPool) init() {
	p.newConnFunc = p.newConn
}

func (p *sslOverMemConnPool) Certificate() []byte {
	return p.certificate
}

func (p *sslOverMemConnPool) newConn() (*mcc.Client, error) {
	ssl_con_str := p.hostName + UrlPortNumberDelimiter + strconv.FormatInt(int64(p.remote_memcached_port), ParseIntBase)
	return NewTLSConn(ssl_con_str, p.userName, p.password, p.certificate, p.san_in_certificate, p.bucketName, p.logger)
}

func (p *sslOverMemConnPool) ConnType() ConnType {
	return SSLOverMem
}

//
// Release connection back to the pool
//
func (p *connPool) Release(client *mcc.Client) {
	//reset connection deadlines
	conn := client.Hijack()

	conn.(net.Conn).SetReadDeadline(time.Date(1, time.January, 0, 0, 0, 0, 0, time.UTC))
	conn.(net.Conn).SetWriteDeadline(time.Date(1, time.January, 0, 0, 0, 0, 0, time.UTC))

	p.lock.RLock()
	defer p.lock.RUnlock()
	if p.clients != nil {
		select {
		case p.clients <- client:
			return
		default:
			//the pool reaches its capacity, drop the client on the floor
			client.Close()
			return
		}
	}
}

func (p *connPool) GetCAS() uint32 {
	p.state_lock.RLock()
	defer p.state_lock.RUnlock()
	return p.cas
}

func (p *connPool) incrementCAS() {
	p.state_lock.Lock()
	defer p.state_lock.Unlock()
	p.cas++
}

func (p *connPool) doesCASMatch(cas uint32) bool {
	p.state_lock.RLock()
	defer p.state_lock.RUnlock()
	if p.cas == cas {
		return true
	}
	return false
}

//
// Release all connections in the connection pool.
//
func (p *connPool) ReleaseConnections(cas uint32) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if !p.doesCASMatch(cas) {
		// no op if cas value does not match
		return
	}

	defer p.incrementCAS()

	if p.clients == nil {
		return
	}

	done := false
	for !done {
		select {
		case client, ok := <-p.clients:
			{
				if ok {
					if client != nil {
						client.Close()
					}
				} else {
					done = true
				}
			}
		default:
			{
				// if there is no more client in the channel
				done = true
			}
		}
	}
}

func (p *connPool) Close() {
	p.ReleaseConnections(p.GetCAS())

	p.lock.Lock()
	defer p.lock.Unlock()
	close(p.clients)
	p.clients = nil

}
func (connPoolMgr *connPoolMgr) GetOrCreatePool(poolNameToCreate string, hostname string, bucketname string, username string, password string, connsize int, plainauth bool) (ConnPool, error) {
	connPoolMgr.map_lock.Lock()
	defer connPoolMgr.map_lock.Unlock()

	pool, ok := connPoolMgr.conn_pools_map[poolNameToCreate]
	if ok {
		_, ok = pool.(*connPool)
		if ok {
			if !pool.Stale() && pool.Password() == password {
				return pool, nil
			} else {
				ConnPoolMgr().logger.Infof("Removing pool %v. stale=%v, new size=%v, old size=%v", poolNameToCreate, pool.Stale(), connsize, pool.MaxConn())
				connPoolMgr.removePool(pool)
			}
		} else {
			return nil, WrongConnTypeError
		}
	}

	var err error
	size := connsize
	if size == 0 {
		size = DefaultConnectionSize
	}
	pool = &connPool{clients: make(chan *mcc.Client, size),
		hostName:   hostname,
		userName:   username,
		password:   password,
		bucketName: bucketname,
		maxConn:    size,
		name:       poolNameToCreate,
		plainAuth:  plainauth,
		lock:       &sync.RWMutex{},
		state_lock: &sync.RWMutex{},
		logger:     log.NewLogger("ConnPool", connPoolMgr.logger.LoggerContext())}
	connPoolMgr.conn_pools_map[poolNameToCreate] = pool

	pool.(*connPool).init()
	return pool, err
}

func (connPoolMgr *connPoolMgr) GetOrCreateSSLOverMemPool(poolNameToCreate string, hostname string, bucketname string, username string, password string, connsize int, remote_mem_port int, cert []byte, san_in_cert bool) (ConnPool, error) {
	connPoolMgr.map_lock.Lock()
	defer connPoolMgr.map_lock.Unlock()

	pool, ok := connPoolMgr.conn_pools_map[poolNameToCreate]
	if ok {
		_, ok = pool.(*sslOverMemConnPool)
		if ok {
			if !pool.Stale() && pool.Password() == password {
				return pool, nil
			} else {
				ConnPoolMgr().logger.Infof("Removing pool %v. stale=%v, new size=%v, old size=%v", poolNameToCreate, pool.Stale(), connsize, pool.MaxConn())
				connPoolMgr.removePool(pool)
			}
		} else {
			connPoolMgr.logger.Errorf("Found existing pool with name=%v, connType=%v\n", pool.Name(), pool.ConnType())
			return nil, WrongConnTypeError
		}
	}

	var err error
	size := connsize
	if size == 0 {
		size = DefaultConnectionSize
	}
	p := &sslOverMemConnPool{
		connPool: connPool{clients: make(chan *mcc.Client, size),
			hostName:   hostname,
			userName:   username,
			password:   password,
			bucketName: bucketname,
			maxConn:    size,
			name:       poolNameToCreate,
			plainAuth:  true,
			lock:       &sync.RWMutex{},
			state_lock: &sync.RWMutex{},
			logger:     log.NewLogger("sslConnPool", connPoolMgr.logger.LoggerContext())},
		remote_memcached_port: remote_mem_port,
		certificate:           cert,
		san_in_certificate:    san_in_cert}
	p.init()

	connPoolMgr.conn_pools_map[poolNameToCreate] = p
	return p, err

}

func (connPoolMgr *connPoolMgr) GetPool(poolName string) ConnPool {
	connPoolMgr.map_lock.RLock()
	defer connPoolMgr.map_lock.RUnlock()
	pool := connPoolMgr.conn_pools_map[poolName]

	if pool != nil {
		connPoolMgr.logger.Infof("Successfully retrieved connection pool with name %v", poolName)
	} else {
		connPoolMgr.logger.Errorf("Could not find connection pool with name %v", poolName)
	}

	return pool
}

func (connPoolMgr *connPoolMgr) FindPoolNamesByPrefix(poolNamePrefix string) []string {
	poolNames := []string{}
	connPoolMgr.map_lock.RLock()
	defer connPoolMgr.map_lock.RUnlock()
	for poolName, _ := range connPoolMgr.conn_pools_map {
		if strings.HasPrefix(poolName, poolNamePrefix) {
			poolNames = append(poolNames, poolName)
		}
	}

	return poolNames
}

func (connPoolMgr *connPoolMgr) SetStaleForPoolsWithNamePrefix(poolNamePrefix string) {
	connPoolMgr.map_lock.RLock()
	defer connPoolMgr.map_lock.RUnlock()
	for poolName, pool := range connPoolMgr.conn_pools_map {
		if strings.HasPrefix(poolName, poolNamePrefix) {
			pool.SetStale(true)
			connPoolMgr.logger.Infof("Set pool %v as stale.", pool.Name())
		}
	}
}

func (connPoolMgr *connPoolMgr) RemovePool(poolName string) {
	connPoolMgr.map_lock.Lock()
	defer connPoolMgr.map_lock.Unlock()
	connPoolMgr.removePool(connPoolMgr.conn_pools_map[poolName])
}

func (connPoolMgr *connPoolMgr) removePool(pool ConnPool) {
	if pool != nil {
		pool.Close()
		delete(connPoolMgr.conn_pools_map, pool.Name())
		connPoolMgr.logger.Infof("Pool %v is removed, all connections are released", pool.Name())
	}
}

func (connPoolMgr *connPoolMgr) fillPool(p ConnPool, connectionSize int) error {
	connPoolMgr.logger.Infof("Fill Pool - poolName=%v,connType=%v, connectionSize=%d\n", p.Name(), p.ConnType, connectionSize)

	//	 initialize the connection pool
	work_load := 10
	num_of_workers := int(math.Ceil(float64(connectionSize) / float64(work_load)))
	index := 0
	waitGrp := &sync.WaitGroup{}
	for i := 0; i < num_of_workers; i++ {
		var connectionsToCreate int
		if index+work_load < connectionSize {
			connectionsToCreate = work_load
		} else {
			connectionsToCreate = connectionSize - index
		}
		f := p.NewConnFunc()
		if f == nil {
			return fmt.Errorf("Pool %v is not properly initialized, no NewConnFunc is set")
		}
		waitGrp.Add(1)
		go func(connectionsToCreate int, waitGrp *sync.WaitGroup, f NewConnFunc) {
			defer waitGrp.Done()
			for i := 0; i < connectionsToCreate; i++ {
				mcClient, err := f()
				if err == nil {
					p.Release(mcClient)
					connPoolMgr.logger.Info("A client connection is established")
				} else {
					connPoolMgr.logger.Errorf("error establishing new connection for pool %v, connectionsToCreate=%v, err=%v", p.Name(), connectionsToCreate, err)
				}
			}
		}(connectionsToCreate, waitGrp, f)

		index = index + connectionsToCreate
	}

	waitGrp.Wait()

	if p.Size() == 0 {
		return fmt.Errorf("Failed to fill connection pool of size %v for %v\n", connectionSize, p.Name())
	}

	connPoolMgr.logger.Infof("Connection pool %s is created with %d clients\n", p.Name(), p.Size())
	return nil

}

func encodeSSLHandShakeMsg(bytes []byte) []byte {
	ret := make([]byte, 4+len(bytes))
	binary.BigEndian.PutUint32(ret[0:4], uint32(len(bytes)))
	copy(ret[4:4+len(bytes)], bytes)
	return ret
}

//return the singleton ConnPoolMgr
func ConnPoolMgr() *connPoolMgr {
	_connPoolMgr.once.Do(func() {
		_connPoolMgr.conn_pools_map = make(map[string]ConnPool)
		_connPoolMgr.logger = log.NewLogger("ConnPoolMgr", log.DefaultLoggerContext)

	})
	return &_connPoolMgr
}

func SetLoggerContexForConnPoolMgr(logger_context *log.LoggerContext) *connPoolMgr {
	connPoolMgr := ConnPoolMgr()
	connPoolMgr.logger = log.NewLogger("ConnPoolMgr", logger_context)
	return connPoolMgr
}

func (connPoolMgr *connPoolMgr) Close() {
	connPoolMgr.map_lock.Lock()
	defer connPoolMgr.map_lock.Unlock()

	for key, pool := range connPoolMgr.conn_pools_map {
		connPoolMgr.logger.Infof("close pool %s", key)
		pool.ReleaseConnections(pool.GetCAS())
	}

	connPoolMgr.conn_pools_map = make(map[string]ConnPool)
}

// plainAuth is set to false only when
// 1. we are connecting to target memcached
// 2. the remote cluster reference is of half-ssl enabled type
func NewConn(hostName string, userName string, password string, bucketName string, plainAuth bool, keepAlivePeriod time.Duration, logger *log.CommonLogger) (conn *mcc.Client, err error) {
	// connect to host
	start_time := time.Now()
	conn, err = mcc.Connect("tcp", hostName)
	if err != nil {
		return nil, err
	}

	logger.Debugf("%vs spent on establish a connection to %v", time.Since(start_time).Seconds(), hostName)

	err = authClient(conn, userName, password, bucketName, plainAuth, logger)
	if err != nil {
		return nil, err
	}

	if keepAlivePeriod > 0 {
		conn.SetKeepAliveOptions(keepAlivePeriod)
	}

	logger.Debugf("%vs spent on authenticate to %v", time.Since(start_time).Seconds(), hostName)
	return conn, nil
}

func NewTLSConn(ssl_con_str string, username string, password string, certificate []byte, san_in_certificate bool, bucketName string, logger *log.CommonLogger) (*mcc.Client, error) {
	if len(certificate) == 0 {
		return nil, fmt.Errorf("No certificate has been provided. Can't establish ssl connection to %v", ssl_con_str)
	}

	logger.Infof("Trying to create a ssl over memcached connection on %v", ssl_con_str)
	conn, _, err := MakeTLSConn(ssl_con_str, certificate, san_in_certificate, logger)
	if err != nil {
		return nil, err
	}

	client, err := mcc.Wrap(conn)
	if err != nil {
		logger.Errorf("Failed to wrap connection. err=%v\n", err)
		conn.Close()
		return nil, err
	}

	err = authClient(client, username, password, bucketName, true /*plainAuth*/, logger)
	if err != nil {
		return nil, err
	}

	logger.Infof("memcached client on ssl connection %v has been created successfully", ssl_con_str)
	return client, nil
}

func MakeTLSConn(ssl_con_str string, certificate []byte, check_server_name bool, logger *log.CommonLogger) (*tls.Conn, *tls.Config, error) {
	// enforce timeout
	errChannel := make(chan error, 2)
	time.AfterFunc(dialer.Timeout, func() {
		errChannel <- ExecutionTimeoutError
	})

	caPool := x509.NewCertPool()
	ok := caPool.AppendCertsFromPEM(certificate)
	if !ok {
		return nil, nil, InvalidCerfiticateError
	}

	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return nil, nil, InvalidCerfiticateError
	}
	cert_remote, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, InvalidCerfiticateError
	}

	tlsConfig := &tls.Config{RootCAs: caPool}
	tlsConfig.BuildNameToCertificate()
	tlsConfig.InsecureSkipVerify = true
	hostname := strings.Split(ssl_con_str, UrlPortNumberDelimiter)[0]
	tlsConfig.ServerName = hostname

	// golang 1.8 added a new curve, X25519, which is not supported by ns_server pre-spock
	// explicitly define curve preferences to get this new curve excluded
	tlsConfig.CurvePreferences = []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.CurveP521}

	// get tcp connection
	rawConn, err := dialer.Dial("tcp", ssl_con_str)

	if err != nil {
		logger.Errorf("Failed to connect to %v, err=%v\n", ssl_con_str, err)
		return nil, nil, err
	}

	tcpConn, ok := rawConn.(*net.TCPConn)
	if !ok {
		// should never get here
		rawConn.Close()
		logger.Errorf("Failed to get tcp connection when connecting to %v\n", ssl_con_str)
		return nil, nil, err
	}

	// always set keep alive
	err = tcpConn.SetKeepAlive(true)
	if err == nil {
		err = tcpConn.SetKeepAlivePeriod(KeepAlivePeriod)
	}
	if err != nil {
		tcpConn.Close()
		logger.Errorf("Failed to set keep alive options when connecting to %v. err=%v\n", ssl_con_str, err)
		return nil, nil, err
	}

	// wrap as tls connection
	tlsConn := tls.Client(tcpConn, tlsConfig)

	// spawn new routine to enforce timeout
	go func() {
		errChannel <- tlsConn.Handshake()
	}()

	err = <-errChannel

	if err != nil {
		tlsConn.Close()
		logger.Errorf("TLS handshake failed when connecting to %v, err=%v\n", ssl_con_str, err)
		return nil, nil, err
	}

	if cert_remote.IsCA {
		connState := tlsConn.ConnectionState()
		peer_certs := connState.PeerCertificates

		opts := x509.VerifyOptions{
			Roots:         tlsConfig.RootCAs,
			CurrentTime:   time.Now(),
			Intermediates: x509.NewCertPool(),
		}

		if check_server_name {
			// need to check server name. get sever name from ssl_con_str
			opts.DNSName = hostname
		} else {
			logger.Debug("remote peer is old and its certificate doesn't have IP SANs, skip verifying ServerName")
		}

		for i, cert := range peer_certs {
			if i == 0 {
				continue
			}
			opts.Intermediates.AddCert(cert)
		}
		_, err = peer_certs[0].Verify(opts)
		if err != nil {
			//close the conn
			tlsConn.Close()
			logger.Errorf("TLS Verify failed when connecting to %v, err=%v\n", ssl_con_str, err)
			return nil, nil, err
		}
	}
	return tlsConn, tlsConfig, nil

}

func DialTCPWithTimeout(network, address string) (net.Conn, error) {
	return dialer.Dial(network, address)
}
