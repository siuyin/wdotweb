// wdotweb - wdot over the web
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"siuyin/log"
	"strings"
)

// lg - global logging variable
var lg = log.New(log.Ldebug)

const egWDotSrc = `:head TB "Title" {
a "Alpha"
b
  a->b
}`

type WDot struct{}

func (wd WDot) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {
	lg.Debug("WDot: " + fmt.Sprint(r.Method, r.URL.Path))
	dat := map[string]string{"title": "wdot web"}
	dat["wdotsrc"] = egWDotSrc
	err := r.ParseForm()
	if err != nil {
		lg.Error(err.Error())
	}
	if len(r.Form["wdotsrc"]) > 0 {
		lg.Debug("writing wdocsrc")
		dat["wdotsrc"] = r.Form["wdotsrc"][0]
		ioutil.WriteFile("public/wdotsrc.wdot", []byte(r.Form["wdotsrc"][0]), 0600)
		output, err := exec.Command("./wdotwebgen").CombinedOutput()
		fmt.Println(string(output))
		if err != nil || strings.Contains(strings.ToLower(string(output)), "error") {
			//lg.Error(err.Error())
			exec.Command("cp", "public/error.png", "public/wdot.png").Run()
			lg.Debug("copied to wdot.png")
		}
		lg.Debug("completed wdocsrc")
	}
	var tmpl *template.Template
	if tmpl == nil {
		tmpl = template.Must(template.ParseFiles("main.tmpl"))
	}
	tmpl.ExecuteTemplate(w, "Main", dat)
}

//
func mapRoot(h http.Handler, redirectTo string) http.Handler {
	lg.Debug("redirecting / to: " + redirectTo)
	return http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		lg.Debug("mapRoot: " + r.URL.Path)
		if r.URL.Path == "/" {
			lg.Debug("mapRootRedirect: " + r.URL.Path)
			http.Redirect(w, r, redirectTo, http.StatusFound)
		}
		lg.Debug("finished redirect: " + r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

/* main function
 */
func main() {
	lg.Info("starting")

	http.Handle("/", mapRoot(http.FileServer(http.Dir("public")), "/wdot"))

	// wdot controller
	var wd WDot
	http.Handle("/wdot", wd)

	err := http.ListenAndServe("0.0.0.0:8188", nil)
	if err != nil {
		lg.Fatal(err.Error())
	}
	lg.Info("ended")
}
