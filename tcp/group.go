package tcp

import 	log "github.com/sirupsen/logrus"


func (c *TcpClients) append(node *TcpClientNode) {
	*c = append(*c, node)
}

func (c *TcpClients) send(msgId int64, data []byte) {
	for _, node := range *c {
		node.Send(msgId, data)
	}
}

func (c *TcpClients) asyncSend(data []byte) {
	for _, node := range *c {
		node.AsyncSend(data)
	}
}

func (c *TcpClients) remove(node *TcpClientNode) {
	for index, n := range *c {
		if n == node {
			*c = append((*c)[:index], (*c)[index+1:]...)
			break
		}
	}
	log.Debugf("#####################remove node, current len %v", len(*c))
}

func (c *TcpClients) close() {
	for _, node := range *c {
		node.close()
	}
}