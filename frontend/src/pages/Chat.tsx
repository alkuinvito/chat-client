import type { user } from "../../wailsjs/go/models";
import MainLayout from "@/components/MainLayout";
import Sidebar from "@/components/sidebar/Sidebar";
import { useEffect, useState } from "react";
import { GetProfile } from "../../wailsjs/go/user/UserService";
import ChatRoom from "@/components/ChatRoom";
import type { TProfileSchema, TResponseSchema } from "@/models";
import { useNavigate } from "react-router";

export default function Chat() {
  const [user, setUser] = useState<TProfileSchema>();
  const [contact, setContact] = useState<user.ContactModel>();

  const navigate = useNavigate();

  const handleSelect = (selected: user.ContactModel) => {
    setContact(selected);
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
      .catch(() => {});
  }, []);

  return (
    <MainLayout className="flex">
      {user && (
        <>
          <Sidebar user={user} onSelect={handleSelect} />
          {contact && <ChatRoom user={user} contact={contact} />}
        </>
      )}
    </MainLayout>
  );
}
