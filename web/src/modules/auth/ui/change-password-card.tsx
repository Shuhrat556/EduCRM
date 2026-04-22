import { zodResolver } from "@hookform/resolvers/zod";
import { AlertCircle, Eye, EyeOff, Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";

import { useAuth } from "@/modules/auth";
import {
  changePasswordSchema,
  profilePasswordSchema,
  type ChangePasswordFormValues,
} from "@/modules/auth/model/change-password-schema";
import { ApiError } from "@/shared/api/types";
import {
  Alert,
  AlertDescription,
  AlertTitle,
} from "@/shared/ui/components/alert";
import { Button } from "@/shared/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/components/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/ui/components/form";
import { Input } from "@/shared/ui/components/input";

type ChangePasswordCardProps = {
  mode: "first" | "profile";
  title: string;
  description: string;
  onSuccess?: () => void;
};

export function ChangePasswordCard({
  mode,
  title,
  description,
  onSuccess,
}: ChangePasswordCardProps) {
  const { changePassword } = useAuth();
  const [serverError, setServerError] = useState<string | null>(null);
  const [show, setShow] = useState({
    current: false,
    next: false,
    confirm: false,
  });

  const schema = mode === "first" ? changePasswordSchema : profilePasswordSchema;

  const form = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
    mode: "onTouched",
  });

  async function onSubmit(values: ChangePasswordFormValues) {
    setServerError(null);
    try {
      await changePassword({
        currentPassword:
          mode === "profile" ? values.currentPassword : undefined,
        newPassword: values.newPassword,
      });
      onSuccess?.();
    } catch (err) {
      setServerError(
        err instanceof ApiError
          ? err.message
          : "Unable to update password. Try again.",
      );
    }
  }

  return (
    <Card className="w-full max-w-[440px] border-border/60 shadow-lg shadow-black/5">
      <CardHeader className="space-y-1 pb-2">
        <CardTitle className="text-xl font-semibold tracking-tight">
          {title}
        </CardTitle>
        <CardDescription className="text-base">{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onSubmit)}
            className="space-y-4"
            noValidate
          >
            {serverError ? (
              <Alert variant="destructive" className="border-destructive/40">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Could not update password</AlertTitle>
                <AlertDescription>{serverError}</AlertDescription>
              </Alert>
            ) : null}

            {mode === "profile" ? (
              <FormField
                control={form.control}
                name="currentPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Current password</FormLabel>
                    <FormControl>
                      <div className="relative">
                        <Input
                          {...field}
                          type={show.current ? "text" : "password"}
                          autoComplete="current-password"
                          className="pr-11"
                          disabled={form.formState.isSubmitting}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute right-0 top-0 h-9 w-9 text-muted-foreground hover:text-foreground"
                          onClick={() =>
                            setShow((s) => ({ ...s, current: !s.current }))
                          }
                          aria-label={
                            show.current ? "Hide password" : "Show password"
                          }
                          disabled={form.formState.isSubmitting}
                        >
                          {show.current ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            ) : null}

            <FormField
              control={form.control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>New password</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        {...field}
                        type={show.next ? "text" : "password"}
                        autoComplete="new-password"
                        className="pr-11"
                        disabled={form.formState.isSubmitting}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-0 top-0 h-9 w-9 text-muted-foreground hover:text-foreground"
                        onClick={() =>
                          setShow((s) => ({ ...s, next: !s.next }))
                        }
                        aria-label={
                          show.next ? "Hide password" : "Show password"
                        }
                        disabled={form.formState.isSubmitting}
                      >
                        {show.next ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm new password</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        {...field}
                        type={show.confirm ? "text" : "password"}
                        autoComplete="new-password"
                        className="pr-11"
                        disabled={form.formState.isSubmitting}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-0 top-0 h-9 w-9 text-muted-foreground hover:text-foreground"
                        onClick={() =>
                          setShow((s) => ({ ...s, confirm: !s.confirm }))
                        }
                        aria-label={
                          show.confirm ? "Hide password" : "Show password"
                        }
                        disabled={form.formState.isSubmitting}
                      >
                        {show.confirm ? (
                          <EyeOff className="h-4 w-4" />
                        ) : (
                          <Eye className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button
              type="submit"
              className="w-full"
              size="lg"
              disabled={form.formState.isSubmitting}
            >
              {form.formState.isSubmitting ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Saving…
                </>
              ) : mode === "first" ? (
                "Continue to dashboard"
              ) : (
                "Update password"
              )}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
