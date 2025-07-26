import { auth, chat } from "../../wailsjs/go/models";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import { useEffect, useState } from "react";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { SendMessage } from "../../wailsjs/go/chat/ChatService";

interface ChatRoomProps {
  user: auth.UserModel;
  room: chat.ChatRoom;
}

interface ChatMessage {
  sender: string;
  message: string;
}

export default function ChatRoom({ user, room }: ChatRoomProps) {
  const [messages, setMessage] = useState<ChatMessage[]>([]);
  const [msg, setMsg] = useState("");

  const handleSubmit = (input: string) => {
    if (input.length >= 1 && input.length <= 250) {
      const message: chat.ChatMessage = {
        sender: user.username,
        message: input,
      };

      SendMessage(room, message)
        .then((res) => {
          if (res.code != 200) {
            LogInfo(`Error ${res.code.toString()} - ${res.code}`);
          } else {
            setMessage((prev) => [...(prev ?? []), message]);
          }
        })
        .catch(() => {})
        .finally(() => {
          setMsg("");
        });
    }
  };

  useEffect(() => {
    EventsOn("msg:new", (ev) => {
      const message = JSON.parse(ev) as ChatMessage;
      setMessage((prev) => [...(prev ?? []), message]);
    });
  }, []);

  return (
    <div className="grow h-full flex flex-col">
      <ul className="grow overflow-y-auto p-2">
        {messages &&
          messages.map((message, i) =>
            message.sender == user.username ? (
              <li
                key={`msg-${i}`}
                className="flex flex-col items-end w-full my-2 text-right"
              >
                <span className="text-xs text-neutral-300">
                  {message.sender}
                </span>
                <span className="w-min max-w-1/2 p-2 bg-neutral-800 rounded-lg">
                  {message.message}
                </span>
              </li>
            ) : (
              <li
                key={`msg-${i}`}
                className="flex flex-col w-full my-2 text-left"
              >
                <span className="text-xs text-neutral-300">
                  {message.sender}
                </span>
                <span className="w-min max-w-1/2 p-2 bg-neutral-800 rounded-lg">
                  {message.message}
                </span>
              </li>
            ),
          )}
      </ul>
      <div className="flex gap-2 p-2">
        <Textarea
          name="message"
          className="border-neutral-700 h-auto min-h-9 text-base resize-none"
          minLength={1}
          maxLength={250}
          rows={1}
          value={msg}
          onInput={(e) => {
            setMsg(e.currentTarget.value);
          }}
        />
        <Button variant="secondary" onClick={() => handleSubmit(msg)}>
          Send
        </Button>
      </div>
    </div>
  );
}
