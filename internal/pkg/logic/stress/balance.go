package stress

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex  int
	rss       []*WeightNode
}

type WeightNode struct {
	weight          int64  // 配置的权重，即在配置文件或初始化时约定好的每个节点的权重
	currentWeight   int64  //节点当前权重，会一直变化
	effectiveWeight int64  //有效权重，初始值为weight, 通讯过程中发现节点异常，则-1 ，之后再次选取本节点，调用成功一次则+1，直达恢复到weight 。 用于健康检查，处理异常节点，降低其权重。
	addr            string // 服务器addr
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("WeightRoundRobinBalance Add params len need 2")
	}

	// 获取值
	addr := params[0]
	weight, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	r.rss = append(r.rss, &WeightNode{
		weight:          weight,
		effectiveWeight: weight, // 初始化時有效权重 = 配置权重值
		currentWeight:   weight, // 初始化時当前权重 = 配置权重值
		addr:            addr,
	})

	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	// 没有服务
	if len(r.rss) == 0 {
		return ""
	}

	var totalWeight int64
	var maxWeightNode *WeightNode
	for key, node := range r.rss {
		// 计算当前状态下所有节点的effectiveWeight之和totalWeight
		totalWeight += node.effectiveWeight

		// 计算currentWeight
		node.currentWeight += node.effectiveWeight

		// 寻找权重最大的
		if maxWeightNode == nil || maxWeightNode.currentWeight < node.currentWeight {
			maxWeightNode = node
			r.curIndex = key
		}
	}

	// 更新选中节点的currentWeight
	maxWeightNode.currentWeight -= totalWeight

	// 返回addr
	return maxWeightNode.addr
}
