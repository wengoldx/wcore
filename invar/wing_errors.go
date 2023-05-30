// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// 00002       2019/06/30   zhaixing       Add function from godfs
// -------------------------------------------------------------------
package invar

import (
	"errors"
	"strings"
)

// WingErr const error with code
type WingErr struct {
	Code int
	Err  error
}

// WingErr const error with code
type WExErr struct {
	Code    int    `json:"code"    description:"Extend error code"`
	Message string `json:"message" description:"Extend error message"`
}

var (
	ErrNotFound            = errors.New("Not fount")
	ErrInvalidNum          = errors.New("Invalid number")
	ErrInvalidAccount      = errors.New("Invalid account")
	ErrInvalidToken        = errors.New("Invalid token")
	ErrInvalidRole         = errors.New("Invalid role")
	ErrInvalidClient       = errors.New("Invalid client")
	ErrInvalidDevice       = errors.New("Invalid device")
	ErrInvalidParams       = errors.New("Invalid params")
	ErrInvalidData         = errors.New("Invalid data")
	ErrInvalidState        = errors.New("Invalid state")
	ErrInvalidPhone        = errors.New("Invalid phone")
	ErrInvalidEmail        = errors.New("Invalid email")
	ErrInvalidOptions      = errors.New("Invalid options")
	ErrInvalidRedisOptions = errors.New("Invalid redis options")
	ErrInvalidConfigs      = errors.New("Invalid config datas")
	ErrInvaildExecTime     = errors.New("Invaild execute time")
	ErrInvalidRealname     = errors.New("Invaild realname")
	ErrTagOffline          = errors.New("Target offline")
	ErrClientOffline       = errors.New("Client offline")
	ErrDupRegister         = errors.New("Duplicated registration")
	ErrDupLogin            = errors.New("Duplicated admin login")
	ErrDupData             = errors.New("Duplicated data")
	ErrDupAccount          = errors.New("Duplicated account")
	ErrDupName             = errors.New("Duplicate name")
	ErrDupKey              = errors.New("Duplicate key")
	ErrTokenExpired        = errors.New("Token expired")
	ErrBadPublicKey        = errors.New("Invalid public key")
	ErrBadPrivateKey       = errors.New("Invalid private key")
	ErrUnkownCharType      = errors.New("Unkown chars type")
	ErrUnperparedState     = errors.New("Unperpared state")
	ErrOrmNotUsing         = errors.New("Orm not using")
	ErrNoneRowFound        = errors.New("None row found")
	ErrNotChanged          = errors.New("Not changed")
	ErrNotInserted         = errors.New("Not inserted")
	ErrSendFailed          = errors.New("Failed to send")
	ErrAuthDenied          = errors.New("Permission denied")
	ErrKeyLenSixteen       = errors.New("Require sixteen-length secret key")
	ErrOverTimes           = errors.New("Over retry times")
	ErrSetFrameNil         = errors.New("Failed clear frame meta")
	ErrOperationNotSupport = errors.New("Operation not support")
	ErrSendHeadBytes       = errors.New("Failed send head bytes")
	ErrSendBodyBytes       = errors.New("Failed send body bytes")
	ErrReadBytes           = errors.New("Error read bytes")
	ErrInternalServer      = errors.New("Internal server error")
	ErrCreateByte          = errors.New("Failed create bytes: system protection")
	ErrFileNotFound        = errors.New("File not found")
	ErrDownloadFile        = errors.New("Failed download file")
	ErrOpenSourceFile      = errors.New("Failed open source file")
	ErrAlreadyConn         = errors.New("Already connected")
	ErrEmptyReponse        = errors.New("Received empty response")
	ErrReadConf            = errors.New("Failed load config file")
	ErrUnexpectedDir       = errors.New("Expect file path not directory")
	ErrWriteMD5            = errors.New("Failed write to md5")
	ErrWriteOut            = errors.New("Failed write out")
	ErrHandleDownload      = errors.New("Failed handle download file")
	ErrFullConnPool        = errors.New("Connection pool is full")
	ErrPoolSize            = errors.New("Thread pool size value must be positive")
	ErrPoolFull            = errors.New("Pool is full, can not take any more")
	ErrCheckDB             = errors.New("Check database: failed retry many times")
	ErrFetchDB             = errors.New("Fetch database connection time out from pool")
	ErrReadFileBody        = errors.New("Failed read file content")
	ErrNilFrame            = errors.New("Frame is null")
	ErrNoStorage           = errors.New("No storage server available")
	ErrUnmatchLen          = errors.New("Unmatch download file length")
	ErrCopyFile            = errors.New("Failed copy file")
	ErrEmptyData           = errors.New("Empty data")
	ErrImgOverSize         = errors.New("Image file size over")
	ErrAudioOverSize       = errors.New("Audio file size over")
	ErrVideoOverSize       = errors.New("Video file size over")
	ErrNoAssociatedExpire  = errors.New("No associated expire")
	ErrUnsupportFormat     = errors.New("Unsupported format data")
	ErrUnsupportedFile     = errors.New("Unsupported file format")
	ErrUnexistKey          = errors.New("Unexist key")
	ErrUnexistRedisKey     = errors.New("Unexist redis key")
	ErrUnexistLifecycle    = errors.New("Unexist lifecycle configs")
	ErrSetLifecycleTag     = errors.New("Failed set file lifecycel tag")
	ErrInactiveAccount     = errors.New("Inactive status account")
)

var (
	WErrNotFound            = &WingErr{0x1000, ErrNotFound}
	WErrInvalidNum          = &WingErr{0x1001, ErrInvalidNum}
	WErrInvalidAccount      = &WingErr{0x1002, ErrInvalidAccount}
	WErrInvalidToken        = &WingErr{0x1003, ErrInvalidToken}
	WErrInvalidRole         = &WingErr{0x1004, ErrInvalidRole}
	WErrInvalidClient       = &WingErr{0x1005, ErrInvalidClient}
	WErrInvalidDevice       = &WingErr{0x1006, ErrInvalidDevice}
	WErrInvalidParams       = &WingErr{0x1007, ErrInvalidParams}
	WErrInvalidData         = &WingErr{0x1008, ErrInvalidData}
	WErrInvalidState        = &WingErr{0x1009, ErrInvalidState}
	WErrInvalidPhone        = &WingErr{0x100A, ErrInvalidPhone}
	WErrInvalidEmail        = &WingErr{0x100B, ErrInvalidEmail}
	WErrInvalidOptions      = &WingErr{0x100C, ErrInvalidOptions}
	WErrInvalidRedisOptions = &WingErr{0x100D, ErrInvalidRedisOptions}
	WErrInvalidConfigs      = &WingErr{0x100E, ErrInvalidConfigs}
	WErrInvaildExecTime     = &WingErr{0x100F, ErrInvaildExecTime}
	WErrInvalidRealname     = &WingErr{0x100F, ErrInvalidRealname}
	WErrTagOffline          = &WingErr{0x1010, ErrTagOffline}
	WErrClientOffline       = &WingErr{0x1011, ErrClientOffline}
	WErrDupRegister         = &WingErr{0x1012, ErrDupRegister}
	WErrDupLogin            = &WingErr{0x1013, ErrDupLogin}
	WErrDupData             = &WingErr{0x1014, ErrDupData}
	WErrDupAccount          = &WingErr{0x1014, ErrDupAccount}
	WErrDupName             = &WingErr{0x1014, ErrDupName}
	WErrDupKey              = &WingErr{0x1014, ErrDupKey}
	WErrTokenExpired        = &WingErr{0x1015, ErrTokenExpired}
	WErrBadPublicKey        = &WingErr{0x1016, ErrBadPublicKey}
	WErrBadPrivateKey       = &WingErr{0x1017, ErrBadPrivateKey}
	WErrUnkownCharType      = &WingErr{0x1018, ErrUnkownCharType}
	WErrUnperparedState     = &WingErr{0x1019, ErrUnperparedState}
	WErrOrmNotUsing         = &WingErr{0x101A, ErrOrmNotUsing}
	WErrNoneRowFound        = &WingErr{0x101B, ErrNoneRowFound}
	WErrNotChanged          = &WingErr{0x101C, ErrNotChanged}
	WErrNotInserted         = &WingErr{0x101D, ErrNotInserted}
	WErrSendFailed          = &WingErr{0x101E, ErrSendFailed}
	WErrAuthDenied          = &WingErr{0x101F, ErrAuthDenied}
	WErrKeyLenSixteen       = &WingErr{0x1020, ErrKeyLenSixteen}
	WErrOverTimes           = &WingErr{0x1021, ErrOverTimes}
	WErrSetFrameNil         = &WingErr{0x1022, ErrSetFrameNil}
	WErrOperationNotSupport = &WingErr{0x1023, ErrOperationNotSupport}
	WErrSendHeadBytes       = &WingErr{0x1024, ErrSendHeadBytes}
	WErrSendBodyBytes       = &WingErr{0x1025, ErrSendBodyBytes}
	WErrReadBytes           = &WingErr{0x1026, ErrReadBytes}
	WErrInternalServer      = &WingErr{0x1027, ErrInternalServer}
	WErrCreateByte          = &WingErr{0x1028, ErrCreateByte}
	WErrFileNotFound        = &WingErr{0x1029, ErrFileNotFound}
	WErrDownloadFile        = &WingErr{0x102A, ErrDownloadFile}
	WErrOpenSourceFile      = &WingErr{0x102B, ErrOpenSourceFile}
	WErrAlreadyConn         = &WingErr{0x102C, ErrAlreadyConn}
	WErrEmptyReponse        = &WingErr{0x102D, ErrEmptyReponse}
	WErrReadConf            = &WingErr{0x102E, ErrReadConf}
	WErrUnexpectedDir       = &WingErr{0x102F, ErrUnexpectedDir}
	WErrWriteMD5            = &WingErr{0x1030, ErrWriteMD5}
	WErrWriteOut            = &WingErr{0x1031, ErrWriteOut}
	WErrHandleDownload      = &WingErr{0x1032, ErrHandleDownload}
	WErrFullConnPool        = &WingErr{0x1033, ErrFullConnPool}
	WErrPoolSize            = &WingErr{0x1034, ErrPoolSize}
	WErrPoolFull            = &WingErr{0x1035, ErrPoolFull}
	WErrCheckDB             = &WingErr{0x1036, ErrCheckDB}
	WErrFetchDB             = &WingErr{0x1037, ErrFetchDB}
	WErrReadFileBody        = &WingErr{0x1038, ErrReadFileBody}
	WErrNilFrame            = &WingErr{0x1039, ErrNilFrame}
	WErrNoStorage           = &WingErr{0x103A, ErrNoStorage}
	WErrUnmatchLen          = &WingErr{0x103B, ErrUnmatchLen}
	WErrCopyFile            = &WingErr{0x103C, ErrCopyFile}
	WErrEmptyData           = &WingErr{0x103D, ErrEmptyData}
	WErrImgOverSize         = &WingErr{0x103E, ErrImgOverSize}
	WErrAudioOverSize       = &WingErr{0x103F, ErrAudioOverSize}
	WErrVideoOverSize       = &WingErr{0x1040, ErrVideoOverSize}
	WErrNoAssociatedExpire  = &WingErr{0x1041, ErrNoAssociatedExpire}
	WErrUnsupportFormat     = &WingErr{0x1042, ErrUnsupportFormat}
	WErrUnsupportedFile     = &WingErr{0x1043, ErrUnsupportedFile}
	WErrUnexistKey          = &WingErr{0x1044, ErrUnexistKey}
	WErrUnexistRedisKey     = &WingErr{0x1045, ErrUnexistRedisKey}
	WErrUnexistLifecycle    = &WingErr{0x1046, ErrUnexistLifecycle}
	WErrSetLifecycleTag     = &WingErr{0x1045, ErrSetLifecycleTag}
	WErrInactiveAccount     = &WingErr{0x1045, ErrInactiveAccount}
)

// Equal tow error if message same on char case
func EqualError(a, b error) bool {
	return a.Error() == b.Error()
}

// Equal tow error if message same ignoral char case
func EqualErrorFold(a, b error) bool {
	return strings.EqualFold(a.Error(), b.Error())
}

// Check if error message contain given error string
func ErrorContain(s, sub error) bool {
	return strings.Contains(s.Error(), sub.Error())
}

// Check if error message start given perfix
func ErrorStart(s, sub error) bool {
	return strings.HasPrefix(s.Error(), sub.Error())
}

// Check if error message start given perfix
func ErrorEnd(s, sub error) bool {
	return strings.HasSuffix(s.Error(), sub.Error())
}

// Check if error message contain given string
func IsError(e error, s string) bool {
	esu, su := strings.ToLower(e.Error()), strings.ToLower(s)
	return strings.Contains(esu, su)
}

// Create a custom extend error from given code and message
func GenWExErr(code int, message string) WExErr {
	return WExErr{Code: code, Message: message}
}

// Transform a WingErr to extend error
func ToWExErr(we WingErr) WExErr {
	return WExErr{Code: we.Code, Message: we.Err.Error()}
}
