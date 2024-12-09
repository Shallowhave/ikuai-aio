package job

import (
	"github.com/NERVEbing/ikuai-aio/api"
	"github.com/NERVEbing/ikuai-aio/config"
	"time"
)

func updateIPGroup(c *config.IKuaiCronIPGROUP, tag string) error {
	var rows []string
	start := time.Now()
	for _, url := range c.Url {
		r, err := fetch(url)
		if err != nil {
			logger(tag, "fetch %s failed, error: %s", url, err)
			continue
		}
		logger(tag, "fetch %s success, rows: %d", url, len(r))
		rows = append(rows, r...)
	}
	logger(tag, "fetch total rows: %d", len(rows))
	if len(rows) == 0 {
		return nil
	}

	client := api.NewClient()
	if err := client.Login(); err != nil {
		return err
	}
	ipGroupShowResp, err := client.IPGroupShow()
	if err != nil {
		return err
	}
	var ids []int
	for _, i := range ipGroupShowResp.Data.Data {
		if i.GroupName == c.Name {
			ids = append(ids, i.ID)
		}
	}
	if err = client.IpGroupDel(ids); err != nil {
		return err
	}
	count, err := client.IpGroupAdd(c.Name, rows, c.Comment)
	if err != nil {
		return err
	}
	logger(tag, "add custom isp unique rows count: %d, duration: %s", count, time.Now().Sub(start).String())

	return nil
}
