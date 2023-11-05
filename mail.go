package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPX, sprRecord, hasDMARC, dmarcRecord, smtpServer\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error:", err)
	}
}

func checkDomain(domain string) {
	var hasMX, hasSPF, hasDMARC, smtpTest bool
	var spfRecord, dmarcRecord string

	mxsRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxsRecords) > 0 {
		hasMX = true
		smtpTest = checkSmtp(mxsRecords[0].Host)
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range txtRecords {
		if record[:3] == "v=sp" {
			hasSPF = true
			spfRecord = record
			break
		}

		// if record[:6] == "v=DMAR" {
		// 	hasDMARC = true
		// 	dmarcRecord = record
		// }
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if record[:8] == "v=DMARC1" {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%s, %t, %t, %s, %t, %s, %t\n", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord, smtpTest)
}

func checkSmtp(host string) bool {
	c, err := smtp.Dial(host + ":25")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	defer c.Close()

	err = c.Hello("check.local")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}

	return true
}
