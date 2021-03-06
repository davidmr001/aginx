package dockerLabels

import (
	"encoding/json"
	"regexp"
	"strconv"
)

var keyRegexp = regexp.MustCompile("aginx.domain(\\.(\\d+))?")
var valueRegexp = regexp.MustCompile("([a-zA-Z0-9-_\\.]*)(,(weight=(\\d+)))?(,(internal))?(,(ssl))?(,(virtual))?(,(nodes))?")

type Label struct {
	Domain   string
	Port     int
	Weight   int  //服务器权重。在费Swarm节点下起作用
	AutoSSL  bool //自动生成证书文件
	Internal bool //使用内部地址
	Virtual  bool //虚拟VIP 只配置一个
	Nodes    bool //外部接口全节点配置，这里还可以筛选
}

type Labels map[int]Label

func (ls *Labels) Has() bool {
	return len(*ls) > 0
}

func (ls *Labels) String() string {
	bs, _ := json.Marshal(ls)
	return string(bs)
}

func FindLabels(labs map[string]string, ignoreSwarmService bool) Labels {
	labels := Labels{}
	if _, has := labs["com.docker.swarm.task.id"]; ignoreSwarmService && has {
		return labels
	}
	for key, value := range labs {
		if keyRegexp.MatchString(key) && valueRegexp.MatchString(value) {
			domain := valueRegexp.FindStringSubmatch(value)
			port := keyRegexp.FindStringSubmatch(key)
			label := Label{Domain: domain[1]}
			label.Weight, _ = strconv.Atoi(domain[4])
			label.Internal = domain[6] == "internal"
			label.Port, _ = strconv.Atoi(port[2])
			label.AutoSSL = domain[8] == "ssl"
			label.Virtual = domain[10] == "virtual"
			label.Nodes = domain[12] == "nodes"
			labels[label.Port] = label
		}
	}
	return labels
}
