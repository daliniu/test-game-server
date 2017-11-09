package service

import (
	"share/config"
	"share/tencent/cmq"
)

// InitQueues 初始化队列
func InitQueues(cmqClient *cmq.Client, queueNames ...string) {
	for _, queueName := range queueNames {
		switch queueName {
		case config.QueueIMRegister:
			cmqClient.CreateQueue(queueName, 1000000, 3, 30, 1024, 172800, 0)
			//case config.TopicGameOver:
			//	cmqClient.CreateTopic(config.TopicGameOver, 65536, 2, config.QueueGameStatsProcess, config.QueueCreditGameOver)
		case config.QueueMallBuyGoodsSuccess, config.QueueMallDeliverGift:
			cmqClient.CreateQueue(queueName, 1000000, 3, 30, 1024, 259200, 0)
		default:
			cmqClient.CreateQueue(queueName, 1000000, 3, 30, 1024, 86400, 0)
		}
	}
}
