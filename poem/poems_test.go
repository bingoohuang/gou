package poem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePoems1(t *testing.T) {
	assert := assert.New(t)

	poems := ParsePoems("./poems_test1.txt")

	assert.Equal([]Poem{
		{
			Title:     "子夜吴歌·秋歌",
			TitleCode: "zywgqg",
			Author:    "唐·李白",
			Lines:     []string{"长安一片月，万户捣衣声。", "秋风吹不尽，总是玉关情。", "何日平胡虏，良人罢远征。"},
			LinesCode: []string{"Ca1p)Whd1s", "Qfc!jZ4ygq", "?rphlLr!yz"},
		},
	}, poems)
}

func TestParsePoems2(t *testing.T) {
	assert := assert.New(t)

	poems := ParsePoems("./poems_test2.txt")

	assert.Equal([]Poem{
		{
			Title:     "子夜吴歌·秋歌",
			TitleCode: "zywgqg",
			Author:    "唐·李白",
			Lines:     []string{"长安一片月，万户捣衣声。", "秋风吹不尽，总是玉关情。", "何日平胡虏，良人罢远征。"},
			LinesCode: []string{"Ca1p)Whd1s", "Qfc!jZ4ygq", "?rphlLr!yz"},
		},
		{
			Title:     "问刘十九",
			TitleCode: "wlsj",
			Author:    "唐·白居易",
			Lines:     []string{"绿蚁新醅酒，红泥小火炉。", "晚来天欲雪，能饮一杯无？"},
			LinesCode: []string{"Lyxp9Hnxhl", "Wlty*Ny1bw"},
		},
		{
			Title:     "凉州词",
			TitleCode: "lzc",
			Author:    "唐·王之涣",
			Lines:     []string{"黄沙远上白云间，一片孤城万仞山。", "羌笛何须怨杨柳，春风不度玉门关。"},
			LinesCode: []string{"Hsysbyj1pgcwr3", "Qd?xyylCf!dymg"},
		},
	}, poems)
}
