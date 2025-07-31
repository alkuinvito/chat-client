import z from "zod";

export type TResponseSchema<T> = {
  code: number;
  data: T;
};

export const RegisterSchema = z
  .object({
    username: z
      .string()
      .min(3, "Username must between 3 to 16 characters")
      .max(16, "Username must between 3 to 16 characters")
      .regex(/^[a-zA-Z0-9]+$/, {
        message: "Username must be alphanumeric",
      }),
    password: z
      .string()
      .min(8, "Password must between 8 to 32 characters")
      .max(32, "Password must between 8 to 32 characters"),
    confirm_password: z.string(),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });

export type TRegisterSchema = z.infer<typeof RegisterSchema>;

export const LoginSchema = z.object({
  username: z
    .string()
    .min(3, "Username must between 3 to 16 characters")
    .max(16, "Username must between 3 to 16 characters")
    .regex(/^[a-zA-Z0-9]+$/, {
      message: "Username must be alphanumeric",
    }),
  password: z
    .string()
    .min(8, "Password must between 8 to 32 characters")
    .max(32, "Password must between 8 to 32 characters"),
});

export type TLoginSchema = z.infer<typeof LoginSchema>;

export type TProfileSchema = {
  id: string;
  username: string;
};

export type TChatRoom = {
  peer_name: string;
  ip: string;
};

export type TPairRequestModel = {
  id: string;
  username: string;
  type: string;
};
