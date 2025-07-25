import type { ReactNode } from "react";
import { twMerge } from "tailwind-merge";

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
    </main>
  );
}
