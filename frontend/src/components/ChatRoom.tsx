import type { TProfileSchema, TResponseSchema } from "@/models";
import type { chat, user } from "../../wailsjs/go/models";
import { EventsOn, LogInfo } from "../../wailsjs/runtime/runtime";
import { useEffect, useRef, useState } from "react";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { SendMessage } from "../../wailsjs/go/chat/ChatService";
import { useVirtualizer } from "@tanstack/react-virtual";

interface ChatRoomProps {
  user: TProfileSchema;
  contact: user.ContactModel;
}

interface ChatMessage {
  sender: string;
  message: string;
}

export default function ChatRoom({ user, contact }: ChatRoomProps) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [message, setMessage] = useState("");

  const parentRef = useRef<HTMLUListElement>(null);

  const virtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 48,
    enabled: true,
  });

  const items = virtualizer.getVirtualItems();

  const handleSubmit = (input: string) => {
    if (input.length >= 1 && input.length <= 250) {
      const message: chat.ChatMessage = {
        sender: user.id,
        message: input,
      };

      SendMessage(contact, message)
        .then((res: TResponseSchema<string>) => {
          if (res.code != 200) {
            LogInfo(`Error ${res.code.toString()} - ${res.data}`);
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
    return EventsOn("msg:new:" + contact.id, (msg: ChatMessage) => {
      setMessages((prev) => [...(prev ?? []), msg]);
    });
  }, [contact]);

  useEffect(() => {
    if (messages.length > 0) {
      virtualizer.scrollToIndex(messages.length - 1, { align: "end" });
    }
  }, [messages]);

  return (
    <div className="grow h-full flex flex-col">
      <ul
        ref={parentRef}
        className="grow overflow-y-auto p-2 pb-0 text-left"
        style={{ willChange: "transform" }}
      >
        <div
          style={{
            height: virtualizer.getTotalSize(),
            width: "100%",
            position: "relative",
          }}
        >
          <div
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              width: "100%",
              transform: `translateY(${items[0]?.start ?? 0}px)`,
            }}
          >
            {virtualizer.getVirtualItems().map((virtualRow) => {
              const message = messages[virtualRow.index];
              const isOwn = message.sender === user.id;

              return (
                <li
                  key={`${virtualRow.key}`}
                  data-index={virtualRow.index}
                  ref={virtualizer.measureElement}
                  className={`flex w-full pb-2 ${isOwn ? "justify-end" : "justify-start"}`}
                >
                  <span className="max-w-1/2 px-2 py-1 bg-neutral-800 rounded-lg whitespace-pre-wrap wrap-break-word">
                    {message.message}
                  </span>
                </li>
              );
            })}
          </div>
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
