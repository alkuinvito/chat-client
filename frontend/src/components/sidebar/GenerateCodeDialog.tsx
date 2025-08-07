import { useEffect, useState } from "react";
import { GeneratePairingCode } from "../../../wailsjs/go/user/UserService";
import { toast } from "sonner";
import { Info, KeyRound } from "lucide-react";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Progress } from "../ui/progress";

export default function GenerateCodeDialog() {
  const [code, setCode] = useState("------");
  const [value, setValue] = useState(0);

  const handleGenerate = () => {
    setCode("------");

    GeneratePairingCode()
      .then((res) => {
        if (res.code === 200) {
          setCode(res.data);
          setValue(100);
        } else {
          toast.error(res.data, { icon: <Info /> });
        }
      })
      .catch(() => {});
  };

  useEffect(() => {
    const duration = 60 * 1000;
    const step = 100;
    const totalSteps = duration / step;
    const decrement = 100 / totalSteps;

    const interval = setInterval(() => {
      setValue((prev) => {
        if (prev <= 0) {
          clearInterval(interval);
          return 0;
        }
        return prev - decrement;
      });
    }, step);

    if (value < 2) {
      setCode("------");
    }

    return () => clearInterval(interval);
  }, [value]);

  return (
    <Dialog>
      <DialogTrigger className="w-full">
        <Button variant="outline" className="w-full">
          <KeyRound />
          Get code
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Generate pairing code</DialogTitle>
        </DialogHeader>
        <div>
          <div className="flex gap-1 items-center mb-3 p-2 bg-neutral-900 border border-neutral-800 text-sm text-neutral-400 rounded-md">
            <Info size={16} />
            <span>Pairing code will expires in 60 seconds</span>
          </div>
          <div className="p-3 border border-neutral-800 rounded-md">
            <span className="text-3xl font-bold tracking-widest select-text">
              {code}
            </span>
            <Progress className="w-full mt-1" value={value} />
          </div>
        </div>
        <DialogFooter>
          <Button
            onClick={() => {
              handleGenerate();
            }}
          >
            <KeyRound />
            Generate
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
