import { useNavigate } from "react-router";
import { useForm } from "react-hook-form";
import {
  LoginSchema,
  TProfileSchema,
  type TLoginSchema,
  type TResponseSchema,
} from "@/models";
import { zodResolver } from "@hookform/resolvers/zod";
import { Login } from "../../wailsjs/go/user/UserService";
import MainLayout from "@/components/MainLayout";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { Info, LockKeyhole, LogInIcon, WifiOff } from "lucide-react";
import logo from "@/assets/images/logo.png";
import { BrowserOpenURL, LogInfo } from "../../wailsjs/runtime";

function UserLogin() {
  const navigate = useNavigate();

  const form = useForm<TLoginSchema>({
    resolver: zodResolver(LoginSchema),
    defaultValues: {
      username: "",
    },
  });

  function onSubmit(data: TLoginSchema) {
    Login(data.username, data.password)
      .then((res: TResponseSchema<TProfileSchema>) => {
        LogInfo(res.code.toString());
        switch (res.code) {
          case 200:
            navigate("/chat");
            break;
          case 404:
            navigate("/register");
            break;
          case 401:
            toast.error("Username and/or password is incorrect", {
              icon: <Info />,
            });
            break;
          default:
            toast.error("Unknown error", {
              icon: <Info />,
            });
            break;
        }
      })
      .catch(() => {});
  }

  return (
    <MainLayout>
      <div className="grow flex flex-col md:flex-row gap-8 items-center justify-center w-full h-full mx-8">
        <div className="max-w-80 w-full float-right clear-both">
          <div>
            <h1 className="mb-1 text-4xl font-bold select-none">Chat Client</h1>
            <div className="flex items-center justify-center gap-3">
              <span className="text-xs text-neutral-500">
                another project from
              </span>
              <button
                className="outline-none border-none text-sm flex gap-1 items-center"
                onClick={() => {
                  BrowserOpenURL("https://wham.my.id");
                }}
              >
                <img src={logo} className="size-6" /> <span>wham.my.id</span>
              </button>
            </div>
          </div>

          <ul className="grid grid-rows-2 gap-4 mt-8">
            <li className="flex items-center justify-start gap-3 p-3 border border-blue-900/50 bg-blue-900/30 text-blue-400/80 text-left text-sm rounded-lg">
              <WifiOff />
              <span>No need for internet, just LAN.</span>
            </li>
            <li className="flex items-center justify-start gap-3 p-3 border border-blue-900/50 bg-blue-900/30 text-blue-400/80 text-left text-sm rounded-lg">
              <LockKeyhole />
              <span>End-to-end encrypted chat.</span>
            </li>
          </ul>
        </div>

        <div className="max-w-80 w-full mx-auto md:mx-0 p-4 border border-neutral-700 bg-neutral-800 rounded-lg">
          <h3 className="text-xl mb-4">Sign in to open chat</h3>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="grid gap-4 w-full"
            >
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input placeholder="Username" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormControl>
                      <Input
                        type="password"
                        placeholder="Password"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit">
                <LogInIcon />
                Login
              </Button>
            </form>
          </Form>
        </div>
      </div>
    </MainLayout>
  );
}

export default UserLogin;
