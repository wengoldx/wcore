// Copyright (c) 2018-2028 Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2022/06/07   jidi           New version
// -------------------------------------------------------------------

package comm

import (
	"github.com/go-ego/gse"
	"unicode"
)

type WordSegment struct {
	segment gse.Segmenter
}

var GseSegmenter *WordSegment

func SegmentHelper() *WordSegment {
	GseSegmenter = &WordSegment{}
	GseSegmenter.segment.LoadDict("source/dictionary/s_1.txt, source/dictionary/t_1.txt")
	return GseSegmenter
}

// CutWord split clothes title, save keywords and set keyword weight
func (c *WordSegment) CutWord(params string) []string {
	key_words := c.filterSpaceSymbols(params)
	words := c.segment.CutSearch(key_words, true)
	return words
}

// filterKeyWord filter the spaces and symbols in the string array
func (c *WordSegment) filterSpaceSymbols(key_words string) string {
	words := []rune{}
	for _, v := range key_words {
		if unicode.IsSpace(v) || unicode.IsPunct(v) {
			continue
		}
		words = append(words, v)
	}
	return string(words)
}
