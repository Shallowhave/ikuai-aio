package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHttpInsecureSkipVerify  = false
	defaultHttpTimeout             = time.Second * 30
	defaultTimezone                = "Asia/Shanghai"
	defaultIKuaiAddr               = "http://192.168.1.1"
	defaultIKuaiUsername           = "admin"
	defaultIKuaiPassword           = "admin"
	defaultIKuaiCronSkipStart      = false
	defaultIKuaiExporterListenAddr = "0.0.0.0:8000"
	defaultIKuaiExporterDisable    = false
)

type Config struct {
	HttpInsecureSkipVerify    bool
	HttpTimeout               time.Duration
	Timezone                  *time.Location
	IKuaiAddr                 string
	IKuaiUsername             string
	IKuaiPassword             string
	IKuaiCronSkipStart        bool
	IKuaiCronCustomISPList    []*IKuaiCronCustomISP
	IKuaiCronStreamDomainList []*IKuaiCronStreamDomain
	IKuaiCronIPGroupList      []*IKuaiCronIPGROUP
	IKuaiExporterDisable      bool
	IKuaiExporterListenAddr   string
}

type IKuaiCronCustomISP struct {
	Cron    string
	Name    string
	Url     []string
	Comment string
}

type IKuaiCronStreamDomain struct {
	Cron      string
	Interface []string
	Url       []string
	SrcAddr   string
	Comment   string
}

type IKuaiCronIPGROUP struct {
	Cron    string
	Name    string
	Url     []string
	Comment string
}

var C *Config

func Load() *Config {
	if C != nil {
		return C
	}

	httpInsecureSkipVerifyStr := getEnv("HTTP_INSECURE_SKIP_VERIFY", strconv.FormatBool(defaultHttpInsecureSkipVerify))
	httpInsecureSkipVerify := httpInsecureSkipVerifyStr == "true"

	httpTimeoutStr := getEnv("HTTP_TIMEOUT", defaultHttpTimeout.String())
	httpTimeout, err := time.ParseDuration(httpTimeoutStr)
	if err != nil {
		log.Fatalln(err)
	}

	timezoneStr := getEnv("TZ", defaultTimezone)
	timezone, err := time.LoadLocation(timezoneStr)
	if err != nil {
		log.Fatalln(err)
	}

	iKuaiAddr := getEnv("IKUAI_ADDR", defaultIKuaiAddr)
	iKuaiUsername := getEnv("IKUAI_USERNAME", defaultIKuaiUsername)
	iKuaiPassword := getEnv("IKUAI_PASSWORD", defaultIKuaiPassword)

	iKuaiCronSkipStartStr := getEnv("IKUAI_CRON_SKIP_START", strconv.FormatBool(defaultIKuaiCronSkipStart))
	iKuaiCronSkipStart := iKuaiCronSkipStartStr == "true"

	iKuaiExporterListenAddr := getEnv("IKUAI_EXPORTER_LISTEN_ADDR", defaultIKuaiExporterListenAddr)
	iKuaiExporterDisableStr := getEnv("IKUAI_EXPORTER_DISABLE", strconv.FormatBool(defaultIKuaiExporterDisable))
	iKuaiExporterDisable := iKuaiExporterDisableStr == "true"

	c := &Config{
		HttpInsecureSkipVerify:  httpInsecureSkipVerify,
		HttpTimeout:             httpTimeout,
		Timezone:                timezone,
		IKuaiAddr:               iKuaiAddr,
		IKuaiUsername:           iKuaiUsername,
		IKuaiPassword:           iKuaiPassword,
		IKuaiCronSkipStart:      iKuaiCronSkipStart,
		IKuaiExporterListenAddr: iKuaiExporterListenAddr,
		IKuaiExporterDisable:    iKuaiExporterDisable,
	}

	c.matchCronCustomISP()
	c.matchCronStreamDomain()
	c.matchCronIpGroup()

	C = c

	return C
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func (c *Config) matchCronCustomISP() {
	re := regexp.MustCompile(`IKUAI_CRON_CUSTOM_ISP_(\d+)`)
	m := map[string]*IKuaiCronCustomISP{}

	for _, env := range os.Environ() {
		match := re.FindStringSubmatch(env)
		if len(match) < 2 {
			continue
		}
		id := match[1]
		key := fmt.Sprintf("IKUAI_CRON_CUSTOM_ISP_%s", id)
		value := os.Getenv(key)
		slice := strings.Split(value, "|")

		cron := ""
		name := ""
		url := ""
		comment := ""

		if len(slice) < 3 {
			continue
		}
		cron = slice[0]
		name = slice[1]
		url = slice[2]
		if len(slice) > 3 {
			comment = slice[3]
		}

		if _, exist := m[id]; !exist {
			m[id] = &IKuaiCronCustomISP{
				Cron:    cron,
				Name:    name,
				Url:     strings.Split(url, ","),
				Comment: comment,
			}
		}
	}

	for _, v := range m {
		c.IKuaiCronCustomISPList = append(c.IKuaiCronCustomISPList, v)
	}
}

func (c *Config) matchCronStreamDomain() {
	re := regexp.MustCompile(`IKUAI_CRON_STREAM_DOMAIN_(\d+)`)
	m := map[string]*IKuaiCronStreamDomain{}

	for _, env := range os.Environ() {
		match := re.FindStringSubmatch(env)
		if len(match) < 2 {
			continue
		}
		id := match[1]
		key := fmt.Sprintf("IKUAI_CRON_STREAM_DOMAIN_%s", id)
		value := os.Getenv(key)
		slice := strings.Split(value, "|")

		cron := ""
		interFace := ""
		url := ""
		srcAddr := ""
		comment := ""

		if len(slice) < 3 {
			continue
		}
		cron = slice[0]
		interFace = slice[1]
		url = slice[2]
		if len(slice) > 3 {
			srcAddr = slice[3]
		}
		if len(slice) > 4 {
			comment = slice[4]
		}

		if _, exist := m[id]; !exist {
			m[id] = &IKuaiCronStreamDomain{
				Cron:      cron,
				Interface: strings.Split(interFace, ","),
				Url:       strings.Split(url, ","),
				SrcAddr:   srcAddr,
				Comment:   comment,
			}
		}
	}

	for _, v := range m {
		c.IKuaiCronStreamDomainList = append(c.IKuaiCronStreamDomainList, v)
	}
}

func (c *Config) matchCronIpGroup() {
	re := regexp.MustCompile(`IKUAI_CRON_IP_GROUP_(\d+)`)
	m := map[string]*IKuaiCronIPGROUP{}
	for _, env := range os.Environ() {
		match := re.FindStringSubmatch(env)
		if len(match) < 2 {
			continue
		}
		id := match[1]
		key := fmt.Sprintf("IKUAI_CRON_IP_GROUP_%s", id)
		value := os.Getenv(key)
		slice := strings.Split(value, "|")
		cron := ""
		name := ""
		url := ""
		comment := ""

		if len(slice) < 3 {
			continue
		}
		cron = slice[0]
		name = slice[1]
		url = slice[2]
		if len(slice) > 3 {
			comment = slice[3]
		}
		if _, exist := m[id]; !exist {
			m[id] = &IKuaiCronIPGROUP{
				Cron:    cron,
				Name:    name,
				Url:     strings.Split(url, ","),
				Comment: comment,
			}
		}

	}
	for _, v := range m {
		c.IKuaiCronIPGroupList = append(c.IKuaiCronIPGroupList, v)
	}
}
