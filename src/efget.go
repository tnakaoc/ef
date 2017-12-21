package main
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net/http"
	"io/ioutil"
	"os/exec"
)
func getURL(url string,fnm string){
	ur,_ := http.Get(url)
	defer ur.Body.Close()
	buf,_ := ioutil.ReadAll(ur.Body)
	ioutil.WriteFile("./"+fnm,buf,os.ModePerm)
}
func main(){
	lanl_path := "https://t2.lanl.gov/nis/data/endf"
	lanl_html := "endfvii.1-n.html"
	getURL(lanl_path+"/"+lanl_html,"index.html")
	file,err := os.Open("index.html")
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line,"raw eval") {
			sta := strings.Index(line,"HREF=")+6
			sto := strings.Index(line,">raw")-1
			if sta<=0||sto<=0||sto<sta { continue }
			path := lanl_path+"/"+line[sta:sto]
			elem := strings.Split(path,"/")
			le   := len(elem)
			if le<2 { continue }
			fmt.Println("trying to download from: ",path)
			exec.Command("wget",path).Run()
			exec.Command("mv",elem[le-1],"tmp.download").Run()
//			getURL(path,"tmp.download")
			fmt.Println("trying to scan.")
			exec.Command("./ef6_scan","tmp.download","scan/"+elem[le-2]+elem[le-1]+".scan").Run()
		}
	}
}
