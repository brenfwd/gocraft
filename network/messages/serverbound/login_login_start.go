package messages

import (
	"log"

	"github.com/brenfwd/gocraft/constants"
	"github.com/brenfwd/gocraft/data"
	"github.com/brenfwd/gocraft/ipc"
	"github.com/brenfwd/gocraft/network/messages"
	"github.com/brenfwd/gocraft/network/messages/clientbound"
	"github.com/google/uuid"
)

func init() {
	messages.RegisterServerbound[LoginServerboundLoginStart](constants.ClientStateLogin, 0x00)
}

type LoginServerboundLoginStart struct {
	messages.Serverbound
	Name       string
	PlayerUUID uuid.UUID
}

func (p *LoginServerboundLoginStart) Handle(i *ipc.ClientIPC) error {
	log.Println(p)
	reason := data.MakeChat().SetText("")
	reason.AddExtra(
		data.MakeChat().SetText("Your account has been restricted\n\n").SetColor(data.ChatColorWhite).SetBold(true),
		data.MakeChat().SetText("Due to a moderation action taken against your account, you will not be able to play on Multiplayer servers for 4371 day(s).\n\n\n"),
		data.MakeChat().SetText("The reason provided for the restriction was:\nSpreading harmful or malicious misinformation with the intent to disrupt a democratic proceeding\n\n\n"),
		data.MakeChat().SetText("For more information, visit https://aka.ms/mcaccres\n\n\n"),
	)
	res := clientbound.LoginClientboundDisconnect{
		Reason: *reason,
	}
	encoded, err := messages.Encode(&res)
	if err != nil {
		return err
	}
	i.SendPacket(&encoded)
	return nil
}
