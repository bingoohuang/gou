package coll

import (
	"testing"
	//. "github.com/ahmetb/go-linq/v3"
	//"github.com/elliotchance/pie/pie"
	//"github.com/jucardi/go-streams/streams"
	//"github.com/stretchr/testify/assert"
	//"github.com/thoas/go-funk"
	//"github.com/wesovilabs/koazee"
)

func TestStreams(t *testing.T) {
	//fruitArray := []string{"peach", "apple", "pear", "plum", "pineapple", "banana", "kiwi", "orange"}
	//fruitsThatStartWithP := streams.FromArray(fruitArray).
	//	Filter(func(v interface{}) bool { return strings.HasPrefix(v.(string), "p") }).
	//	OrderBy(func(a interface{}, b interface{}) int { return strings.Compare(a.(string), b.(string)) }).
	//	ToArray().([]string)
	//assert.Equal(t, []string{"peach", "pear", "pineapple", "plum"}, fruitsThatStartWithP)
	//
	//streams.FromArray(fruitArray).
	//	Filter(func(v interface{}) bool { return strings.HasPrefix(v.(string), "p") }).
	//	ForEach(func(v interface{}) { fmt.Printf("%v\n", v) })
	//
	//var result []string
	//From(fruitArray).
	//	WhereT(func(s string) bool { return strings.HasPrefix(s, "p") }).
	//	OrderByT(func(s string) int { return len(s) }).
	//	ToSlice(&result)
	//
	//assert.Equal(t, []string{"pear", "plum", "peach", "pineapple"}, result)
	//
	//pieStr := pie.Strings{"Bob", "Sally", "John", "Jane"}.
	//	Filter(func(v string) bool { return strings.HasPrefix(v, "J") }).Join(",")
	//assert.Equal(t, "John,Jane", pieStr)
	//
	//funkResult := funk.Filter([]string{"pear", "plum", "peach", "pineapple"},
	//	func(v string) bool { return strings.HasPrefix(v, "pe") }).([]string)
	//assert.Equal(t, []string{"pear", "peach"}, funkResult)
	//
	//kr := koazee.StreamOf(fruitArray).Filter(func(s string) bool { return strings.HasPrefix(s, "p") }).Do().Out().Val()
	//assert.Equal(t, []string{"peach", "pear", "plum", "pineapple"}, kr)
}
