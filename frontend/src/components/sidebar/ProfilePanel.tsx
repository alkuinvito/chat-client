import type { TProfileSchema } from "@/models";

interface ProfilePanelProfile {
  user: TProfileSchema;
}

export default function ProfilePanel({ user }: ProfilePanelProfile) {
  return (
    <div className="flex items-center w-full border-b border-neutral-900">
      <div className="grow grid px-4 py-2 h-full text-left">
        <span className="line-clamp-1 text-lg">{user.username}</span>
        <span className="text-xs text-neutral-400 line-clamp-1">{user.id}</span>
      </div>
    </div>
  );
}
