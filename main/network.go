package main

type Network struct {
	rt RoutingTable
}

func NewNetwork(rt *RoutingTable) *Network {
	network := &Network{}
	network.rt = *rt
	return network
}

func Listen(ip string, port int) {
	// TODO

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}






func (network *Network) handleRequest() {
	//Unpack packet
	//packet := ....

	network.rt.AddContact(packet.sender)

	switch packet.messageType {
		case "SendPingMessage":
			//Måste svara nåt! Svara med sig själv?

		case "SendFindContactMessage":
			answerData := network.rt.FindClosestContacts(packet.target, 20)

			//Pack answer. answerData, packet.messageID, ....

		case "SendFindDataMessage":
			//Check if data in storage, packet.target(ID to data)

			//If Data 
			//send file How?

			//Else
			answerData := network.rt.FindClosestContacts(packet.target, 20)

		case "SendStoreMessage":
			//Store data
			//Kolla om data redan finns?

			//Svara med något?


		default:
			//Not a valid message


	}

	//Send marshalled response on UDP



}