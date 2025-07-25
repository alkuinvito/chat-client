import { chat } from "../../wailsjs/go/models";
import { useState } from "react";
import { GetRooms } from "../../wailsjs/go/chat/ChatService";
import { RotateCw } from "lucide-react";

interface ChatRoomProps {
  onSelect: (room: chat.ChatRoom) => void;
}

export default function ChatRoom({ onSelect }: ChatRoomProps) {
  const [rooms, setRooms] = useState<chat.ChatRoom[]>();
  const [isLoading, setIsLoading] = useState(false);

  const getRooms = () => {
    setIsLoading(true);

    GetRooms()
      .then((res) => {
        setRooms(res);
        setIsLoading(false);
      })
      .catch((e) => {
        setIsLoading(false);
      });
  };

  return (
    <div className="grow flex flex-col max-w-[280px] h-full bg-neutral-800">
      <div className="flex items-center w-full h-12 border-b border-neutral-900">
        <div className="grow flex items-center px-4 h-full text-left border-r border-neutral-900">
          <span className="select-none">Jajang</span>
        </div>
        <button
          className="size-12 hover:bg-neutral-700 transition-colors"
          onClick={getRooms}
        >
          <RotateCw className="mx-auto" size={20} />
        </button>
      </div>
      {isLoading ? (
        <div className="grow h-full flex items-center justify-center">
          <span className="text-neutral-500 select-none">
            Searching available peer(s)...
          </span>
        </div>
      ) : rooms ? (
        <ul className="grow h-full">
          {rooms.map((room) => (
            <li key={room.ip}>
              <button
                className="grid p-2 w-full text-left hover:bg-neutral-700 transition-colors"
                onClick={() => onSelect(room)}
              >
                <span className="select-none line-clamp-1">
                  {room.peer_name}
                </span>
                <span className="select-none line-clamp-1 text-xs text-neutral-400">
                  {room.ip}
                </span>
              </button>
            </li>
          ))}
        </ul>
      ) : (
        <div className="grow h-full flex items-center justify-center">
          <span className="text-neutral-500 select-none">No peer found.</span>
        </div>
      )}
    </div>
  );
}
