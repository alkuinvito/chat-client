import { useEffect, useState, useCallback } from "react";
import type { user } from "../../../wailsjs/go/models";
import { GetContacts } from "../../../wailsjs/go/user/UserService";
import type { TResponseSchema } from "@/models";
import { toast } from "sonner";
import { Info } from "lucide-react";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { Input } from "../ui/input";

interface ContactsProps {
  onSelect: (contact: user.ContactModel) => void;
}

export default function ContactsPanel({ onSelect }: ContactsProps) {
  const [current, setCurrent] = useState<user.ContactModel>();
  const [contacts, setContacts] = useState<user.ContactModel[]>([]);
  const [filteredContacts, setFilteredContacts] = useState<user.ContactModel[]>(
    [],
  );
  const [searchTerm, setSearchTerm] = useState("");

  const getContacts = () => {
    GetContacts()
      .then((res: TResponseSchema<user.ContactModel[]>) => {
        switch (res.code) {
          case 200:
            setContacts(res.data);
            setFilteredContacts(res.data);
            break;
          default:
            toast.error("Error retrieving contacts", { icon: <Info /> });
            break;
        }
      })
      .catch(() => {});
  };

  const debounce = (func: (...args: any[]) => void, delay: number) => {
    let timeoutId: NodeJS.Timeout;
    return (...args: any[]) => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => {
        func(...args);
      }, delay);
    };
  };

  const handleSearch = useCallback(
    debounce((term: string) => {
      if (term.trim() === "") {
        setFilteredContacts(contacts);
      } else {
        const filtered = contacts.filter(
          (contact) =>
            contact.username.toLowerCase().includes(term.toLowerCase()) ||
            contact.id.toLowerCase().includes(term.toLowerCase()),
        );
        setFilteredContacts(filtered);
      }
    }, 300),
    [contacts],
  );

  useEffect(() => {
    getContacts();

    return EventsOn("pair:new", (contact: user.ContactModel) => {
      toast("New contact added: " + contact.username, { icon: <Info /> });
      setContacts((prev) => [...(prev || []), contact]);
      setFilteredContacts((prev) => [...(prev || []), contact]);
    });
  }, []);

  useEffect(() => {
    handleSearch(searchTerm);
  }, [searchTerm, handleSearch]);

  return (
    <div>
      <div className="p-2">
        <Input
          placeholder="Search contacts..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>
      <ul className="grow h-full overflow-y-auto">
        {filteredContacts.map((contact) => (
          <li key={contact.id}>
            <button
              className="grid px-4 py-2 w-full text-left enabled:hover:bg-neutral-700 transition-colors disabled:bg-neutral-700 disabled:hover:bg-neutral-600"
              onClick={() => {
                setCurrent(contact);
                onSelect(contact);
              }}
              disabled={contact === current}
            >
              <span className="select-none line-clamp-1">
                {contact.username}
              </span>
              <span className="select-none line-clamp-1 text-xs text-neutral-400">
                {contact.id}
              </span>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
