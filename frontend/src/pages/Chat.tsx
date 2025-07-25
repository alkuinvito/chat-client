import { chat } from "../../wailsjs/go/models";
import { LogInfo } from "../../wailsjs/runtime/runtime";
import MainLayout from "@/components/MainLayout";
import ChatRoom from "../components/ChatRoom";

export default function Chat() {
  const handleSelectRoom = (selected: chat.ChatRoom) => {
    LogInfo(selected.peer_name);
  };

  return (
    <MainLayout className="flex">
      <ChatRoom onSelect={handleSelectRoom} />
      <div></div>
    </MainLayout>
  );
}
