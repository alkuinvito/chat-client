import { useEffect, useState, useCallback } from "react";
import type { chat, user } from "../../../wailsjs/go/models";
import { GetContacts } from "../../../wailsjs/go/user/UserService";
import type { ContactList, TResponseSchema } from "@/models";
import { toast } from "sonner";
import { Info } from "lucide-react";
import { EventsOn } from "../../../wailsjs/runtime/runtime";
import { Input } from "../ui/input";

interface ContactsProps {
  onSelect: (contact: user.ContactModel) => void;
}

export default function ContactsPanel({ onSelect }: ContactsProps) {
  const [current, setCurrent] = useState<user.ContactModel>();
  const [contacts, setContacts] = useState<ContactList[]>([]);
  const [filteredContacts, setFilteredContacts] = useState<ContactList[]>([]);
  const [searchTerm, setSearchTerm] = useState("");

  const getContacts = () => {
    GetContacts()
      .then((res: TResponseSchema<user.ContactModel[]>) => {
        switch (res.code) {
          case 200:
            const contactList = res.data.map((contact) => ({
              contact: contact,
              unreadMessage: 0,
            }));

            setContacts(contactList);
            setFilteredContacts(contactList);
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
            contact.contact.username
              .toLowerCase()
              .includes(term.toLowerCase()) ||
            contact.contact.id.toLowerCase().includes(term.toLowerCase()),
        );
        setFilteredContacts(filtered);
      }
    }, 300),
    [contacts],
  );

  const handleOpenChat = (contact: user.ContactModel) => {
    // reset the unread message counter
    const targetContact = {
      contact,
      unreadMessage: 0,
    };
    const contactLists = contacts.filter((c) => c.contact.id !== contact.id);
    setContacts([...contactLists, targetContact]);
    setFilteredContacts([...contactLists, targetContact]);

    setCurrent(contact);
  };

  useEffect(() => {
    getContacts();

    // listen for new contact added
    const unsubscribeNewContact = EventsOn(
      "pair:new",
      (contact: user.ContactModel) => {
        toast("New contact added: " + contact.username, { icon: <Info /> });
        const contactList = {
          contact: contact,
          unreadMessage: 0,
        };
        setContacts((prev) => [...(prev || []), contactList]);
        setFilteredContacts((prev) => [...(prev || []), contactList]);
      },
    );

    // listen for new messages
    const unsubscribeNewMsg = EventsOn("msg:new", (msg: chat.ChatMessage) => {
      const contact = contacts.filter((c) => c.contact.id === msg.peer_id);
      if (!contact || contact.length !== 1) return;

      // set new unread message
      contact[0].unreadMessage++;
      if (contact[0].unreadMessage >= 99) {
        contact[0].unreadMessage = 99;
      }

      const contactList = contacts.filter((c) => c.contact.id !== msg.peer_id);
      setContacts([...contactList, contact[0]]);
      setFilteredContacts([...contactList, contact[0]]);
    });

    // unsubscribe all event listeners
    return () => {
      unsubscribeNewContact();
      unsubscribeNewMsg();
    };
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
          <li key={contact.contact.id}>
            <button
              className="px-4 py-2 w-full enabled:hover:bg-neutral-700 transition-colors disabled:bg-neutral-700 disabled:hover:bg-neutral-600"
              onClick={() => {
                handleOpenChat(contact.contact);
                onSelect(contact.contact);
              }}
              disabled={contact.contact === current}
            >
              <div className="flex justify-between items-center gap-2">
                <div className="grid text-left">
                  <span className="select-none line-clamp-1">
                    {contact.contact.username}
                  </span>
                  <span className="select-none line-clamp-1 text-xs text-neutral-400">
                    {contact.contact.id}
                  </span>
                </div>
                {contact.unreadMessage > 0 && (
                  <div className="shrink-0 flex justify-center items-center size-5 text-center p-1 bg-neutral-700 rounded-full">
                    <span className="text-[0.625rem] text-neutral-200 leading-2">
                      {contact.unreadMessage}
                    </span>
                  </div>
                )}
              </div>
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
