import type { chat, discovery } from "../../wailsjs/go/models";
import MainLayout from "@/components/MainLayout";
import PeerList from "@/components/PeerList";
import { useEffect, useState } from "react";
import { GetProfile, RequestPair } from "../../wailsjs/go/user/UserService";
import ChatRoom from "@/components/ChatRoom";
import type { TProfileSchema, TResponseSchema } from "@/models";
import { useNavigate } from "react-router";

export default function Chat() {
  const [user, setUser] = useState<TProfileSchema>();
  const [peer, setPeer] = useState<discovery.PeerModel>();

  const navigate = useNavigate();

  const handleSelect = (selected: discovery.PeerModel) => {
    RequestPair(selected)
      .then((res: TResponseSchema<string>) => {
        switch (res.code) {
          case 200:
            alert(res.data);
            break;
          default:
            alert(res.data);
            break;
        }
      })
      .catch((e) => {});
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
          <PeerList user={user} onSelect={handleSelect} />
          {peer && <ChatRoom user={user} peer={peer} />}
        </>
      )}
    </MainLayout>
  );
}
