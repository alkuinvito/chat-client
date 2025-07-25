import { GetRooms } from "../../wailsjs/go/chat/ChatService";
import { LogInfo } from "../../wailsjs/runtime/runtime";

export default function Chat() {
  const getRooms = () => {
    GetRooms().then((rooms) => {
      for (const room of rooms) {
        LogInfo(room);
      }
    });
  };

  return (
    <div>
      <h1>Chat here</h1>
      <button onClick={getRooms}>Refresh</button>
    </div>
  );
}
