import type { user } from "../../../wailsjs/go/models";
import type { TProfileSchema } from "@/models";
import ProfilePanel from "./ProfilePanel";
import ContactsPanel from "./ContactsPanel";
import PairDialog from "./PairDialog";
import GenerateCodeDialog from "./GenerateCodeDialog";

interface SidebarProps {
  user: TProfileSchema;
  onSelect: (contact: user.ContactModel) => void;
}

export default function Sidebar({ user, onSelect }: SidebarProps) {
  return (
    <div className="flex flex-col w-screen max-w-[280px] h-full bg-neutral-800">
      <ProfilePanel user={user} />
      <div>
        <div className="grid grid-cols-2 gap-2 p-2 border-b border-b-neutral-900">
          <PairDialog />
          <GenerateCodeDialog />
        </div>
        <ContactsPanel onSelect={onSelect} />
      </div>
    </div>
  );
}
