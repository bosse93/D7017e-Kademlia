package main

type Network struct {
	rt RoutingTable
}

func NewNetwork(rt *RoutingTable) *Network {
	Network.rt = *rt
	return Network
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

	Network.rt.AddContact(packet.sender)

	switch packet.messageType {
		case "SendPingMessage":
			//Måste svara nåt! Svara med sig själv?

		case "SendFindContactMessage":
			answerData := Network.rt.FindClosestContacts(packet.target, 20)

			//Pack answer. answerData, packet.messageID, ....

		case "SendFindDataMessage":
			//Check if data in storage, packet.target(ID to data)

			//If Data 
			//send file How?

			//Else
			answerData := Network.rt.FindClosestContacts(packet.target, 20)

		case "SendStoreMessage":
			//Store data
			//Kolla om data redan finns?

			//Svara med något?


		default:
			//Not a valid message


	}

	//Send marshalled response on UDP



}