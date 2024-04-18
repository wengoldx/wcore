// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/wengoldx/wing/invar"
	"github.com/wengoldx/wing/logger"
)

/*
 * ID Card Descriptions :
 *
 * card string length fixed at 18 chars, the first 17 chars must be number.
 * e.g. [42 11 26 19810328 93 5 2]
 *
 * (1).  1 ~  2 digits represent the code of province.
 * (2).  3 ~  4 digits represent the code of city.
 * (3).  5 ~  6 digits represent the code of district.
 * (4).  7 ~ 14 digits represent year, month, and day of birth.
 * (5). 15 ~ 16 digits represent the code of the local police station.
 * (6).      17 diget represent gender, odd numbers for male, other for female.
 * (7).      18 diget is the verification code, must be 0 ~9 or X char.
 */

const validateCodes = "10X98765432"

var validateWeights = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// Verify ID Card internal, just simple validate card number only
func IsVaildIDCard(card string) bool {
	card = strings.ToUpper(card)
	if cardlen := len(card); cardlen != 18 {
		logger.E("Invalid ID Card:", card, "lenght:", cardlen)
		return false
	}

	num, last := card[:17], card[17:]
	if !invar.RegexNumber(num) {
		logger.E("Not digites of ID Card:", num)
		return false
	}
	return validateCardNumbers(num, last)
}

// Return birthday as time from given ID Card string
func CardBirthday(card string) (*time.Time, error) {
	if len(card) != 18 {
		return nil, invar.ErrInvalidParams
	}

	birthday := card[6:14]
	bt, err := ParseTime(DateNoneHyphen, birthday)
	if err != nil {
		logger.E("Parse card birthday:", birthday, "err:", err)
		return nil, err
	}
	return &bt, nil
}

// Return gender from given ID Card string, true is male, false is female
func CardGender(card string) (bool, error) {
	if len(card) != 18 {
		return false, invar.ErrInvalidParams
	}

	genderMask, err := strconv.Atoi(string(card[16]))
	if err != nil {
		logger.E("Parse card:", card, "gender, err:", err)
		return false, err
	}
	return genderMask%2 == 1 /* male: 1, female: 2 */, nil
}

// Validate card number self if valide by last code char
// see more http://www.360doc.com/content/22/0112/12/74433059_1012930821.shtml
func validateCardNumbers(num, last string) bool {
	sum := 0
	for i := 0; i < 17; i++ {
		if digit, err := strconv.Atoi(num[i : i+1]); err != nil {
			logger.E("Invalid digit number at:", i, "err:", err)
			return false
		} else {
			sum += digit * validateWeights[i]
		}
	}

	index := sum % 11
	code := validateCodes[index : index+1]
	return code == last
}
