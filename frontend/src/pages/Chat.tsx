import { auth, chat } from "../../wailsjs/go/models";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import MainLayout from "@/components/MainLayout";
import ChatList from "@/components/ChatList";
import { useEffect, useState } from "react";
import { GetProfile } from "../../wailsjs/go/auth/AuthService";
import ChatRoom from "@/components/ChatRoom";

export default function Chat() {
  const [user, setUser] = useState<auth.UserModel>();
  const [room, setRoom] = useState<chat.ChatRoom>();

  const handleSelectRoom = (selected: chat.ChatRoom) => {
    setRoom(selected);
  };

  useEffect(() => {
    GetProfile()
      .then((res) => {
        setUser(res);
      })
      .catch((e) => {});
  }, []);

  return (
    <MainLayout className="flex">
      {user && (
        <>
          <ChatList user={user} onSelect={handleSelectRoom} />
          {room && <ChatRoom user={user} room={room} />}
        </>
      )}
    </MainLayout>
  );
}
