import type { TProfileSchema, TResponseSchema } from "@/models";
import type { chat, user } from "../../wailsjs/go/models";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { useEffect, useRef, useState } from "react";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";
import { GetMessages, SendMessage } from "../../wailsjs/go/chat/ChatService";
import { useVirtualizer } from "@tanstack/react-virtual";
import { toast } from "sonner";
import { ArrowDown, Info } from "lucide-react";

interface ChatRoomProps {
  user: TProfileSchema;
  contact: user.ContactModel;
}

export default function ChatRoom({ user, contact }: ChatRoomProps) {
  const [messages, setMessages] = useState<chat.ChatMessage[]>([]);
  const [message, setMessage] = useState("");
  const [hasMore, setHasMore] = useState(true);
  const [autoScroll, setAutoScroll] = useState(true);
  const [cursor, setCursor] = useState(0);

  const parentRef = useRef<HTMLUListElement>(null);

  const virtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 48,
    enabled: true,
  });

  const items = virtualizer.getVirtualItems();

  const getMessages = (peerId: string, cursor: number) => {
    GetMessages(peerId, cursor)
      .then((res: TResponseSchema<chat.ChatMessage[]>) => {
        switch (res.code) {
          case 200:
            if (res.data.length === 0) {
              setHasMore(false);
              break;
            }
            setHasMore(true);
            setMessages((prev) => [...res.data, ...(prev ?? [])]);
            break;
          case 404:
            setHasMore(false);
            toast.info("No older messages", { icon: <Info /> });
            break;
          default:
            toast.error("Error retrieving older messages", { icon: <Info /> });
            break;
        }
      })
      .catch(() => {});
  };

  const handleSubmit = (input: string) => {
    if (input.length >= 1 && input.length <= 250) {
      setMessage("");

      const message = {
        sender: user.id,
        message: input,
      };

      SendMessage(contact, message)
        .then((res: TResponseSchema<chat.ChatMessage>) => {
          switch (res.code) {
            case 200:
              setMessages((prev) => [...(prev ?? []), res.data]);
              break;
            case 404:
              toast.error("Peer is appear to be offline", { icon: <Info /> });
              break;
            default:
              toast.error("Failed to send message", { icon: <Info /> });
              break;
          }
        })
        .catch(() => {});
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e.currentTarget.value);
    }
  };

  const scrollToBottom = () => {
    virtualizer.scrollToIndex(messages.length - 1, { align: "end" });
  };

  useEffect(() => {
    const elem = parentRef.current;
    if (!elem) return;

    const handleScroll = () => {
      const distanceFromBottom =
        elem.scrollHeight - elem.scrollTop - elem.clientHeight;

      setAutoScroll(distanceFromBottom < 50);
    };

    elem.addEventListener("scroll", handleScroll);

    return () => elem.removeEventListener("scroll", handleScroll);
  }, []);

  useEffect(() => {
    getMessages(contact.id, 0);

    // listen for new message
    const unsubscribeMsg = EventsOn("msg:new", (msg: chat.ChatMessage) => {
      if (msg.peer_id === contact.id) {
        setMessages((prev) => [...(prev ?? []), msg]);
      }
    });

    return () => {
      unsubscribeMsg();
    };
  }, [contact]);

  useEffect(() => {
    // check if user manually scrolled up
    if (!autoScroll) {
      // check if user on oldest message
      if (items[0]?.index === 0 && hasMore) {
        // make sure the getMessages called once
        if (messages[0].id != cursor) {
          getMessages(contact.id, messages[0].id || 0);
          setCursor(messages[0].id);
        }
      }
    }
  }, [items, autoScroll, hasMore, contact.id, messages]);

  useEffect(() => {
    if (autoScroll) {
      if (messages.length > 0) {
        scrollToBottom();
      }
    }
  }, [messages]);

  return (
    <div className="grow h-full flex flex-col relative">
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
                  <span className="max-w-1/2 px-2 py-1 bg-neutral-800 rounded-lg whitespace-pre-wrap wrap-break-word !select-text">
                    {message.message}
                  </span>
                </li>
              );
            })}
          </div>
        </div>
      </ul>
      <div
        className={
          "absolute bottom-16 right-1/2 translate-x-1/2 " +
          (autoScroll ? "hidden opacity-0" : "block opacity-100")
        }
      >
        <Button
          variant="outline"
          onClick={scrollToBottom}
          className="bg-neutral-900 rounded-full"
        >
          <ArrowDown />
          Scroll to bottom
        </Button>
      </div>
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
