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

var (
	ErrNotFound            = errors.New("Not fount")
	ErrInvalidNum          = errors.New("Invalid number")
	ErrInvalidAccount      = errors.New("Invalid account")
	ErrInvalidToken        = errors.New("Invalid token")
	ErrInvalidClient       = errors.New("Invalid client")
	ErrInvalidDevice       = errors.New("Invalid device")
	ErrInvalidParams       = errors.New("Invalid params")
	ErrInvalidData         = errors.New("Invalid data")
	ErrInvalidState        = errors.New("Invalid state")
	ErrTagOffline          = errors.New("Target offline")
	ErrClientOffline       = errors.New("Client offline")
	ErrDupRegister         = errors.New("Duplicated registration")
	ErrDupLogin            = errors.New("Duplicated admin login")
	ErrDupData             = errors.New("Duplicated data")
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
	ErrFileNotFound        = errors.New("File not found")
	ErrInternalServer      = errors.New("Internal server error")
	ErrDownloadFile        = errors.New("Failed download file")
	ErrCreateByte          = errors.New("Failed create bytes: system protection")
	ErrAlreadyConn         = errors.New("Already connected")
	ErrOpenSourceFile      = errors.New("Failed open source file")
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
	ErrUnsupportedFile     = errors.New("Unsupported file format")
	ErrInvalidConfigs      = errors.New("Invalid config datas")
	ErrInvalidRedisOptions = errors.New("Invalid redis options")
	ErrUnexistRedisKey     = errors.New("Unexist redis key")
	ErrNoAssociatedExpire  = errors.New("No associated expire")
	ErrUnsupportFormat     = errors.New("Unsupported format data")
	ErrInvalidOptions      = errors.New("Invalid options")
	ErrUnexistKey          = errors.New("Unexist key")
	ErrInvaildExecTime     = errors.New("Invaild execute time")
	ErrLifecycleUnexist    = errors.New("The lifecycle configuration does not exist")
)

var (
	WErrNotFound            = &WingErr{0x1000, ErrNotFound}
	WErrInvalidNum          = &WingErr{0x1001, ErrInvalidNum}
	WErrInvalidAccount      = &WingErr{0x1002, ErrInvalidAccount}
	WErrInvalidToken        = &WingErr{0x1003, ErrInvalidToken}
	WErrInvalidClient       = &WingErr{0x1004, ErrInvalidClient}
	WErrInvalidDevice       = &WingErr{0x1005, ErrInvalidDevice}
	WErrInvalidParams       = &WingErr{0x1006, ErrInvalidParams}
	WErrInvalidData         = &WingErr{0x1007, ErrInvalidData}
	WErrInvalidState        = &WingErr{0x1008, ErrInvalidState}
	WErrTagOffline          = &WingErr{0x1009, ErrTagOffline}
	WErrClientOffline       = &WingErr{0x100A, ErrClientOffline}
	WErrDupRegister         = &WingErr{0x100B, ErrDupRegister}
	WErrDupLogin            = &WingErr{0x100C, ErrDupLogin}
	WErrDupData             = &WingErr{0x100D, ErrDupData}
	WErrTokenExpired        = &WingErr{0x100E, ErrTokenExpired}
	WErrBadPublicKey        = &WingErr{0x100F, ErrBadPublicKey}
	WErrBadPrivateKey       = &WingErr{0x1010, ErrBadPrivateKey}
	WErrUnkownCharType      = &WingErr{0x1011, ErrUnkownCharType}
	WErrUnperparedState     = &WingErr{0x1012, ErrUnperparedState}
	WErrOrmNotUsing         = &WingErr{0x1013, ErrOrmNotUsing}
	WErrNoneRowFound        = &WingErr{0x1014, ErrNoneRowFound}
	WErrNotChanged          = &WingErr{0x1015, ErrNotChanged}
	WErrNotInserted         = &WingErr{0x1016, ErrNotInserted}
	WErrSendFailed          = &WingErr{0x1017, ErrSendFailed}
	WErrAuthDenied          = &WingErr{0x1018, ErrAuthDenied}
	WErrKeyLenSixteen       = &WingErr{0x1019, ErrKeyLenSixteen}
	WErrOverTimes           = &WingErr{0x101A, ErrOverTimes}
	WErrSetFrameNil         = &WingErr{0x101B, ErrSetFrameNil}
	WErrOperationNotSupport = &WingErr{0x101C, ErrOperationNotSupport}
	WErrSendHeadBytes       = &WingErr{0x101D, ErrSendHeadBytes}
	WErrSendBodyBytes       = &WingErr{0x101E, ErrSendBodyBytes}
	WErrReadBytes           = &WingErr{0x101F, ErrReadBytes}
	WErrFileNotFound        = &WingErr{0x1020, ErrFileNotFound}
	WErrInternalServer      = &WingErr{0x1021, ErrInternalServer}
	WErrDownloadFile        = &WingErr{0x1022, ErrDownloadFile}
	WErrCreateByte          = &WingErr{0x1023, ErrCreateByte}
	WErrAlreadyConn         = &WingErr{0x1024, ErrAlreadyConn}
	WErrOpenSourceFile      = &WingErr{0x1025, ErrOpenSourceFile}
	WErrEmptyReponse        = &WingErr{0x1026, ErrEmptyReponse}
	WErrReadConf            = &WingErr{0x1027, ErrReadConf}
	WErrUnexpectedDir       = &WingErr{0x1028, ErrUnexpectedDir}
	WErrWriteMD5            = &WingErr{0x1029, ErrWriteMD5}
	WErrWriteOut            = &WingErr{0x102A, ErrWriteOut}
	WErrHandleDownload      = &WingErr{0x102B, ErrHandleDownload}
	WErrFullConnPool        = &WingErr{0x102C, ErrFullConnPool}
	WErrPoolSize            = &WingErr{0x102D, ErrPoolSize}
	WErrPoolFull            = &WingErr{0x102E, ErrPoolFull}
	WErrCheckDB             = &WingErr{0x102F, ErrCheckDB}
	WErrFetchDB             = &WingErr{0x1030, ErrFetchDB}
	WErrReadFileBody        = &WingErr{0x1031, ErrReadFileBody}
	WErrNilFrame            = &WingErr{0x1032, ErrNilFrame}
	WErrNoStorage           = &WingErr{0x1033, ErrNoStorage}
	WErrUnmatchLen          = &WingErr{0x1034, ErrUnmatchLen}
	WErrCopyFile            = &WingErr{0x1035, ErrCopyFile}
	WErrEmptyData           = &WingErr{0x1036, ErrEmptyData}
	WErrImgOverSize         = &WingErr{0x1037, ErrImgOverSize}
	WErrAudioOverSize       = &WingErr{0x1038, ErrAudioOverSize}
	WErrVideoOverSize       = &WingErr{0x1039, ErrVideoOverSize}
	WErrUnsupportedFile     = &WingErr{0x103A, ErrUnsupportedFile}
	WErrInvalidConfigs      = &WingErr{0x103B, ErrInvalidConfigs}
	WErrInvalidRedisOptions = &WingErr{0x103C, ErrInvalidRedisOptions}
	WErrUnexistRedisKey     = &WingErr{0x103D, ErrUnexistRedisKey}
	WErrNoAssociatedExpire  = &WingErr{0x103E, ErrNoAssociatedExpire}
	WErrUnsupportFormat     = &WingErr{0x1040, ErrUnsupportFormat}
	WErrInvalidOptions      = &WingErr{0x1041, ErrInvalidOptions}
	WErrUnexistKey          = &WingErr{0x1042, ErrUnexistKey}
	WErrInvaildExecTime     = &WingErr{0x1043, ErrInvaildExecTime}
	WErrLifecycleUnexist    = &WingErr{0x1044, ErrLifecycleUnexist}
)

// Equal tow error if message same on char case
func EqualError(a, b error) bool {
	return a.Error() == b.Error()
}

// Equal tow error if message same ignoral char case
func EqualErrorFold(a, b error) bool {
	return strings.EqualFold(a.Error(), b.Error())
}

// Check if error message contain given string
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
