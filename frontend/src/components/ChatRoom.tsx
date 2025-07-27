import type { TProfileSchema, TResponseSchema } from "@/models";
import type { chat, discovery } from "../../wailsjs/go/models";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import { useEffect, useRef, useState } from "react";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { SendMessage } from "../../wailsjs/go/chat/ChatService";
import { useVirtualizer } from "@tanstack/react-virtual";

interface ChatRoomProps {
  user: TProfileSchema;
  peer: discovery.PeerModel;
}

interface ChatMessage {
  sender: string;
  message: string;
}

export default function ChatRoom({ user, peer }: ChatRoomProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [message, setMessage] = useState("");

  const parentRef = useRef<HTMLUListElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 48,
    enabled: true,
  });

  const handleSubmit = (input: string) => {
    if (input.length >= 1 && input.length <= 250) {
      const message: chat.ChatMessage = {
        sender: user.username,
        message: input,
      };

      SendMessage(peer, message)
        .then((res: TResponseSchema<string>) => {
          if (res.code != 200) {
            LogInfo(`Error ${res.code.toString()} - ${res.code}`);
          } else {
            setMessages((prev) => [...(prev ?? []), message]);
          }
        })
        .catch(() => {})
        .finally(() => {
          setMessage("");
        });
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e.currentTarget.value);
    }
  };

  useEffect(() => {
    EventsOn("msg:new", (ev) => {
      const msg = JSON.parse(ev) as ChatMessage;
      setMessages((prev) => [...(prev ?? []), msg]);
    });
  }, []);

  useEffect(() => {
    if (messages.length > 0) {
      rowVirtualizer.scrollToIndex(messages.length - 1, { align: "end" });
    }
  }, [messages]);

  return (
    <div className="grow h-full flex flex-col">
      <ul
        ref={parentRef}
        className="grow overflow-y-auto p-2 text-left"
        style={{ willChange: "transform" }}
      >
        <div
          style={{
            height: `${rowVirtualizer.getTotalSize()}px`,
            position: "relative",
          }}
        >
          {rowVirtualizer.getVirtualItems().map((virtualRow) => {
            const message = messages[virtualRow.index];
            const isOwn = message.sender === user.username;

            return (
              <li
                key={`msg-${virtualRow.key}`}
                className={`absolute w-full ${isOwn ? "flex justify-end" : ""}`}
                style={{
                  top: 0,
                  transform: `translateY(${virtualRow.start}px)`,
                }}
              >
                <span className="max-w-1/2 px-2 py-1 bg-neutral-800 rounded-lg whitespace-pre-wrap">
                  {message.message}
                </span>
              </li>
            );
          })}
        </div>
      </ul>
      <div className="flex gap-2 p-2">
        <Textarea
          name="message"
          className="border-neutral-700 h-auto min-h-9 text-base resize-none"
          minLength={1}
          maxLength={250}
          rows={1}
          value={message}
          onInput={(e) => {
            setMessage(e.currentTarget.value);
          }}
          onKeyDown={handleKeyDown}
        />
        <Button
          variant="secondary"
          onClick={() => handleSubmit(message)}
          disabled={message.length < 1 || message.length > 250}
        >
          Send
        </Button>
      </div>
    </div>
  );
}
