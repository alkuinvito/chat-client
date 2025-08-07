import { Loader2 } from "lucide-react";
import { twMerge } from "tailwind-merge";

interface LoaderProps {
  className?: string;
}

export function Loader({ className }: LoaderProps) {
  return <Loader2 className={twMerge("size-6 animate-spin", className)} />;
}
