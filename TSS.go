package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"
)

func sendColorRequest(color string, cfg TSSConfig) {
	requestURL := fmt.Sprintf("http://%s:%s/%s", cfg.Host, cfg.Port, color)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return
	}
	fmt.Printf("client: status code: %d\n", res.StatusCode)
}

type TSSConfig struct {
	Host string
	Port string
}

func readconfig() TSSConfig {
	content, err := ioutil.ReadFile("./tssconfig.json")
	if err != nil {
		fmt.Println("Error when opening file: ", err)
		os.Exit(1)
	}
	var cfg TSSConfig
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		fmt.Println("Error during Unmarshal(): ", err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	lastStatus := ""
	hitEoF := false

	var names = map[string]string{
		"Available":       "Green",
		"Away":            "Yellow",
		"BeRightBack":     "Orange",
		"Busy":            "Orange",
		"ConnectionError": "Gray",
		"DoNotDisturb":    "Red",
		"InAMeeting":      "Orange",
		// "NewActivity":Color{12,12,12,255},
		"Offline":    "Gray",
		"OnThePhone": "Red",
		"Presenting": "Red",
		"Unknown":    "Gray",
	}

	builder := strings.Builder{}
	builder.WriteString("Added (")
	first := true
	for key := range names {
		if !first {
			builder.WriteString("|")
		} else {
			first = false
		}
		builder.WriteString(key)
	}
	builder.WriteString(")")
	regexStr := builder.String()
	fmt.Println(regexStr)
	re := regexp.MustCompile(regexStr) //`Added ([a-zA-Z]+)`

	cfg := readconfig()
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("error finding appdata folder: %s\n", err)
		os.Exit(1)
	}
	file, err := os.OpenFile(configDir+"\\Microsoft\\Teams\\logs.txt", os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err != nil {
		fmt.Printf("error reading teams log: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// without this sleep you would hogg the CPU
				time.Sleep(500 * time.Millisecond)
				if !hitEoF {
					hitEoF = true
					fmt.Printf("%q\n", lastStatus)
					sendColorRequest(names[lastStatus], cfg)
				}
				// truncated ?
				truncated, errTruncated := isTruncated(file)
				if errTruncated != nil {
					break
				}
				if truncated {
					// seek from start
					_, errSeekStart := file.Seek(0, io.SeekStart)
					if errSeekStart != nil {
						break
					}
				}
				continue
			}
			break
		}

		if strings.Contains(line, "StatusIndicatorStateService: Added") {
			submatches := re.FindStringSubmatch(line)
			if submatches == nil {
				continue
			}
			lastStatus = submatches[1]
			if hitEoF {
				fmt.Printf("%q > %q\n", lastStatus, names[lastStatus])
				sendColorRequest(names[lastStatus], cfg)
			}
			// fmt.Printf("%s\n", string(line))
		}
	}
}

func isTruncated(file *os.File) (bool, error) {
	// current read position in a file
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	// file stat to get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return currentPos > fileInfo.Size(), nil
}
