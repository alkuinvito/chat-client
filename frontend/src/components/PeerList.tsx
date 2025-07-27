import type { discovery } from "../../wailsjs/go/models";
import { useState } from "react";
import type { TProfileSchema, TResponseSchema } from "@/models";
import { GetPeers } from "../../wailsjs/go/discovery/DiscoveryService";
import { Info, RotateCw } from "lucide-react";
import { toast } from "sonner";

interface PeerListProps {
  user: TProfileSchema;
  onSelect: (peer: discovery.PeerModel) => void;
}

export default function PeerList({ user, onSelect }: PeerListProps) {
  const [rooms, setRooms] = useState<discovery.PeerModel[]>();
  const [current, setCurrent] = useState<discovery.PeerModel>();
  const [isLoading, setIsLoading] = useState(false);

  const getRooms = () => {
    setIsLoading(true);

    GetPeers()
      .then((res: TResponseSchema<discovery.PeerModel[]>) => {
        switch (res.code) {
          case 200:
            setRooms(res.data);
            break;
          default:
            toast.error("Error when searching peers", { icon: <Info /> });
            break;
        }
      })
      .catch((e) => {})
      .finally(() => {
        setIsLoading(false);
      });
  };

  return (
    <div className="flex flex-col w-screen max-w-[280px] h-full bg-neutral-800">
      <div className="flex items-center w-full h-12 border-b border-neutral-900">
        <div className="grow flex items-center px-4 h-full text-left border-r border-neutral-900">
          <span className="select-none">{user.username}</span>
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
        <ul className="grow h-full overflow-y-auto">
          {rooms.map((room) => (
            <li key={room.ip}>
              <button
                className="grid px-4 py-2 w-full text-left enabled:hover:bg-neutral-700 transition-colors disabled:bg-neutral-700 disabled:hover:bg-neutral-600"
                onClick={() => {
                  setCurrent(room);
                  onSelect(room);
                }}
                disabled={room == current}
              >
                <span className="select-none line-clamp-1">
                  {room.username}
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
