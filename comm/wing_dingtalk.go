// Copyright (c) 2018-2019 WING All Rights Reserved.
//
// Author : yangping
// Email  : youhei_yp@163.com
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package comm

import (
	"encoding/json"
	"fmt"
	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
	"github.com/wengoldx/wing/secure"
	"strings"
	"time"
)

const (
	DTalkMsgText       = "text"       // text message content type
	DTalkMsgLink       = "link"       // link message content type
	DTalkMsgMarkdown   = "markdown"   // markdown message content type
	DTalkMsgActionCard = "actionCard" // action card message content type
	DTalkFeedCard      = "feedCard"   // feed card message content type
)

// -------------------------------------------------------------------
// WARNING :
// Do NOT change the json labels of below structs, it must
// same as DingTalk offical APIs define.
// -------------------------------------------------------------------

type DTAt struct {
	Mobiles []string `json:"atMobiles"`
	UserIDs []string `json:"atUserIds"`
	AtAll   bool     `json:"isAtAll"`
}
type DTButton struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}

type DTText struct {
	Content string `json:"content"`
}

type DTMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DTActionCard struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	BtnLayer    string `json:"btnOrientation"`
	SingleTitle string `json:"singleTitle"`
	SingleURL   string `json:"singleURL"`
}

type DTSplitAction struct {
	Title    string     `json:"title"`
	Text     string     `json:"text"`
	BtnLayer string     `json:"btnOrientation"`
	Btns     []DTButton `json:"btns"`
}

type DTLink struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	PicURL string `json:"picUrl"`
	MsgURL string `json:"messageUrl"`
}

type DTFeedLink struct {
	Title  string `json:"title"`
	PicURL string `json:"picURL"`
	MsgURL string `json:"messageURL"`
}

type DTFeedCard struct {
	Links []DTFeedLink `json:"links"`
}

// -------------------------------------------------------------------

// DTMsgText Text type message
type DTMsgText struct {
	At      DTAt   `json:"at"`
	Text    DTText `json:"text"`
	MsgType string `json:"msgtype"`
}

// DTMsgLink Link type message
type DTMsgLink struct {
	Link    DTLink `json:"link"`
	MsgType string `json:"msgtype"`
}

// DTMsgMarkdown Markdown type message
type DTMsgMarkdown struct {
	Text    DTMarkdown `json:"markdown"`
	At      DTAt       `json:"at"`
	MsgType string     `json:"msgtype"`
}

// DTMsgActionCard Action card type message with one click action
type DTMsgActionCard struct {
	Text    DTActionCard `json:"actionCard"`
	MsgType string       `json:"msgtype"`
}

// DTMsgSplitAction Action card type message with split button actions
type DTMsgSplitAction struct {
	Text    DTSplitAction `json:"actionCard"`
	MsgType string        `json:"msgtype"`
}

// DTMsgFeedCard Feed card type message
type DTMsgFeedCard struct {
	Card    DTFeedCard `json:"feedCard"`
	MsgType string     `json:"msgtype"`
}

// -------------------------------------------------------------------

// DTalkSender message sender for DingTalk custom robot, it just support
// inited with keywords, secure token functions, but not ips range sets.
//
// WARNING :
//
// Notice that the sender may not success @ sameone or members of group chat
// when using DingTalk user ids and the target robot have no enterprise ownership,
// so recommend use DingTalk user phone number to @ sameone or members when
// you not ensure the robot if have enterprise ownership.
//
// USAGES :
//
// the below only show send text type messages usages, the others type message as same.
// see more access link https://developers.dingtalk.com/document/robots/custom-robot-access
//
// [CODE:]
//	sender := comm.DTalkSender{
//		WebHook: "https://xxxx", Keyword: "KEY to filter message", Secure: "secure token"
//	}
//
//	// at sameone or members of group chat by user phone number
//	atMobiles := []string{"130xxxxxxxx","150xxxxxxxx"}
//
//	// at sameone or members of group chat by user id
//	atUserIds := []string{"userid1","userid2"}
//
//	// Usage 1 :
//	// send text message filter by keyword without at sameone
//	sender.SendText("message content", nil, nil, false)
//
//	// Usage 2 :
//	// send text message filter by keyword and at group chat members
//	sender.SendText("message content", nil, atUserIds, false)
//	sender.SendText("message content", atMobiles, nil, false)
//	sender.SendText("message content", atMobiles, atUserIds, false)
//
//	// Usage 3 :
//	// send text message filter by keyword and at all group chat members
//	sender.SendText("message content", nil, nil, true)
//
//	// Usage 4 :
//	// send text message with secure token
//	sender.SendText("message content", atMobiles, atUserIds, false, true)
//	sender.SendText("message content", nil, nil, true, true)
//	sender.SendText("message content", nil, nil, false, true)
// [:CODE]
type DTalkSender struct {
	WebHook string // custom group chat robot access webhook
	Keyword string // message content keyword filter
	Secure  string // robot secure signture
}

// SetSecure set DingTalk sender secure signture, it may remove
// all leading and trailing white space.
func (s *DTalkSender) SetSecure(secure string) {
	s.Secure = strings.TrimSpace(secure)
}

// UsingKey using keyword to check message content if valid, , it
// may remove all leading and trailing white space.
func (s *DTalkSender) UsingKey(keyword string) {
	s.Keyword = strings.TrimSpace(keyword)
}

// signURL sign timestamp and signture datas with send webhook
func (s *DTalkSender) signURL() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	signstr := fmt.Sprintf("%d\n%s", timestamp, s.Secure)
	signtrue := secure.SignSHA256(s.Secure, signstr)
	return fmt.Sprintf("%s&timestamp=%d&sign=%s", s.WebHook, timestamp, signtrue)
}

// checkKeyAndURL sign post url when using secure, or check keyword
// from message content if using keywords filter.
func (s *DTalkSender) checkKeyAndURL(content string, isSecure ...bool) (string, string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		logger.E("Empty message content, abort send!")
		return "", "", invar.ErrInvalidData
	}

	posturl := s.WebHook
	if len(isSecure) > 0 && isSecure[0] {
		posturl = s.signURL()
	} else if s.Keyword == "" || !strings.Contains(content, s.Keyword) {
		logger.E("Empty keyword, or not found keyword in message content!")
		return "", "", invar.ErrInvalidToken
	}
	return content, posturl, nil
}

// send send given message and check response result
func (s *DTalkSender) send(posturl string, data interface{}) error {
	resp, err := HttpPost(posturl, data)
	if err != nil {
		logger.E("Failed send text message to DingTalk group chat")
		return invar.ErrSendFailed
	}

	result := &struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(resp, result); err != nil {
		logger.E("Failed unmarshal send result:", string(resp))
		return err
	}

	logger.D("Send text message result:{code:", result.Errcode, "msg:", result.Errmsg, "}")
	if strings.ToLower(strings.TrimSpace(result.Errmsg)) != "ok" {
		return invar.ErrSendFailed
	}
	return nil
}

// SendText send text message, it support at sameone or members of group chat.
//
//	{
//		"at": {
//			"atMobiles": [ "180xxxxxx" ],
//			"atUserIds": [ "user123" ],
//			"isAtAll": false
//		},
//		"text": { "content": "the weather is nice today" },
//		"msgtype": "text"
// }
func (s *DTalkSender) SendText(content string, atMobiles, atUserIDs []string, isAtAll bool, isSecure ...bool) error {
	msg, posturl, err := s.checkKeyAndURL(content, isSecure...)
	if err != nil {
		return err
	}

	if atMobiles == nil {
		atMobiles = []string{}
	}
	if atUserIDs == nil {
		atUserIDs = []string{}
	}

	logger.D("Send text type message")
	return s.send(posturl, &DTMsgText{
		At:      DTAt{Mobiles: atMobiles, UserIDs: atUserIDs, AtAll: isAtAll},
		Text:    DTText{Content: msg},
		MsgType: DTalkMsgText,
	})
}

// SendLink send link message, it not support at anyone but have a picture and web link.
//
//	{
//		"msgtype": "link",
//		"link": {
//			"text": "the weather is nice today",
//			"title": "Hellow",
//			"picUrl": "https://link/picture.png",
//			"messageUrl": "https://link/message/url"
//		}
//	}
func (s *DTalkSender) SendLink(title, text, picURL, msgURL string, isSecure ...bool) error {
	if title == "" || text == "" {
		logger.E("Empty title or text in link message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	logger.D("Send link type message")
	return s.send(posturl, &DTMsgLink{
		Link:    DTLink{Title: title, Text: text, PicURL: picURL, MsgURL: msgURL},
		MsgType: DTalkMsgLink,
	})
}

// SendMarkdown send markdown type message, it support anyone and pick, message link urls.
//
//	{
//		"msgtype": "markdown",
//		"markdown": {
//			"title": "Hellow",
//			"text": "### the weather is nice today \n > yes"
//		},
//		"at": {
//			"atMobiles": [ "150XXXXXXXX" ],
//			"atUserIds": [ "user123" ],
//			"isAtAll": false
//		}
//	}
func (s *DTalkSender) SendMarkdown(title, text string, atMobiles, atUserIds []string, isAtAll bool, isSecure ...bool) error {
	if title == "" || text == "" {
		logger.E("Empty title or text in markdown message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	logger.D("Send markdown type message")
	return s.send(posturl, &DTMsgMarkdown{
		Text:    DTMarkdown{Title: title, Text: text},
		At:      DTAt{Mobiles: atMobiles, UserIDs: atUserIds, AtAll: isAtAll},
		MsgType: DTalkMsgMarkdown,
	})
}

// SendActionCard send action card type message, it not support at anyone but has a single link.
//
//	{
//		"actionCard": {
//			"title": "Hellow",
//			"text": "the weather is nice today",
//			"btnOrientation": "0",
//			"singleTitle" : "Click to chat",
//			"singleURL" : "https://actioncard/single/url"
//		},
//		"msgtype": "actionCard"
//	}
func (s *DTalkSender) SendActionCard(title, text, singleTitle, singleURL string, isSecure ...bool) error {
	if title == "" || text == "" || singleTitle == "" || singleURL == "" {
		logger.E("Empty input params in action card message")
		return invar.ErrInvalidData
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	logger.D("Send action card type message with single action")
	return s.send(posturl, &DTMsgActionCard{
		Text:    DTActionCard{Title: title, Text: text, BtnLayer: "0", SingleTitle: singleTitle, SingleURL: singleURL},
		MsgType: DTalkMsgActionCard,
	})
}

// SendActionCard2 send action card type message with multiple buttons.
//
//	{
//		"msgtype": "actionCard",
//		"actionCard": {
//			"title": "Hellow",
//			"text": "the weather is nice today",
//			"btnOrientation": "0",
//			"btns": [
//				{ "title": "Others",   "actionURL": "https://actioncard/other/url" },
//				{ "title": "See more", "actionURL": "https://actioncard/more/url"  }
//			]
//		}
//	}
func (s *DTalkSender) SendActionCard2(title, text string, btns []DTButton, isVertical bool, isSecure ...bool) error {
	if title == "" || text == "" {
		logger.E("Empty title or text in action card message")
		return invar.ErrInvalidData
	}

	// check all buttons if valid
	for _, btn := range btns {
		if btn.Title == "" || btn.ActionURL == "" {
			logger.E("Invalid action card button data!")
			return invar.ErrInvalidData
		}
	}

	_, posturl, err := s.checkKeyAndURL(title+text, isSecure...)
	if err != nil {
		return err
	}

	vertical := CondiString(isVertical, "0", "1")
	logger.D("Send action card type message with multips buttons")
	return s.send(posturl, &DTMsgSplitAction{
		Text:    DTSplitAction{Title: title, Text: text, BtnLayer: vertical, Btns: btns},
		MsgType: DTalkMsgActionCard,
	})
}

// SendFeedCard send feed card type message, it not support at anyone.
//
//	{
//		"msgtype":"feedCard",
//		"feedCard": {
//			"links": [
//				{ "title": "Hellow 1", "messageURL": "https://feedcard/message/url/1", "picURL": "https://feedcard/picture1.png" },
//				{ "title": "Hellow 2", "messageURL": "https://feedcard/message/url/2", "picURL": "https://feedcard/picture2.png" }
//			]
//		}
//	}
func (s *DTalkSender) SendFeedCard(links []DTFeedLink, isSecure ...bool) error {
	titles := ""
	// check all feed links if valid
	for _, link := range links {
		if link.Title == "" || link.PicURL == "" || link.MsgURL == "" {
			logger.E("Invalid feed card link data!")
			return invar.ErrInvalidData
		}
		titles += link.Title
	}

	_, posturl, err := s.checkKeyAndURL(titles, isSecure...)
	if err != nil {
		return err
	}

	logger.D("Send feed card type message")
	return s.send(posturl, &DTMsgFeedCard{
		Card:    DTFeedCard{Links: links},
		MsgType: DTalkFeedCard,
	})
}
