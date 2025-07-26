import { auth, chat } from "../../wailsjs/go/models";
import MainLayout from "@/components/MainLayout";
import ChatList from "@/components/ChatList";
import { useEffect, useState } from "react";
import { GetProfile } from "../../wailsjs/go/auth/AuthService";
import ChatRoom from "@/components/ChatRoom";
import type { TProfileSchema, TResponseSchema } from "@/models";
import { LogInfo } from "../../wailsjs/runtime/runtime";
import { useNavigate } from "react-router";

export default function Chat() {
  const [user, setUser] = useState<TProfileSchema>();
  const [room, setRoom] = useState<chat.ChatRoom>();

  const navigate = useNavigate();

  const handleSelectRoom = (selected: chat.ChatRoom) => {
    setRoom(selected);
  };

  useEffect(() => {
    GetProfile()
      .then((res: TResponseSchema<TProfileSchema>) => {
        switch (res.code) {
          case 200:
            setUser(res.data);
            break;
          case 404:
            navigate("/");
            break;
        }
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
