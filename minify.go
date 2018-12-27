package go_utils

import (
	"bytes"
	"fmt"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"strings"
)

func MinifyHtml(htmlString string, devMode bool) string {
	if devMode {
		return htmlString
	}

	mini := minify.New()
	mini.AddFunc("text/html", html.Minify)
	minified, _ := mini.String("text/html", htmlString)
	return minified
}

func MinifyCss(cssString string, devMode bool) string {
	if devMode {
		return cssString
	}

	mini := minify.New()
	mini.AddFunc("text/css", css.Minify)

	minifiedCss, err := mini.String("text/css", cssString)
	if err != nil {
		fmt.Println("mini css:", err.Error())
	}

	return minifiedCss
}

func MinifyJs(jsString string, devMode bool) string {
	if devMode {
		return jsString
	}

	mini := minify.New()
	mini.AddFunc("text/javascript", js.Minify)
	minifiedJs, err := mini.String("text/javascript", jsString)
	if err != nil {
		fmt.Println("mini js:", err.Error())
	}

	return minifiedJs
}

func MergeResJs(mustAsset func(name string) []byte, statics ...string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(mustAsset("res/" + static))
		scripts.Write([]byte(";"))
	}

	return scripts.String()
}

func MergeResCss(mustAsset func(name string) []byte, statics ...string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(mustAsset("res/" + static))
		scripts.Write([]byte("\n"))
	}

	return scripts.String()
}

func MergeJs(mustAsset func(name string) []byte, statics []string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(mustAsset(static))
		scripts.Write([]byte(";"))
	}

	return scripts.String()
}

func MergeCss(mustAsset func(name string) []byte, statics []string) string {
	var scripts bytes.Buffer
	for _, static := range statics {
		scripts.Write(mustAsset(static))
		scripts.Write([]byte("\n"))
	}

	return scripts.String()
}

func FilterAssetNames(assetNames []string, suffix string) []string {
	filtered := make([]string, 0)
	for _, assetName := range assetNames {
		if strings.HasSuffix(assetName, suffix) {
			filtered = append(filtered, assetName)
		}
	}

	return filtered
}

func FilterAssetNamesOrdered(assetNames []string, suffix string, orderedNames ...string) []string {
	filtered := make([]string, 0)

	for _, assetName := range orderedNames {
		if IndexOf(assetName+suffix, assetNames) >= 0 {
			filtered = append(filtered, assetName)
		}
	}
	for _, assetName := range assetNames {
		if strings.HasSuffix(assetName, suffix) && IndexOf(strings.TrimSuffix(assetName, suffix), assetNames) < 0 {
			filtered = append(filtered, assetName)
		}
	}

	return filtered
}
