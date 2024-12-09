package api

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

func (c *Client) IPGroupShow() (*IPGroupShowResp, error) {
	req := &CallReq{
		FuncName: "ipgroup",
		Action:   "show",
		Param: map[string]string{
			"TYPE": "data",
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.request(iKuaiCallPath, b)
	if err != nil {
		return nil, err
	}

	var mod IPGroupShowResp
	if err = json.Unmarshal(resp, &mod); err != nil {
		return nil, err
	}
	if mod.Result != 30000 {
		return nil, errors.New(mod.ErrMsg)
	}

	return &mod, nil
}

func (c *Client) IpGroupDel(ids []int) error {
	id := ""
	if len(ids) == 0 {
		return nil
	}
	var idStr []string
	for _, i := range ids {
		idStr = append(idStr, strconv.Itoa(i))
	}
	id = strings.Join(idStr, ",")
	req := &CallReq{
		FuncName: "ipgroup",
		Action:   "del",
		Param: map[string]string{
			"id": id,
		},
	}
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := c.request(iKuaiCallPath, b)
	if err != nil {
		return err
	}

	var mod IPGroupDelResp
	if err = json.Unmarshal(resp, &mod); err != nil {
		return err
	}
	if mod.Result != 30000 {
		return errors.New(mod.ErrMsg)
	}

	return nil
}

func (c *Client) IpGroupAdd(groupName string, addrPool []string, comment string) (int, error) {
	m := map[string]bool{}
	if len(comment) == 0 {
		comment = "ikuai-aio"
	}
	for _, i := range addrPool {
		i = parseIPv4(i)
		if len(i) == 0 {
			continue
		}
		if _, exist := m[i]; !exist {
			m[i] = false
		}
		comment = comment + ",%20" + comment

	}

	addrPool = make([]string, 0, len(m))
	for row := range m {
		addrPool = append(addrPool, row)
	}

	chunkSize := 5000
	ipGroupSlices := chunkSliceStr(addrPool, chunkSize)
	for _, slice := range ipGroupSlices {
		req := &CallReq{
			FuncName: "ipgroup",
			Action:   "add",
			Param: map[string]string{
				"group_name": groupName,
				"addr_pool":  strings.Join(slice, ","),
				"comment":    comment,
				"type":       "1",
				"NewRow":     "true",
			},
		}

		b, err := json.Marshal(req)
		if err != nil {
			return 0, err
		}
		resp, err := c.request(iKuaiCallPath, b)
		if err != nil {
			return 0, err
		}

		var mod IPGroupAddResp
		if err = json.Unmarshal(resp, &mod); err != nil {
			return 0, err
		}
		if mod.Result != 30000 {
			return 0, errors.New(mod.ErrMsg)
		}
	}

	return len(addrPool), nil
}
