import { Loader2 } from "lucide-react";
import { twMerge } from "tailwind-merge";

interface LoaderProps {
  className?: string;
}

export function Loader({ className }: LoaderProps) {
  return <Loader2 className={twMerge("h-5 w-5 animate-spin", className)} />;
}
