package go_utils

import (
	"fmt"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
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
