package main

//import hexutils "github.com/rkbalgi/utils"
import "github.com/rkbalgi/go/crypto"
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/rkbalgi/go/utils"
	"reflect"
	"sync"
	//"os"
	//	"strings"
	//"strconv"
	//"bufio"
	//	"io"
)

type CommandHeader struct {
	Header      string `size:"12"`
	CommandName string `size:"2"`
	MacData     []byte `size:"4"`
}

var letter_map map[string]int

func letter_sum(month string) int {

	sum := 0
	for i := 0; i < len(month); i++ {
		sum += letter_map[month[i:i+1]]
	}

	return sum
}

func main() {

	fmt.Println("Hello World")
	letter_map = make(map[string]int, 26)
	tmp := "abcdefghijklmnopqrstuvwxyz"
	for i := 1; i < 27; i++ {
		letter_map[tmp[i-1:i]] = i
	}

	fmt.Println("january ", letter_sum("s"))

	cmd := "NC"

	fmt.Println(hex.EncodeToString([]byte(cmd)))

	what_type := [...]int{1, 2, 3}

	fmt.Println(reflect.TypeOf(what_type).Kind().String())

	//fmt.Printf("%s-%s\n",command_header.Header,command_header.CommandName)
	//fmt.Println(command_header)
	//os.Exit(0);

	var testKey []byte = []byte("helloworld111112") //67876543")
	testData, _ := hex.DecodeString("3657f432e456feda01")
	fmt.Printf("Key: %s Data:  %s\n", hex.EncodeToString(testKey), hex.EncodeToString(testData))
	encryptedData, err := crypto.EncryptTripleDesEde2(testKey, testData, crypto.Iso9797M2Padding)
	fmt.Println(utils.HexToString(encryptedData), " - ")
	decryptedData, err := crypto.DecryptTripleDesEde2(testKey, encryptedData, crypto.Iso9797M2Padding)
	if bytes.Equal(testData, decryptedData) {
		fmt.Println("Test OK")
	} else {
		fmt.Print("Test Failed")
	}
	if err == nil {
		fmt.Println(utils.HexToString(encryptedData), " - ", utils.HexToString(decryptedData))
		//fmt.Println(crypto.Iso9797M1Padding," ", crypto.Iso9797M2Padding)
	}
	fmt.Printf("%08b\n", 10)

	var goGrp *sync.WaitGroup = new(sync.WaitGroup)
	go func() {
		fmt.Println("Hello World")
		goGrp.Done()
	}()

	//goGrp.Add(1);
	goGrp.Wait()

	fmt.Println(crypto.RotN(3, "xyzragha"))

	/*var arr = [10]uint32{
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10}
	defer my_func(10)
	fmt.Printf("%T %d \n", arr, cap(arr))
	var total float32 = 0
	var i int
	var val uint32
	for i, val = range arr {
		fmt.Print(i, val, "\n")
		total += float32(val)
	}

	avg := total / float32(len(arr))
	fmt.Printf("%T - %f\n", avg, avg)*/

}

func my_func(x int) {
	fmt.Println(x)
}
