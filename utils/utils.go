package utils

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"github.com/couchbase/cbauth"
	"github.com/couchbase/go-couchbase"
	mc "github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	base "github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

var NonExistentBucketError error = errors.New("Bucket doesn't exist")

type Utilities struct {
	logger_utils *log.CommonLogger
}

/**
 * NOTE - ideally we want to be able to pass in utility interfaces so we can do much
 * better unit testing with mocks. This constructor should be used in main() and then
 * passed down level by levels.
 * Currently, this method is being called in many places, and each place that is using
 * this method should ideally be using a passed in interface from another parent level.
 */
func NewUtilities() *Utilities {
	retVar := &Utilities{
		logger_utils: log.NewLogger("Utils", log.DefaultLoggerContext),
	}
	return retVar
}

/**
 * This utils file contains both regular non-REST related utilities as well as REST related.
 * The first section is non-REST related utility functions
 */

type BucketBasicStats struct {
	ItemCount int `json:"itemCount"`
}

//Only used by unit test
//TODO: replace with go-couchbase bucket stats API
type CouchBucket struct {
	Name string           `json:"name"`
	Stat BucketBasicStats `json:"basicStats"`
}

func (u *Utilities) GetNonExistentBucketError() error {
	return NonExistentBucketError
}

//func (u *Utilities) GetLoggerUtils (*log.CommonLogger) {
//	return u.logger_utils
//}

func (u *Utilities) loggerForFunc(logger *log.CommonLogger) *log.CommonLogger {
	var l *log.CommonLogger
	if logger != nil {
		l = logger
	} else {
		l = u.logger_utils
	}
	return l
}

func (u *Utilities) ValidateSettings(defs base.SettingDefinitions,
	settings map[string]interface{},
	logger *log.CommonLogger) error {
	var l *log.CommonLogger = u.loggerForFunc(logger)

	l.Debugf("Start validate setting=%v, defs=%v", settings, defs)
	var err *base.SettingsError = nil
	for key, def := range defs {
		val, ok := settings[key]
		if !ok && def.Required {
			if err == nil {
				err = base.NewSettingsError()
			}
			err.Add(key, errors.New("required, but not supplied"))
		} else {
			if val != nil && def.Data_type != reflect.PtrTo(reflect.TypeOf(val)) {
				if err == nil {
					err = base.NewSettingsError()
				}
				err.Add(key, errors.New(fmt.Sprintf("expected type is %v, supplied type is %v",
					def.Data_type, reflect.TypeOf(val))))
			}
		}
	}
	if err != nil {
		l.Infof("setting validation result = %v", *err)
		return *err
	}
	return nil
}

func (u *Utilities) RecoverPanic(err *error) {
	if r := recover(); r != nil {
		*err = errors.New(fmt.Sprint(r))
	}
}

func (u *Utilities) LocalPool(localConnectStr string) (couchbase.Pool, error) {
	localURL := fmt.Sprintf("http://%s", localConnectStr)
	client, err := couchbase.ConnectWithAuth(localURL, cbauth.NewAuthHandler(nil))
	if err != nil {
		return couchbase.Pool{}, u.NewEnhancedError(fmt.Sprintf("Error connecting to couchbase. url=%v", u.UrlForLog(localURL)), err)
	}
	return client.GetPool("default")
}

// Get bucket in local cluster
func (u *Utilities) LocalBucket(localConnectStr, bucketName string) (*couchbase.Bucket, error) {
	u.logger_utils.Debugf("Getting local bucket name=%v\n", bucketName)

	pool, err := u.LocalPool(localConnectStr)
	if err != nil {
		return nil, err
	}

	bucket, err := pool.GetBucket(bucketName)
	if err != nil {
		return nil, u.NewEnhancedError(fmt.Sprintf("Error getting bucket, %v, from pool.", bucketName), err)
	}

	u.logger_utils.Debugf("Got local bucket successfully name=%v\n", bucket.Name)
	return bucket, err
}

func (u *Utilities) UnwrapError(infos map[string]interface{}) (err error) {
	if infos != nil && len(infos) > 0 {
		err = infos["error"].(error)
	}
	return err
}

// returns an enhanced error with erroe message being "msg + old error message"
func (u *Utilities) NewEnhancedError(msg string, err error) error {
	return errors.New(msg + "\n err = " + err.Error())
}

func (u *Utilities) GetMapFromExpvarMap(expvarMap *expvar.Map) map[string]interface{} {
	regMap := make(map[string]interface{})

	expvarMap.Do(func(keyValue expvar.KeyValue) {
		valueStr := keyValue.Value.String()
		// first check if valueStr is an integer
		valueInt, err := strconv.Atoi(valueStr)
		if err == nil {
			regMap[keyValue.Key] = valueInt
		} else {
			// then check if valueStr is a float
			valueFloat, err := strconv.ParseFloat(valueStr, 64)
			if err == nil {
				regMap[keyValue.Key] = valueFloat
			} else {
				// should never happen
				u.logger_utils.Errorf("Invalid value in expvarMap. Only float and integer values are supported")
			}
		}
	})
	return regMap
}

//convert the format returned by go-memcached StatMap - map[string]string to map[uint16]uint64
func (u *Utilities) ParseHighSeqnoStat(vbnos []uint16, stats_map map[string]string, highseqno_map map[uint16]uint64) error {
	var err error
	for _, vbno := range vbnos {
		stats_key := fmt.Sprintf(base.VBUCKET_HIGH_SEQNO_STAT_KEY_FORMAT, vbno)
		highseqnostr, ok := stats_map[stats_key]
		if !ok || highseqnostr == "" {
			err = fmt.Errorf("Can't find high seqno for vbno=%v in stats map. Source topology may have changed.\n", vbno)
			return err
		}
		highseqno, err := strconv.ParseUint(highseqnostr, 10, 64)
		if err != nil {
			u.logger_utils.Warnf("high seqno for vbno=%v in stats map is not a valid uint64. high seqno=%v\n", vbno, highseqnostr)
			err = fmt.Errorf("high seqno for vbno=%v in stats map is not a valid uint64. high seqno=%v\n", vbno, highseqnostr)
			return err
		}
		highseqno_map[vbno] = highseqno
	}
	return nil
}

//convert the format returned by go-memcached StatMap - map[string]string to map[uint16][]uint64
func (u *Utilities) ParseHighSeqnoAndVBUuidFromStats(vbnos []uint16, stats_map map[string]string, high_seqno_and_vbuuid_map map[uint16][]uint64) {
	for _, vbno := range vbnos {
		high_seqno_stats_key := fmt.Sprintf(base.VBUCKET_HIGH_SEQNO_STAT_KEY_FORMAT, vbno)
		highseqnostr, ok := stats_map[high_seqno_stats_key]
		if !ok {
			u.logger_utils.Warnf("Can't find high seqno for vbno=%v in stats map. Source topology may have changed.\n", vbno)
			continue
		}
		high_seqno, err := strconv.ParseUint(highseqnostr, 10, 64)
		if err != nil {
			u.logger_utils.Warnf("high seqno for vbno=%v in stats map is not a valid uint64. high seqno=%v\n", vbno, highseqnostr)
			continue
		}

		vbuuid_stats_key := fmt.Sprintf(base.VBUCKET_UUID_STAT_KEY_FORMAT, vbno)
		vbuuidstr, ok := stats_map[vbuuid_stats_key]
		if !ok {
			u.logger_utils.Warnf("Can't find vbuuid for vbno=%v in stats map. Source topology may have changed.\n", vbno)
			continue
		}
		vbuuid, err := strconv.ParseUint(vbuuidstr, 10, 64)
		if err != nil {
			u.logger_utils.Warnf("vbuuid for vbno=%v in stats map is not a valid uint64. vbuuid=%v\n", vbno, vbuuidstr)
			continue
		}

		high_seqno_and_vbuuid_map[vbno] = []uint64{high_seqno, vbuuid}
	}
}

// encode data in a map into a byte array, which can then be used as
// the body part of a http request
// so far only five types are supported: string, int, bool, LogLevel, []byte
// which should be sufficient for all cases at hand
func (u *Utilities) EncodeMapIntoByteArray(data map[string]interface{}) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	params := make(url.Values)
	for key, val := range data {
		var strVal string
		switch val.(type) {
		case string:
			strVal = val.(string)
		case int:
			strVal = strconv.FormatInt(int64(val.(int)), base.ParseIntBase)
		case bool:
			strVal = strconv.FormatBool(val.(bool))
		case log.LogLevel:
			strVal = val.(log.LogLevel).String()
		case []byte:
			strVal = string(val.([]byte))
		default:
			return nil, base.IncorrectValueTypeInMapError(key, val, "string/int/bool/LogLevel/[]byte")
		}
		params.Add(key, strVal)
	}

	return []byte(params.Encode()), nil
}

func (u *Utilities) UrlForLog(urlStr string) string {
	result, err := url.Parse(urlStr)
	if err == nil {
		if result.User != nil {
			result.User = url.UserPassword(result.User.Username(), "xxxx")
		}
		return result.String()
	} else {
		return urlStr
	}
}

func (u *Utilities) GetMatchedKeys(expression string, keys []string) (map[string][][]int, error) {
	u.logger_utils.Infof("GetMatchedKeys expression=%v, expression in bytes=%v\n", expression, []byte(expression))
	if !utf8.ValidString(expression) {
		return nil, errors.New("expression is not valid utf8")
	}
	for _, key := range keys {
		u.logger_utils.Infof("key=%v, key_bytes=%v\n", key, []byte(key))
		if !utf8.ValidString(key) {
			return nil, errors.New("key is not valid utf8")
		}
	}

	regExp, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	matchesMap := make(map[string][][]int)

	for _, key := range keys {
		var matches [][]int
		if u.RegexpMatch(regExp, []byte(key)) {
			matches = regExp.FindAllStringIndex(key, -1)
		} else {
			matches = make([][]int, 0)
		}
		u.logger_utils.Debugf("key=%v, matches with byte index=%v\n", key, matches)
		convertedMatches, err := u.convertByteIndexToRuneIndex(key, matches)
		if err != nil {
			return nil, err
		}
		matchesMap[key] = convertedMatches
	}

	return matchesMap, nil
}

func (u *Utilities) RegexpMatch(regExp *regexp.Regexp, key []byte) bool {
	return regExp.Match(key)
}

// given a matches map, convert the indices from byte index to rune index
func (u *Utilities) convertByteIndexToRuneIndex(key string, matches [][]int) ([][]int, error) {
	convertedMatches := make([][]int, 0)
	if len(key) == 0 || len(matches) == 0 {
		return matches, nil
	}

	// parse key and build a byte index to rune index map
	indexMap := make(map[int]int)
	byteIndex := 0
	runeIndex := 0
	keyBytes := []byte(key)
	keyLen := len(key)
	for {
		indexMap[byteIndex] = runeIndex
		if byteIndex < keyLen {
			_, runeLen := utf8.DecodeRune(keyBytes[byteIndex:])
			byteIndex += runeLen
			runeIndex++
		} else {
			break
		}
	}

	u.logger_utils.Debugf("key=%v, indexMap=%v\n", key, indexMap)

	var ok bool
	for _, match := range matches {
		convertedMatch := make([]int, 2)
		convertedMatch[0], ok = indexMap[match[0]]
		if !ok {
			// should not happen
			errMsg := u.InvalidRuneIndexErrorMessage(key, match[0])
			u.logger_utils.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		convertedMatch[1], ok = indexMap[match[1]]
		if !ok {
			// should not happen
			errMsg := u.InvalidRuneIndexErrorMessage(key, match[1])
			u.logger_utils.Errorf(errMsg)
			return nil, errors.New(errMsg)
		}
		convertedMatches = append(convertedMatches, convertedMatch)
	}

	return convertedMatches, nil
}

func (u *Utilities) InvalidRuneIndexErrorMessage(key string, index int) string {
	return fmt.Sprintf("byte index, %v, in match for key, %v, is not a starting index for a rune", index, key)
}

func (u *Utilities) LocalBucketUUID(local_connStr string, bucketName string, logger *log.CommonLogger) (string, error) {
	return u.BucketUUID(local_connStr, bucketName, "", "", nil, false, logger)
}

func (u *Utilities) LocalBucketPassword(local_connStr string, bucketName string, logger *log.CommonLogger) (string, error) {
	return u.BucketPassword(local_connStr, bucketName, "", "", nil, false, logger)
}

func (u *Utilities) ReplicationStatusNotFoundError(topic string) error {
	return fmt.Errorf("Cannot find replication status for topic %v", topic)
}

func (u *Utilities) BucketNotFoundError(bucketName string) error {
	return fmt.Errorf("Bucket `%v` not found.", bucketName)
}

// creates a local memcached connection.
// always use plain auth
func (u *Utilities) GetMemcachedConnection(serverAddr, bucketName, userAgent string,
	keepAlivePeriod time.Duration, logger *log.CommonLogger) (mcc.ClientIface, error) {
	logger.Infof("GetMemcachedConnection serverAddr=%v, bucketName=%v\n", serverAddr, bucketName)
	if serverAddr == "" {
		err := fmt.Errorf("Failed to get memcached connection because serverAddr is empty. bucketName=%v, userAgent=%v", bucketName, userAgent)
		logger.Warnf(err.Error())
		return nil, err
	}
	username, password, err := cbauth.GetMemcachedServiceAuth(serverAddr)
	logger.Debugf("memcached auth: username=%v, password=%v, err=%v\n", username, password, err)
	if err != nil {
		return nil, err
	}

	return u.GetRemoteMemcachedConnection(serverAddr, username, password, bucketName, userAgent,
		true /*plainAuth*/, keepAlivePeriod, logger)
}

func (u *Utilities) GetRemoteMemcachedConnection(serverAddr, username, password, bucketName, userAgent string,
	plainAuth bool, keepAlivePeriod time.Duration, logger *log.CommonLogger) (mcc.ClientIface, error) {
	conn, err := base.NewConn(serverAddr, username, password, bucketName, plainAuth, keepAlivePeriod, logger)
	if err != nil {
		return nil, err
	}

	err = u.SendHELO(conn, userAgent, base.HELOTimeout, base.HELOTimeout, logger)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// send helo with specified user agent string to memcached
// the helo is purely informational, for the identification of the client
// unsuccessful response is not treated as errors
func (u *Utilities) SendHELO(client mcc.ClientIface, userAgent string, readTimeout, writeTimeout time.Duration,
	logger *log.CommonLogger) (err error) {
	heloReq := u.ComposeHELORequest(userAgent, false /*enableDataType*/)

	var response *mc.MCResponse
	response, err = u.sendHELORequest(client, heloReq, userAgent, readTimeout, writeTimeout, logger)
	if err != nil {
		logger.Errorf("Received error response from HELO command. userAgent=%v, err=%v.", userAgent, err)
	} else if response.Status != mc.SUCCESS {
		logger.Warnf("Received unexpected response from HELO command. userAgent=%v, response status=%v.", userAgent, response.Status)
	} else {
		logger.Infof("Successfully sent HELO command with userAgent=%v", userAgent)
	}
	return
}

// send helo to memcached with data type (including xattr) feature enabled
// used exclusively by xmem nozzle
// we need to know whether data type is indeed enabled from helo response
// unsuccessful response is treated as errors
func (u *Utilities) SendHELOWithXattrFeature(client mcc.ClientIface, userAgent string, readTimeout, writeTimeout time.Duration,
	logger *log.CommonLogger) (xattrEnabled bool, err error) {
	heloReq := u.ComposeHELORequest(userAgent, true /*enableXattr*/)

	var response *mc.MCResponse
	response, err = u.sendHELORequest(client, heloReq, userAgent, readTimeout, writeTimeout, logger)
	if err != nil {
		logger.Errorf("Received error response from HELO command. userAgent=%v, err=%v.", userAgent, err)
	} else if response.Status != mc.SUCCESS {
		errMsg := fmt.Sprintf("Received unexpected response from HELO command. userAgent=%v, response status=%v.", userAgent, response.Status)
		logger.Error(errMsg)
		err = errors.New(errMsg)
	} else {
		// helo succeeded. parse response body for features enabled
		bodyLen := len(response.Body)
		if (bodyLen & 1) != 0 {
			// body has to have even number of bytes
			errMsg := fmt.Sprintf("Received response body with odd number of bytes from HELO command. userAgent=%v, response body=%v.", userAgent, response.Body)
			logger.Error(errMsg)
			err = errors.New(errMsg)
			return
		}
		pos := 0
		for {
			if pos >= bodyLen {
				break
			}
			feature := binary.BigEndian.Uint16(response.Body[pos : pos+2])
			if feature == base.HELO_FEATURE_XATTR {
				xattrEnabled = true
				break
			}
			pos += 2
		}
		logger.Infof("Successfully sent HELO command with userAgent=%v. xattrEnabled=%v", userAgent, xattrEnabled)
	}
	return
}

func (u *Utilities) sendHELORequest(client mcc.ClientIface, heloReq *mc.MCRequest, userAgent string, readTimeout, writeTimeout time.Duration,
	logger *log.CommonLogger) (response *mc.MCResponse, err error) {

	conn := client.Hijack()
	conn.(net.Conn).SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err = conn.Write(heloReq.Bytes())
	conn.(net.Conn).SetWriteDeadline(time.Time{})
	if err != nil {
		logger.Warnf("Error sending HELO command. userAgent=%v, err=%v.", userAgent, err)
		return
	}

	conn.(net.Conn).SetReadDeadline(time.Now().Add(readTimeout))
	response, err = client.Receive()
	conn.(net.Conn).SetReadDeadline(time.Time{})
	return
}

// compose a HELO command
func (u *Utilities) ComposeHELORequest(userAgent string, enableXattr bool) *mc.MCRequest {
	var value []byte
	if enableXattr {
		value = make([]byte, 4)
		// tcp nodelay
		binary.BigEndian.PutUint16(value[0:2], base.HELO_FEATURE_TCP_NO_DELAY)
		// Xattr
		binary.BigEndian.PutUint16(value[2:4], base.HELO_FEATURE_XATTR)
	} else {
		value = make([]byte, 2)
		// tcp nodelay
		binary.BigEndian.PutUint16(value[0:2], base.HELO_FEATURE_TCP_NO_DELAY)
	}
	return &mc.MCRequest{
		Key:    []byte(userAgent),
		Opcode: mc.HELLO,
		Body:   value,
	}
}

func (u *Utilities) GetIntSettingFromSettings(settings map[string]interface{}, settingName string) (int, error) {
	settingObj := u.GetSettingFromSettings(settings, settingName)
	if settingObj == nil {
		return -1, nil
	}

	setting, ok := settingObj.(int)
	if !ok {
		return -1, fmt.Errorf("Setting %v is of wrong type", settingName)
	}

	return setting, nil
}

func (u *Utilities) GetStringSettingFromSettings(settings map[string]interface{}, settingName string) (string, error) {
	settingObj := u.GetSettingFromSettings(settings, settingName)
	if settingObj == nil {
		return "", nil
	}

	setting, ok := settingObj.(string)
	if !ok {
		return "", fmt.Errorf("Setting %v is of wrong type", settingName)
	}

	return setting, nil
}

func (u *Utilities) GetSettingFromSettings(settings map[string]interface{}, settingName string) interface{} {
	if settings == nil {
		return nil
	}

	setting, ok := settings[settingName]
	if !ok {
		return nil
	}

	return setting
}

func (u *Utilities) GetMemcachedClient(serverAddr, bucketName string, kv_mem_clients map[string]mcc.ClientIface,
	userAgent string, keepAlivePeriod time.Duration, logger *log.CommonLogger) (mcc.ClientIface, error) {
	client, ok := kv_mem_clients[serverAddr]
	if ok {
		return client, nil
	} else {
		if bucketName == "" {
			err := fmt.Errorf("Failed to get memcached client because of unexpected empty bucketName. serverAddr=%v, userAgent=%v", serverAddr, userAgent)
			logger.Warnf(err.Error())
			return nil, err
		}

		var client, err = u.GetMemcachedConnection(serverAddr, bucketName, userAgent, keepAlivePeriod, logger)
		if err == nil {
			kv_mem_clients[serverAddr] = client
			return client, nil
		} else {
			return nil, err
		}
	}
}

func (u *Utilities) GetServerVBucketsMap(connStr, bucketName string, bucketInfo map[string]interface{}) (map[string][]uint16, error) {
	vbucketServerMapObj, ok := bucketInfo[base.VBucketServerMapKey]
	if !ok {
		return nil, fmt.Errorf("Error getting vbucket server map from bucket info. connStr=%v, bucketName=%v, bucketInfo=%v\n", connStr, bucketName, bucketInfo)
	}
	vbucketServerMap, ok := vbucketServerMapObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Vbucket server map is of wrong type. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMapObj)
	}

	// get server list
	serverListObj, ok := vbucketServerMap[base.ServerListKey]
	if !ok {
		return nil, fmt.Errorf("Error getting server list from vbucket server map. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMap)
	}
	serverList, ok := serverListObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Server list is of wrong type. connStr=%v, bucketName=%v, serverList=%v\n", connStr, bucketName, serverListObj)
	}

	servers := make([]string, len(serverList))
	for index, serverName := range serverList {
		serverNameStr, ok := serverName.(string)
		if !ok {
			return nil, fmt.Errorf("Server name is of wrong type. connStr=%v, bucketName=%v, serverName=%v\n", connStr, bucketName, serverName)
		}
		servers[index] = serverNameStr
	}

	// get vbucket "map"
	vbucketMapObj, ok := vbucketServerMap[base.VBucketMapKey]
	if !ok {
		return nil, fmt.Errorf("Error getting vbucket map from vbucket server map. connStr=%v, bucketName=%v, vbucketServerMap=%v\n", connStr, bucketName, vbucketServerMap)
	}
	vbucketMap, ok := vbucketMapObj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Vbucket map is of wrong type. connStr=%v, bucketName=%v, vbucketMap=%v\n", connStr, bucketName, vbucketMapObj)
	}

	serverVBMap := make(map[string][]uint16)

	for vbno, indexListObj := range vbucketMap {
		indexList, ok := indexListObj.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Index list is of wrong type. connStr=%v, bucketName=%v, indexList=%v\n", connStr, bucketName, indexListObj)
		}
		if len(indexList) == 0 {
			return nil, fmt.Errorf("Index list is empty. connStr=%v, bucketName=%v, vbno=%v\n", connStr, bucketName, vbno)
		}
		indexFloat, ok := indexList[0].(float64)
		if !ok {
			return nil, fmt.Errorf("Master index is of wrong type. connStr=%v, bucketName=%v, index=%v\n", connStr, bucketName, indexList[0])
		}
		indexInt := int(indexFloat)
		if indexInt >= len(servers) {
			return nil, fmt.Errorf("Master index is out of range. connStr=%v, bucketName=%v, index=%v\n", connStr, bucketName, indexInt)
		} else if indexInt < 0 {
			// During rebalancing or topology changes, it's possible ns_server may return a -1 for index. Callers should treat it as an transient error.
			return nil, fmt.Errorf(fmt.Sprintf("%v connStr=%v, bucketName=%v, index=%v\n", base.ErrorMasterNegativeIndex, connStr, bucketName, indexInt))
		}

		server := servers[indexInt]
		var vbList []uint16
		vbList, ok = serverVBMap[server]
		if !ok {
			vbList = make([]uint16, 0)
		}
		vbList = append(vbList, uint16(vbno))
		serverVBMap[server] = vbList
	}
	return serverVBMap, nil
}

// get bucket type setting from bucket info
func (u *Utilities) GetBucketTypeFromBucketInfo(bucketName string, bucketInfo map[string]interface{}) (string, error) {
	bucketType := ""
	bucketTypeObj, ok := bucketInfo[base.BucketTypeKey]
	if !ok {
		return "", fmt.Errorf("Error looking up bucket type of bucket %v", bucketName)
	} else {
		bucketType, ok = bucketTypeObj.(string)
		if !ok {
			return "", fmt.Errorf("bucketType on bucket %v is of wrong type.", bucketName)
		}
	}
	return bucketType, nil
}

// check if a bucket belongs to an elastic search (es) cluster by looking for "authType" field in bucket info.
// if not found, cluster is es
func (u *Utilities) CheckWhetherClusterIsESBasedOnBucketInfo(bucketInfo map[string]interface{}) bool {
	_, ok := bucketInfo[base.AuthTypeKey]
	return !ok
}

// get conflict resolution type setting from bucket info
func (u *Utilities) GetConflictResolutionTypeFromBucketInfo(bucketName string, bucketInfo map[string]interface{}) (string, error) {
	conflictResolutionType := base.ConflictResolutionType_Seqno
	conflictResolutionTypeObj, ok := bucketInfo[base.ConflictResolutionTypeKey]
	if ok {
		conflictResolutionType, ok = conflictResolutionTypeObj.(string)
		if !ok {
			return "", fmt.Errorf("ConflictResolutionType on bucket %v is of wrong type.", bucketName)
		}
	}
	return conflictResolutionType, nil
}

// get EvictionPolicy setting from bucket info
func (u *Utilities) GetEvictionPolicyFromBucketInfo(bucketName string, bucketInfo map[string]interface{}) (string, error) {
	evictionPolicy := ""
	evictionPolicyObj, ok := bucketInfo[base.EvictionPolicyKey]
	if ok {
		evictionPolicy, ok = evictionPolicyObj.(string)
		if !ok {
			return "", fmt.Errorf("EvictionPolicy on bucket %v is of wrong type.", bucketName)
		}
	}
	return evictionPolicy, nil
}

/**
 * The second section is couchbase REST related utility functions
 */
func (u *Utilities) GetMemcachedSSLPortMap(hostName, username, password string, certificate []byte, sanInCertificate bool, bucket string, logger *log.CommonLogger) (map[string]uint16, error) {
	ret := make(map[string]uint16)

	logger.Infof("GetMemcachedSSLPort, hostName=%v\n", hostName)
	bucketInfo, err := u.GetClusterInfo(hostName, base.BPath+bucket, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return nil, err
	}

	nodesExt, ok := bucketInfo[base.NodeExtKey]
	if !ok {
		return nil, u.BucketInfoParseError(bucketInfo, logger)
	}

	nodesExtArray, ok := nodesExt.([]interface{})
	if !ok {
		return nil, u.BucketInfoParseError(bucketInfo, logger)
	}

	for _, nodeExt := range nodesExtArray {

		nodeExtMap, ok := nodeExt.(map[string]interface{})
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		hostname, err := u.GetHostNameFromNodeInfo(hostName, nodeExtMap, logger)
		if err != nil {
			return nil, err
		}

		service, ok := nodeExtMap[base.ServicesKey]
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		services_map, ok := service.(map[string]interface{})
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		kv_port, ok := services_map[base.KVPortKey]
		if !ok {
			// the node may not have kv services. skip the node
			continue
		}
		kvPortFloat, ok := kv_port.(float64)
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		hostAddr := base.GetHostAddr(hostname, uint16(kvPortFloat))

		kv_ssl_port, ok := services_map[base.KVSSLPortKey]
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		kvSSLPortFloat, ok := kv_ssl_port.(float64)
		if !ok {
			return nil, u.BucketInfoParseError(bucketInfo, logger)
		}

		ret[hostAddr] = uint16(kvSSLPortFloat)
	}
	logger.Infof("memcached ssl port map=%v\n", ret)

	return ret, nil
}

func (u *Utilities) BucketInfoParseError(bucketInfo map[string]interface{}, logger *log.CommonLogger) error {
	errMsg := "Error parsing memcached ssl port of remote cluster."
	detailedErrMsg := errMsg + fmt.Sprintf("bucketInfo=%v", bucketInfo)
	logger.Errorf(detailedErrMsg)
	return fmt.Errorf(errMsg)
}

func (u *Utilities) HttpsHostAddr(hostAddr string, logger *log.CommonLogger) (string, error, bool) {
	hostName := base.GetHostName(hostAddr)
	sslPort, err, isInternalError := u.GetSSLPort(hostAddr, logger)
	if err != nil {
		return "", err, isInternalError
	}
	return base.GetHostAddr(hostName, sslPort), nil, false
}

func (u *Utilities) GetSSLPort(hostAddr string, logger *log.CommonLogger) (uint16, error, bool) {
	portInfo := make(map[string]interface{})
	err, statusCode := u.QueryRestApiWithAuth(hostAddr, base.SSLPortsPath, false, "", "", nil, false, base.MethodGet, "", nil, 0, &portInfo, nil, false, logger)
	if err != nil || statusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed on calling %v, err=%v, statusCode=%v", base.SSLPortsPath, err, statusCode), false
	}
	sslPort, ok := portInfo[base.SSLPortKey]
	if !ok {
		errMsg := "Failed to parse port info. ssl port is missing."
		logger.Errorf("%v. portInfo=%v", errMsg, portInfo)
		return 0, fmt.Errorf(errMsg), true
	}

	sslPortFloat, ok := sslPort.(float64)
	if !ok {
		return 0, fmt.Errorf("ssl port is of wrong type. Expected type: float64; Actual type: %s", reflect.TypeOf(sslPort)), true
	}

	return uint16(sslPortFloat), nil, false
}

func (u *Utilities) GetClusterInfoWStatusCode(hostAddr, path, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (map[string]interface{}, error, int) {
	clusterInfo := make(map[string]interface{})
	err, statusCode := u.QueryRestApiWithAuth(hostAddr, path, false, username, password, certificate, sanInCertificate, base.MethodGet, "", nil, 0, &clusterInfo, nil, false, logger)
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed on calling host=%v, path=%v, err=%v, statusCode=%v", hostAddr, path, err, statusCode), statusCode
	}
	return clusterInfo, nil, statusCode
}

func (u *Utilities) GetClusterInfo(hostAddr, path, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (map[string]interface{}, error) {
	clusterInfo, err, _ := u.GetClusterInfoWStatusCode(hostAddr, path, username, password, certificate, sanInCertificate, logger)
	return clusterInfo, err
}

func (u *Utilities) GetClusterUUID(hostAddr, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (string, error) {
	clusterInfo, err := u.GetClusterInfo(hostAddr, base.PoolsPath, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return "", err
	}
	clusterUUIDObj, ok := clusterInfo[base.UUIDKey]
	if !ok {
		return "", fmt.Errorf("Cannot find uuid key in cluster info. hostAddr=%v, clusterInfo=%v\n", hostAddr, clusterInfo)
	}
	clusterUUID, ok := clusterUUIDObj.(string)
	if !ok {
		// cluster uuid is "[]" for unintialized cluster
		_, ok = clusterUUIDObj.([]interface{})
		if ok {
			return "", fmt.Errorf("cluster %v is not initialized. clusterUUIDObj=%v\n", hostAddr, clusterUUIDObj)
		} else {
			return "", fmt.Errorf("uuid key in cluster info is not of string type. hostAddr=%v, clusterUUIDObj=%v\n", hostAddr, clusterUUIDObj)
		}
	}
	return clusterUUID, nil
}

// get a list of node infos with full info
// this api calls xxx/pools/nodes, which returns full node info including clustercompatibility, etc.
// the catch is that this xxx/pools/nodes is not supported by elastic search cluster
func (u *Utilities) GetNodeListWithFullInfo(hostAddr, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) ([]interface{}, error) {
	clusterInfo, err := u.GetClusterInfo(hostAddr, base.NodesPath, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return nil, err
	}

	return u.GetNodeListFromInfoMap(clusterInfo, logger)

}

// get a list of node infos with minimum info
// this api calls xxx/pools/default, which returns a subset of node info such as hostname
// this api can/needs to be used when connecting to elastic search cluster, which supports xxx/pools/default
func (u *Utilities) GetNodeListWithMinInfo(hostAddr, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) ([]interface{}, error) {
	clusterInfo, err := u.GetClusterInfo(hostAddr, base.DefaultPoolPath, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return nil, err
	}

	return u.GetNodeListFromInfoMap(clusterInfo, logger)

}

func (u *Utilities) GetClusterUUIDAndNodeListWithMinInfo(hostAddr, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (string, []interface{}, error) {
	defaultPoolInfo, err := u.GetClusterInfo(hostAddr, base.DefaultPoolPath, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return "", nil, err
	}

	clusterUUID, err := u.GetClusterUUIDFromDefaultPoolInfo(defaultPoolInfo, logger)
	if err != nil {
		return "", nil, err
	}

	nodeList, err := u.GetNodeListFromInfoMap(defaultPoolInfo, logger)

	return clusterUUID, nodeList, err

}

// get bucket info
// a specialized case of GetClusterInfo
func (u *Utilities) GetBucketInfo(hostAddr, bucketName, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (map[string]interface{}, error) {
	bucketInfo := make(map[string]interface{})
	err, statusCode := u.QueryRestApiWithAuth(hostAddr, base.DefaultPoolBucketsPath+bucketName, false, username, password, certificate, sanInCertificate, base.MethodGet, "", nil, 0, &bucketInfo, nil, false, logger)
	if err == nil && statusCode == http.StatusOK {
		return bucketInfo, nil
	}
	if statusCode == http.StatusNotFound {
		return nil, u.GetNonExistentBucketError()
	} else {
		logger.Errorf("Failed to get bucket info for bucket '%v'. host=%v, err=%v, statusCode=%v", bucketName, hostAddr, err, statusCode)
		return nil, fmt.Errorf("Failed to get bucket info.")
	}
}

// get bucket uuid
func (u *Utilities) BucketUUID(hostAddr, bucketName, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (string, error) {
	bucketInfo, err := u.GetBucketInfo(hostAddr, bucketName, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return "", err
	}

	return u.GetBucketUuidFromBucketInfo(bucketName, bucketInfo, logger)
}

// get bucket password
func (u *Utilities) BucketPassword(hostAddr, bucketName, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (string, error) {
	bucketInfo, err := u.GetBucketInfo(hostAddr, bucketName, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return "", err
	}

	return u.GetBucketPasswordFromBucketInfo(bucketName, bucketInfo, logger)
}

func (u *Utilities) GetLocalBuckets(hostAddr string, logger *log.CommonLogger) (map[string]string, error) {
	return u.GetBuckets(hostAddr, "", "", nil, false, logger)
}

// return a map of buckets
// key = bucketName, value = bucketUUID
func (u *Utilities) GetBuckets(hostAddr, username, password string, certificate []byte, sanInCertificate bool, logger *log.CommonLogger) (map[string]string, error) {
	bucketListInfo := make([]interface{}, 0)
	err, statusCode := u.QueryRestApiWithAuth(hostAddr, base.DefaultPoolBucketsPath, false, username, password, certificate, sanInCertificate, base.MethodGet, "", nil, 0, &bucketListInfo, nil, false, logger)
	if err != nil || statusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed on calling host=%v, path=%v, err=%v, statusCode=%v", hostAddr, base.DefaultPoolBucketsPath, err, statusCode)
	}

	return u.GetBucketsFromInfoMap(bucketListInfo, logger)
}

func (u *Utilities) GetBucketsFromInfoMap(bucketListInfo []interface{}, logger *log.CommonLogger) (map[string]string, error) {
	buckets := make(map[string]string)
	for _, bucketInfo := range bucketListInfo {
		bucketInfoMap, ok := bucketInfo.(map[string]interface{})
		if !ok {
			errMsg := fmt.Sprintf("bucket info is not of map type.  bucket info=%v", bucketInfo)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		bucketNameInfo, ok := bucketInfoMap[base.BucketNameKey]
		if !ok {
			errMsg := fmt.Sprintf("bucket info does not contain bucket name.  bucket info=%v", bucketInfoMap)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		bucketName, ok := bucketNameInfo.(string)
		if !ok {
			errMsg := fmt.Sprintf("bucket name is not of string type.  bucket name=%v", bucketNameInfo)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		bucketUUIDInfo, ok := bucketInfoMap[base.UUIDKey]
		if !ok {
			errMsg := fmt.Sprintf("bucket info does not contain bucket uuid.  bucket info=%v", bucketInfoMap)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		bucketUUID, ok := bucketUUIDInfo.(string)
		if !ok {
			errMsg := fmt.Sprintf("bucket uuid is not of string type.  bucket uuid=%v", bucketUUIDInfo)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		buckets[bucketName] = bucketUUID
	}

	return buckets, nil
}

// get a number of fields in bucket for validation purpose
// 1. bucket type
// 2. bucket uuid
// 3. bucket conflict resolution type
// 4. bucket eviction policy
// 5. bucket server vb map
func (u *Utilities) BucketValidationInfo(hostAddr, bucketName, username, password string, certificate []byte, sanInCertificate bool,
	logger *log.CommonLogger) (bucketInfo map[string]interface{}, bucketType string, bucketUUID string, bucketConflictResolutionType string,
	bucketEvictionPolicy string, bucketKVVBMap map[string][]uint16, err error) {
	bucketInfo, err = u.GetBucketInfo(hostAddr, bucketName, username, password, certificate, sanInCertificate, logger)
	if err != nil {
		return
	}

	bucketType, err = u.GetBucketTypeFromBucketInfo(bucketName, bucketInfo)
	if err != nil {
		err = fmt.Errorf("Error retrieving BucketType setting on bucket %v. err=%v", bucketName, err)
		return
	}
	bucketUUID, err = u.GetBucketUuidFromBucketInfo(bucketName, bucketInfo, logger)
	if err != nil {
		err = fmt.Errorf("Error retrieving UUID setting on bucket %v. err=%v", bucketName, err)
		return
	}
	bucketConflictResolutionType, err = u.GetConflictResolutionTypeFromBucketInfo(bucketName, bucketInfo)
	if err != nil {
		err = fmt.Errorf("Error retrieving ConflictResolutionType setting on bucket %v. err=%v", bucketName, err)
		return
	}
	bucketEvictionPolicy, err = u.GetEvictionPolicyFromBucketInfo(bucketName, bucketInfo)
	if err != nil {
		err = fmt.Errorf("Error retrieving EvictionPolicy setting on bucket %v. err=%v", bucketName, err)
		return
	}
	bucketKVVBMap, err = u.GetServerVBucketsMap(hostAddr, bucketName, bucketInfo)
	if err != nil {
		err = fmt.Errorf("Error retrieving server vb map on bucket %v. err=%v", bucketName, err)
		return
	}
	return
}

func (u *Utilities) GetBucketUuidFromBucketInfo(bucketName string, bucketInfo map[string]interface{}, logger *log.CommonLogger) (string, error) {
	bucketUUID := ""
	bucketUUIDObj, ok := bucketInfo[base.UUIDKey]
	if !ok {
		return "", fmt.Errorf("Error looking up uuid of bucket %v", bucketName)
	} else {
		bucketUUID, ok = bucketUUIDObj.(string)
		if !ok {
			return "", fmt.Errorf("Uuid of bucket %v is of wrong type", bucketName)
		}
	}
	return bucketUUID, nil
}

func (u *Utilities) GetClusterUUIDFromDefaultPoolInfo(defaultPoolInfo map[string]interface{}, logger *log.CommonLogger) (string, error) {
	bucketsObj, ok := defaultPoolInfo[base.BucketsKey]
	if !ok {
		errMsg := fmt.Sprintf("Cannot find buckets key in default pool info. defaultPoolInfo=%v\n", defaultPoolInfo)
		logger.Error(errMsg)
		return "", errors.New(errMsg)
	}
	bucketsInfo, ok := bucketsObj.(map[string]interface{})
	if !ok {
		errMsg := fmt.Sprintf("buckets in default pool info is not of map type. buckets=%v\n", bucketsObj)
		logger.Error(errMsg)
		return "", errors.New(errMsg)
	}
	uriObj, ok := bucketsInfo[base.URIKey]
	if !ok {
		errMsg := fmt.Sprintf("Cannot find uri key in buckets info. bucketsInfo=%v\n", bucketsInfo)
		logger.Error(errMsg)
		return "", errors.New(errMsg)
	}
	uri, ok := uriObj.(string)
	if !ok {
		errMsg := fmt.Sprintf("uri in buckets info is not of string type. uri=%v\n", uriObj)
		logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	return u.GetClusterUUIDFromURI(uri)
}

func (u *Utilities) GetClusterUUIDFromURI(uri string) (string, error) {
	// uri is in the form of /pools/default/buckets?uuid=d5dea23aa7ee3771becb3fcdb46ff956
	searchKey := base.UUIDKey + "="
	index := strings.LastIndex(uri, searchKey)
	if index < 0 {
		return "", fmt.Errorf("uri does not contain uuid. uri=%v", uri)
	}
	return uri[index+len(searchKey):], nil
}

func (u *Utilities) GetClusterCompatibilityFromBucketInfo(bucketName string, bucketInfo map[string]interface{}, logger *log.CommonLogger) (int, error) {
	nodeList, err := u.GetNodeListFromInfoMap(bucketInfo, logger)
	if err != nil {
		return 0, err
	}

	clusterCompatibility, err := u.GetClusterCompatibilityFromNodeList(nodeList)
	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	return clusterCompatibility, nil
}

func (u *Utilities) GetBucketPasswordFromBucketInfo(bucketName string, bucketInfo map[string]interface{}, logger *log.CommonLogger) (string, error) {
	bucketPassword := ""
	bucketPasswordObj, ok := bucketInfo[base.SASLPasswordKey]
	if !ok {
		return "", fmt.Errorf("Error looking up password of bucket %v", bucketName)
	} else {
		bucketPassword, ok = bucketPasswordObj.(string)
		if !ok {
			return "", fmt.Errorf("Password of bucket %v is of wrong type", bucketName)
		}
	}
	return bucketPassword, nil
}

func (u *Utilities) GetNodeListFromInfoMap(infoMap map[string]interface{}, logger *log.CommonLogger) ([]interface{}, error) {
	// get node list from the map
	nodes, ok := infoMap[base.NodesKey]
	if !ok {
		errMsg := fmt.Sprintf("info map contains no nodes. info map=%v", infoMap)
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	nodeList, ok := nodes.([]interface{})
	if !ok {
		errMsg := fmt.Sprintf("nodes is not of list type. type of nodes=%v", reflect.TypeOf(nodes))
		logger.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	// only return the nodes that are active
	activeNodeList := make([]interface{}, 0)
	for _, node := range nodeList {
		nodeInfoMap, ok := node.(map[string]interface{})
		if !ok {
			errMsg := fmt.Sprintf("node info is not of map type. type=%v", reflect.TypeOf(node))
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		clusterMembershipObj, ok := nodeInfoMap[base.ClusterMembershipKey]
		if !ok {
			// this could happen when target is elastic search cluster (or maybe very old couchbase cluster?)
			// consider the node to be "active" to be safe
			errMsg := fmt.Sprintf("node info map does not contain cluster membership. node info map=%v ", nodeInfoMap)
			logger.Debug(errMsg)
			activeNodeList = append(activeNodeList, node)
			continue
		}
		clusterMembership, ok := clusterMembershipObj.(string)
		if !ok {
			// play safe and return the node as active
			errMsg := fmt.Sprintf("cluster membership is not string type. type=%v ", reflect.TypeOf(clusterMembershipObj))
			logger.Warn(errMsg)
			activeNodeList = append(activeNodeList, node)
			continue
		}
		if clusterMembership == "" || clusterMembership == base.ClusterMembership_Active {
			activeNodeList = append(activeNodeList, node)
		}
	}

	return activeNodeList, nil
}

func (u *Utilities) GetClusterCompatibilityFromNodeList(nodeList []interface{}) (int, error) {
	if len(nodeList) > 0 {
		firstNode, ok := nodeList[0].(map[string]interface{})
		if !ok {
			return 0, fmt.Errorf("node info is of wrong type. node info=%v", nodeList[0])
		}
		clusterCompatibility, ok := firstNode[base.ClusterCompatibilityKey]
		if !ok {
			return 0, fmt.Errorf("Can't get cluster compatibility info. node info=%v\n If replicating to ElasticSearch node, use XDCR v1.", nodeList[0])
		}
		clusterCompatibilityFloat, ok := clusterCompatibility.(float64)
		if !ok {
			return 0, fmt.Errorf("cluster compatibility is not of int type. type=%v", reflect.TypeOf(clusterCompatibility))
		}
		return int(clusterCompatibilityFloat), nil
	}

	return 0, fmt.Errorf("node list is empty")
}

func (u *Utilities) GetNodeNameListFromNodeList(nodeList []interface{}, connStr string, logger *log.CommonLogger) ([]string, error) {
	nodeNameList := make([]string, 0)

	for _, node := range nodeList {
		nodeInfoMap, ok := node.(map[string]interface{})
		if !ok {
			errMsg := fmt.Sprintf("node info is not of map type. type of node info=%v", reflect.TypeOf(node))
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		hostAddr, err := u.GetHostAddrFromNodeInfo(connStr, nodeInfoMap, logger)
		if err != nil {
			errMsg := fmt.Sprintf("cannot get hostname from node info %v", nodeInfoMap)
			logger.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		nodeNameList = append(nodeNameList, hostAddr)
	}
	return nodeNameList, nil
}

func (u *Utilities) GetHostAddrFromNodeInfo(adminHostAddr string, nodeInfo map[string]interface{}, logger *log.CommonLogger) (string, error) {
	var hostAddr string
	var ok bool
	hostAddrObj, ok := nodeInfo[base.HostNameKey]
	if !ok {
		logger.Infof("hostname is missing from node info %v. This could happen in local test env where target cluster consists of a single node, %v. Just use that node.\n", nodeInfo, adminHostAddr)
		hostAddr = adminHostAddr
	} else {
		hostAddr, ok = hostAddrObj.(string)
		if !ok {
			return "", fmt.Errorf("Error constructing ssl port map for target cluster %v. host name, %v, is of wrong type\n", hostAddr, hostAddrObj)
		}
	}

	return hostAddr, nil
}

func (u *Utilities) GetHostNameFromNodeInfo(adminHostAddr string, nodeInfo map[string]interface{}, logger *log.CommonLogger) (string, error) {
	hostAddr, err := u.GetHostAddrFromNodeInfo(adminHostAddr, nodeInfo, logger)
	if err != nil {
		return "", err
	}
	return base.GetHostName(hostAddr), nil
}

//convenient api for rest calls to local cluster
func (u *Utilities) QueryRestApi(baseURL string,
	path string,
	preservePathEncoding bool,
	httpCommand string,
	contentType string,
	body []byte,
	timeout time.Duration,
	out interface{},
	logger *log.CommonLogger) (error, int) {
	return u.QueryRestApiWithAuth(baseURL, path, preservePathEncoding, "", "", nil, false, httpCommand, contentType, body, timeout, out, nil, false, logger)
}

func (u *Utilities) EnforcePrefix(prefix string, str string) string {
	var ret_str string = str
	if !strings.HasPrefix(str, prefix) {
		ret_str = prefix + str
	}
	return ret_str
}

func (u *Utilities) RemovePrefix(prefix string, str string) string {
	ret_str := strings.Replace(str, prefix, "", 1)
	return ret_str
}

//this expect the baseURL doesn't contain username and password
//if username and password passed in is "", assume it is local rest call,
//then call cbauth to add authenticate information
func (u *Utilities) QueryRestApiWithAuth(
	baseURL string,
	path string,
	preservePathEncoding bool,
	username string,
	password string,
	certificate []byte,
	san_in_certificate bool,
	httpCommand string,
	contentType string,
	body []byte,
	timeout time.Duration,
	out interface{},
	client *http.Client,
	keep_client_alive bool,
	logger *log.CommonLogger) (error, int) {
	http_client, req, err := u.prepareForRestCall(baseURL, path, preservePathEncoding, username, password, certificate, san_in_certificate, httpCommand, contentType, body, client, logger)
	if err != nil {
		return err, 0
	}

	err, statusCode := u.doRestCall(req, timeout, out, http_client, logger)
	u.cleanupAfterRestCall(keep_client_alive, err, http_client, logger)

	return err, statusCode
}

func (u *Utilities) prepareForRestCall(baseURL string,
	path string,
	preservePathEncoding bool,
	username string,
	password string,
	certificate []byte,
	san_in_certificate bool,
	httpCommand string,
	contentType string,
	body []byte,
	client *http.Client,
	logger *log.CommonLogger) (*http.Client, *http.Request, error) {
	var l *log.CommonLogger = u.loggerForFunc(logger)
	var ret_client *http.Client = client
	req, host, err := u.ConstructHttpRequest(baseURL, path, preservePathEncoding, username, password, certificate, httpCommand, contentType, body, l)
	if err != nil {
		return nil, nil, err
	}

	if ret_client == nil {
		ret_client, err = u.GetHttpClient(certificate, san_in_certificate, host, l)
		if err != nil {
			l.Errorf("Failed to get client for request, err=%v, req=%v\n", err, req)
			return nil, nil, err
		}
	}
	return ret_client, req, nil
}

func (u *Utilities) cleanupAfterRestCall(keep_client_alive bool, err error, client *http.Client, logger *log.CommonLogger) {
	if !keep_client_alive || u.IsSeriousNetError(err) {
		if client != nil && client.Transport != nil {
			transport, ok := client.Transport.(*http.Transport)
			if ok {
				if u.IsSeriousNetError(err) {
					logger.Debugf("Encountered %v, close all idle connections for this http client.\n", err)
				}
				transport.CloseIdleConnections()
			}
		}
	}
}

func (u *Utilities) doRestCall(req *http.Request,
	timeout time.Duration,
	out interface{},
	client *http.Client,
	logger *log.CommonLogger) (error, int) {
	var l *log.CommonLogger = u.loggerForFunc(logger)
	if timeout > 0 {
		client.Timeout = timeout
	} else if client.Timeout != base.DefaultHttpTimeout {
		client.Timeout = base.DefaultHttpTimeout
	}

	res, err := client.Do(req)
	if err == nil && res != nil && res.Body != nil {
		defer res.Body.Close()
		bod, e := ioutil.ReadAll(io.LimitReader(res.Body, res.ContentLength))
		if e != nil {
			l.Infof("Failed to read response body, err=%v\n req=%v\n", e, req)
			return e, res.StatusCode
		}
		if out != nil {
			err_marshal := json.Unmarshal(bod, out)
			if err_marshal != nil {
				l.Infof("Failed to unmarshal the response as json, err=%v, bod=%v\n req=%v\n", err_marshal, bod, req)
				out = bod
			} else {
				l.Debugf("out=%v\n", out)
			}
		} else {
			l.Debugf("out is nil")
		}
		return nil, res.StatusCode
	}

	return err, 0

}

//convenient api for rest calls to local cluster
func (u *Utilities) InvokeRestWithRetry(baseURL string,
	path string,
	preservePathEncoding bool,
	httpCommand string,
	contentType string,
	body []byte,
	timeout time.Duration,
	out interface{},
	client *http.Client,
	keep_client_alive bool,
	logger *log.CommonLogger, num_retry int) (error, int, *http.Client) {
	return u.InvokeRestWithRetryWithAuth(baseURL, path, preservePathEncoding, "", "", nil, false, true, httpCommand, contentType, body, timeout, out, client, keep_client_alive, logger, num_retry)
}

func (u *Utilities) InvokeRestWithRetryWithAuth(baseURL string,
	path string,
	preservePathEncoding bool,
	username string,
	password string,
	certificate []byte,
	san_in_certificate bool,
	insecureSkipVerify bool,
	httpCommand string,
	contentType string,
	body []byte,
	timeout time.Duration,
	out interface{},
	client *http.Client,
	keep_client_alive bool,
	logger *log.CommonLogger, num_retry int) (error, int, *http.Client) {

	var http_client *http.Client = nil
	var ret_err error
	var statusCode int
	var req *http.Request = nil
	backoff_time := 500 * time.Millisecond

	for i := 0; i < num_retry; i++ {
		http_client, req, ret_err = u.prepareForRestCall(baseURL, path, preservePathEncoding, username, password, certificate, san_in_certificate, httpCommand, contentType, body, client, logger)
		if ret_err == nil {
			ret_err, statusCode = u.doRestCall(req, timeout, out, http_client, logger)
		}

		if ret_err == nil {
			break
		}

		logger.Infof("Received error when making rest call. baseURL=%v, path=%v, ret_err=%v, statusCode=%v, num_retry=%v\n", baseURL, path, ret_err, statusCode, i)

		//cleanup the idle connection if the error is serious network error
		u.cleanupAfterRestCall(true, ret_err, http_client, logger)

		//backoff
		backoff_time = backoff_time + backoff_time
		time.Sleep(backoff_time)
	}

	return ret_err, statusCode, http_client

}

func (u *Utilities) GetHttpClient(certificate []byte, san_in_certificate bool, ssl_con_str string, logger *log.CommonLogger) (*http.Client, error) {
	var client *http.Client
	if len(certificate) != 0 {
		//https
		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM(certificate)
		if !ok {
			return nil, base.InvalidCerfiticateError
		}

		//using a separate tls connection to verify certificate
		//it can be changed in 1.4 when DialTLS is avaialbe in http.Transport
		conn, tlsConfig, err := base.MakeTLSConn(ssl_con_str, certificate, san_in_certificate, logger)
		if err != nil {
			return nil, err
		}
		conn.Close()

		tr := &http.Transport{TLSClientConfig: tlsConfig, Dial: base.DialTCPWithTimeout}
		client = &http.Client{Transport: tr,
			Timeout: base.DefaultHttpTimeout}

	} else {
		client = &http.Client{Timeout: base.DefaultHttpTimeout}
	}
	return client, nil
}

func (u *Utilities) maybeAddAuth(req *http.Request, username string, password string) {
	if username != "" && password != "" {
		req.Header.Set("Authorization", "Basic "+
			base64.StdEncoding.EncodeToString([]byte(username+":"+password)))
	}
}

//this expect the baseURL doesn't contain username and password
//if username and password passed in is "", assume it is local rest call,
//then call cbauth to add authenticate information
func (u *Utilities) ConstructHttpRequest(
	baseURL string,
	path string,
	preservePathEncoding bool,
	username string,
	password string,
	certificate []byte,
	httpCommand string,
	contentType string,
	body []byte,
	logger *log.CommonLogger) (*http.Request, string, error) {
	var baseURL_new string

	//process the URL
	if len(certificate) == 0 {
		baseURL_new = u.EnforcePrefix("http://", baseURL)
	} else {
		baseURL_new = u.EnforcePrefix("https://", baseURL)
	}
	url, err := couchbase.ParseURL(baseURL_new)
	if err != nil {
		return nil, "", err
	}

	var l *log.CommonLogger = u.loggerForFunc(logger)
	var req *http.Request = nil

	if !preservePathEncoding {
		if q := strings.Index(path, "?"); q > 0 {
			url.Path = path[:q]
			url.RawQuery = path[q+1:]
		} else {
			url.Path = path
		}

		req, err = http.NewRequest(httpCommand, url.String(), bytes.NewBuffer(body))
		if err != nil {
			return nil, "", err
		}
	} else {
		// use url.Opaque to preserve encoding
		url.Opaque = "//"

		index := strings.Index(baseURL_new, "//")
		if index < len(baseURL_new)-2 {
			url.Opaque += baseURL_new[index+2:]
		}
		url.Opaque += path

		req, err = http.NewRequest(httpCommand, baseURL_new, bytes.NewBuffer(body))
		if err != nil {
			return nil, "", err
		}

		// get the original Opaque back
		req.URL.Opaque = url.Opaque
	}

	if contentType == "" {
		contentType = base.DefaultContentType
	}
	req.Header.Set(base.ContentType, contentType)

	req.Header.Set(base.UserAgent, base.GoxdcrUserAgent)

	// username is nil when calling /nodes/self/xdcrSSLPorts on target
	// other username can be nil only in local rest calls
	if username == "" && path != base.SSLPortsPath {
		err := cbauth.SetRequestAuth(req)
		if err != nil {
			l.Errorf("Failed to set authentication to request, req=%v\n", req)
			return nil, "", err
		}
	} else {
		req.SetBasicAuth(username, password)
	}

	//TODO: log request would log password barely
	l.Debugf("http request=%v\n", req)

	return req, url.Host, nil
}

// encode http request into wire format
// it differs from HttpRequest.Write() in that it preserves the Content-Length in the header,
// and ignores Body in request
func (u *Utilities) EncodeHttpRequest(req *http.Request) ([]byte, error) {
	reqBytes := make([]byte, 0)
	reqBytes = append(reqBytes, []byte(req.Method)...)
	reqBytes = append(reqBytes, []byte(" ")...)
	reqBytes = append(reqBytes, []byte(req.URL.String())...)
	reqBytes = append(reqBytes, []byte(" HTTP/1.1\r\n")...)

	hasHost := false
	for key, value := range req.Header {
		if key == "Host" {
			hasHost = true
		}
		if value != nil && len(value) > 0 {
			reqBytes = u.EncodeHttpRequestHeader(reqBytes, key, value[0])
		} else {
			reqBytes = u.EncodeHttpRequestHeader(reqBytes, key, "")
		}
	}
	if !hasHost {
		// ensure that host name is in header
		reqBytes = u.EncodeHttpRequestHeader(reqBytes, "Host", req.Host)
	}

	// add extra "\r\n" as separator for Body
	reqBytes = append(reqBytes, []byte("\r\n")...)

	if req.Body != nil {
		defer req.Body.Close()

		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		reqBytes = append(reqBytes, bodyBytes...)
	}
	return reqBytes, nil
}

func (u *Utilities) EncodeHttpRequestHeader(reqBytes []byte, key, value string) []byte {
	reqBytes = append(reqBytes, []byte(key)...)
	reqBytes = append(reqBytes, []byte(": ")...)
	reqBytes = append(reqBytes, []byte(value)...)
	reqBytes = append(reqBytes, []byte("\r\n")...)
	return reqBytes
}

func (u *Utilities) IsSeriousNetError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	netError, ok := err.(*net.OpError)
	return err == syscall.EPIPE ||
		err == io.EOF ||
		strings.Contains(errStr, "EOF") ||
		strings.Contains(errStr, "use of closed network connection") ||
		strings.Contains(errStr, "connection reset by peer") ||
		strings.Contains(errStr, "http: can't write HTTP request on broken connection") ||
		(ok && (!netError.Temporary() && !netError.Timeout()))
}

func (u *Utilities) NewTCPConn(hostName string) (*net.TCPConn, error) {
	conn, err := base.DialTCPWithTimeout(base.NetTCP, hostName)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, fmt.Errorf("Failed to set up connection to %v", hostName)
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		// should never get here
		conn.Close()
		return nil, fmt.Errorf("The connection to %v returned is not TCP type", hostName)
	}

	// same settings as erlang xdcr
	err = tcpConn.SetKeepAlive(true)
	if err == nil {
		err = tcpConn.SetKeepAlivePeriod(base.KeepAlivePeriod)
	}
	if err == nil {
		err = tcpConn.SetNoDelay(false)
	}

	if err != nil {
		tcpConn.Close()
		return nil, fmt.Errorf("Error setting options on the connection to %v. err=%v", hostName, err)
	}

	return tcpConn, nil
}

/**
 * Executes a anonymous function that returns an error. If the error is non nil, retry with exponential backoff.
 * Returns base.ErrorFailedAfterRetry if operation times out, nil otherwise.
 * Max retries == the times to retry in additional to the initial try, should the initial try fail
 * initialWait == Initial time with which to start
 * Factor == exponential backoff factor based off of initialWait
 */
func (u *Utilities) ExponentialBackoffExecutor(name string, initialWait time.Duration, maxRetries int, factor int, op ExponentialOpFunc) error {
	waitTime := initialWait
	for i := 0; i <= maxRetries; i++ {
		if op() == nil {
			return nil
		} else if i != maxRetries {
			u.logger_utils.Warnf("ExponentialBackoffExecutor for %v encountered error. Sleeping %v\n",
				name, waitTime)
			time.Sleep(waitTime)
			waitTime *= time.Duration(factor)
		}
	}
	return base.ErrorFailedAfterRetry
}

/*
 * This method has an additional parameter, finCh, than ExponentialBackoffExecutor. When finCh is closed,
 * this method can abort earlier.
 */
func (u *Utilities) ExponentialBackoffExecutorWithFinishSignal(name string, initialWait time.Duration, maxRetries int, factor int, op ExponentialOpFunc2, param interface{}, finCh chan bool) (interface{}, error) {
	waitTime := initialWait
	for i := 0; i <= maxRetries; i++ {
		select {
		case <-finCh:
			u.logger_utils.Warnf("ExponentialBackoffExecutorWithFinishSignal for %v aborting because of finch closure\n", name)
			break
		default:
			result, err := op(param)
			if err == nil {
				return result, nil
			} else if i != maxRetries {
				u.logger_utils.Warnf("ExponentialBackoffExecutorWithFinishSignal for %v encountered error. Sleeping %v\n",
					name, waitTime)
				base.WaitForTimeoutOrFinishSignal(waitTime, finCh)
				waitTime *= time.Duration(factor)
			}
		}
	}
	return nil, base.ErrorFailedAfterRetry
}
