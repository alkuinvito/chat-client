import { useNavigate } from "react-router";
import { useForm } from "react-hook-form";
import {
  RegisterSchema,
  TProfileSchema,
  type TRegisterSchema,
  type TResponseSchema,
} from "@/models";
import { zodResolver } from "@hookform/resolvers/zod";
import { Register } from "../../wailsjs/go/user/UserService";
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
import { Info } from "lucide-react";

function UserRegister() {
  const navigate = useNavigate();

  const form = useForm<TRegisterSchema>({
    resolver: zodResolver(RegisterSchema),
    defaultValues: {
      username: "",
    },
  });

  function onSubmit(data: TRegisterSchema) {
    Register(data.username, data.password)
      .then((res: TResponseSchema<TProfileSchema>) => {
        if (res.code === 200) {
          navigate("/login");
        } else {
          toast.error("Unknown error", {
            icon: <Info />,
          });
        }
      })
      .catch((e) => {});
  }

  return (
    <MainLayout className="flex flex-col justify-center gap-16">
      <div>
        <h1 className="text-4xl font-bold select-none">Chat Client</h1>
      </div>
      <div className="mx-auto">
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="grid gap-4 w-screen max-w-80"
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
                    <Input type="password" placeholder="Password" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Register</Button>
          </form>
        </Form>
      </div>
    </MainLayout>
  );
}

export default UserRegister;
