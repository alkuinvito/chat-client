import type { ReactNode } from "react";
import { twMerge } from "tailwind-merge";
import { Toaster } from "./ui/sonner";

interface MainLayoutProps {
  children: ReactNode;
  className?: string;
}

export default function MainLayout({ children, className }: MainLayoutProps) {
  return (
    <main
      className={twMerge(
        "h-screen bg-neutral-900 text-white flex items-center",
        className,
      )}
    >
      {children}
      <Toaster richColors />
    </main>
  );
}
