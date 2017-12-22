package main
import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strconv"
	"math"
	"strings"
)
func cnvFloat(str string) float64{
	if !(strings.Contains(str,"+")||strings.Contains(str,"-")) {
		A,err:=strconv.ParseFloat(strings.TrimSpace(str),64)
		if err!= nil {
			return 0
		}
		return A
	}
	var buf = []byte(strings.TrimSpace(str))
	sign := 1.
	for i,s := range buf {
		if (s=='-'||s=='+')&&i!=0 {
			if s=='-' { sign=-1. }
			A,err:=strconv.ParseFloat(string(buf[:i]),64)
			if err != nil { return 0 }
			B,err:=strconv.ParseFloat(string(buf[i+1:]),64)
			if err != nil { return 0 }
			return A*math.Pow(10.,sign*B)
		}
	}
	return 0.
}
func getv(str string,ind int) float64{
	pos:=[]uint{0,11,22,33,44,55,66,70,72,75,80}
	if ind>(len(pos)) {
		return 0
	}
	return cnvFloat(str[pos[ind]:pos[ind+1]])
}
func parsev(str string) []float64{
	res:=make([]float64,0,10)
	for i:=0;i<10;i++ {
		res=append(res,getv(str,i))
	}
	return res
}
func main(){
	if len(os.Args) == 1 {
		fmt.Println("usage : ",os.Args[0]," [filename] ([output])");
		return
	}
	file,err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("input file open error.")
		log.Fatal(err)
		return
	}
	defer file.Close()
	ofn := func() string{
		if len(os.Args) == 2 {
			return "scan.dat"
		}
		return os.Args[2]
	}()
	ofs,err := os.Create(ofn)
	if err != nil {
		fmt.Println("output file open error.")
		ofs,err = os.Create("/dev/stdout")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	defer ofs.Close()
	scanner := bufio.NewScanner(file)
	is_ef6 := func() bool {
		for i:=0;i<20;i++ {
			scanner.Scan()
			if strings.Contains(scanner.Text(),"------ENDF-6 FORMAT") {
				return true
			}
		}
		return false
	}()
	if !is_ef6 {
		fmt.Println("This file may not a ENDF-6 format file")
		fmt.Println("exit")
		return
	}
	for scanner.Scan() {
		vals := parsev(scanner.Text())
		mf:=int(vals[7])
		mt:=int(vals[8])
		ln:=int(vals[9])
		if mf==2&&mt==151 {
			if ln==5 {
				pari := -2.*(vals[2]-0.5)
				numE := int(vals[5])
				for i:=0;i<numE;i++ {
					scanner.Scan()
					resdata := parsev(scanner.Text())
					fmt.Fprintf(ofs,"%e\t%e\t%e\t%e\t%e\n",resdata[0],pari*resdata[1],resdata[2],resdata[3],resdata[4])
				}
				scanner.Scan()
				vals =parsev(scanner.Text())
				mf   =int(vals[7])
				mt   =int(vals[8])
				pari = -2.*(vals[2]-0.5)
				numE =int(vals[5])
				if mf!=2||mt!=151 {
					break
				}
				for i:=0;i<numE;i++ {
					scanner.Scan()
					resdata := parsev(scanner.Text())
					fmt.Fprintf(ofs,"%e\t%e\t%e\t%e\t%e\n",resdata[0],pari*resdata[1],resdata[2],resdata[3],resdata[4])
				}
				break
			}
		}
	}
	return
}
