package job

import (
	"github.com/NERVEbing/ikuai-aio/api"
	"github.com/NERVEbing/ikuai-aio/config"
)

func updateCustomISP(c *config.IKuaiCronCustomISP) error {
	var rows []string
	for _, url := range c.Url {
		r, err := fetch(url)
		if err != nil {
			logger("updateCustomISP", "fetch %s error: %s", url, err)
			continue
		}
		logger("updateCustomISP", "fetch %s success, rows: %d", url, len(r))
		rows = append(rows, r...)
	}
	logger("updateCustomISP", "fetch total: %d", len(rows))
	if len(rows) == 0 {
		return nil
	}

	client := api.NewClient()
	if err := client.Login(); err != nil {
		return err
	}
	customISPShowResp, err := client.CustomISPShow()
	if err != nil {
		return err
	}
	var ids []int
	for _, i := range customISPShowResp.Data.Data {
		ids = append(ids, i.ID)
	}
	if err = client.CustomISPDel(ids); err != nil {
		return err
	}
	count, err := client.CustomISPAdd(c.Name, rows, c.Comment)
	if err != nil {
		return err
	}
	logger("updateCustomISP", "update custom isp success, count: %d", count)

	return nil
}
