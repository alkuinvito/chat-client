import { useState, useEffect } from "react";
import {
  RequestPairing,
  ScanPeers,
} from "../../../wailsjs/go/user/UserService";
import type { discovery, user } from "wailsjs/go/models";
import { toast } from "sonner";
import { Info, RadioTower, Send } from "lucide-react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Loader } from "../Loader";
import { InputOTP, InputOTPGroup, InputOTPSlot } from "../ui/input-otp";

export default function PairDialog() {
  const [currPeer, setCurrPeer] = useState<discovery.PeerModel>();
  const [code, setCode] = useState("");
  const [peers, setPeers] = useState<discovery.PeerModel[]>([]);
  const [isLoading, setLoading] = useState(false);

  const handleSubmit = () => {
    if (!currPeer || code.length != 6) return;

    const req: user.RequestPairSchema = {
      id: currPeer.id,
      username: currPeer.username,
      code: code,
    };

    RequestPairing(req)
      .then((res) => {
        if (res.code === 200) {
          toast.success("Paired successfully", { icon: <Info /> });
        } else {
          toast.error(res.data, { icon: <Info /> });
        }
      })
      .catch(() => {});
  };

  const scanPeers = () => {
    setLoading(true);
    setCurrPeer(undefined);

    ScanPeers()
      .then((res) => {
        if (res.code === 200) {
          setPeers(res.data);
        } else {
          toast.error("Error fetching pair requests", { icon: <Info /> });
        }
      })
      .catch(() => {})
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    scanPeers();
  }, []);

  return (
    <Dialog>
      <DialogTrigger className="w-full">
        <Button variant="outline" className="w-full">
          <RadioTower />
          Pair
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Send request to pair</DialogTitle>
        </DialogHeader>
        <div className="text-left grid gap-2">
          <div>
            <span>
              1. Ask peer to <b>generate</b> 6-digit pairing code
            </span>
          </div>

          <div>
            <span>
              2. Select available peer with matching <b>name</b> and <b>ID</b>
            </span>
            <ul className="h-64 overflow-y-auto border border-neutral-800 rounded-sm">
              {isLoading ? (
                <div className="h-full w-full flex justify-center items-center">
                  <Loader className="size-8" />
                </div>
              ) : peers ? (
                peers.map((peer) => (
                  <li key={peer.id} className="w-full">
                    <button
                      className="grid w-full px-2 py-1 enabled:hover:bg-neutral-800 disabled:bg-neutral-800  text-left"
                      onClick={() => {
                        setCurrPeer(peer);
                      }}
                      disabled={peer.id === currPeer?.id}
                    >
                      <div className="flex justify-between">
                        <span>{peer.username}</span>
                        <span className="text-sm text-neutral-400">
                          {peer.ip}
                        </span>
                      </div>
                      <span className="text-xs text-neutral-400">
                        {peer.id}
                      </span>
                    </button>
                  </li>
                ))
              ) : (
                <div className="h-full w-full flex justify-center items-center">
                  <span className="text-neutral-400">No peers found.</span>
                </div>
              )}
            </ul>
          </div>

          <div>
            <span>
              3. Input <b>pairing code</b> generated from peer
            </span>
            <div className="flex justify-center items-center w-full">
              <InputOTP
                maxLength={6}
                onChange={(e) => {
                  setCode(e);
                }}
              >
                <InputOTPGroup>
                  <InputOTPSlot index={0} />
                  <InputOTPSlot index={1} />
                  <InputOTPSlot index={2} />
                  <InputOTPSlot index={3} />
                  <InputOTPSlot index={4} />
                  <InputOTPSlot index={5} />
                </InputOTPGroup>
              </InputOTP>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button
            variant="outline"
            className="mr-2"
            onClick={() => {
              scanPeers();
            }}
          >
            Refresh
          </Button>
          <Button
            onClick={() => {
              handleSubmit();
            }}
            disabled={!currPeer || code.length != 6}
          >
            <Send />
            Pair
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
