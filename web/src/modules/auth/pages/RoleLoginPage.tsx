import { zodResolver } from "@hookform/resolvers/zod";
import { AlertCircle, Eye, EyeOff, Loader2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, Navigate, useLocation, useNavigate } from "react-router-dom";

import { getEnv } from "@/shared/config/env";
import {
  formatRoleLabel,
  getDefaultDashboardPath,
  loginPathForRole,
  primaryRole,
  resolvePostLoginRedirect,
  useAuth,
  type Role,
} from "@/modules/auth";
import {
  loginFormSchema,
  type LoginFormValues,
} from "@/modules/auth/model/login-schema";
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
import { Checkbox } from "@/shared/ui/components/checkbox";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/ui/components/form";
import { Input } from "@/shared/ui/components/input";
import { AppLoading } from "@/shared/ui/feedback/app-loading";

type RoleLoginPageProps = {
  expectedRole: Role;
  subtitle: string;
};

export function RoleLoginPage({ expectedRole, subtitle }: RoleLoginPageProps) {
  const {
    status,
    login,
    user,
  } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const from =
    (location.state as { from?: { pathname: string } } | undefined)?.from
      ?.pathname;

  const [showPassword, setShowPassword] = useState(false);
  const [serverError, setServerError] = useState<string | null>(null);

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      identifier: "",
      password: "",
      rememberMe: false,
    },
    mode: "onTouched",
  });

  if (status === "loading" && !form.formState.isSubmitting) {
    return <AppLoading fullPage label="Restoring session…" />;
  }

  if (status === "authenticated" && user) {
    const pr = primaryRole(user.roles);
    if (pr === expectedRole) {
      const dest = resolvePostLoginRedirect(from, user.roles, {
        mustChangePassword: user.mustChangePassword,
      });
      return <Navigate to={dest} replace />;
    }
    return <Navigate to={getDefaultDashboardPath(user.roles)} replace />;
  }

  async function onSubmit(values: LoginFormValues) {
    setServerError(null);
    try {
      const sessionUser = await login({
        identifier: values.identifier,
        password: values.password,
        rememberMe: values.rememberMe,
      });
      navigate(
        resolvePostLoginRedirect(from, sessionUser.roles, {
          mustChangePassword: sessionUser.mustChangePassword,
        }),
        { replace: true },
      );
    } catch (err) {
      setServerError(
        err instanceof ApiError
          ? err.message
          : "Unable to sign in. Check your connection and try again.",
      );
    }
  }

  const title = getEnv().VITE_APP_NAME;
  const loginPath = loginPathForRole(expectedRole);

  return (
    <div className="flex min-h-dvh items-center justify-center bg-gradient-to-b from-background via-background to-muted/40 p-4">
      <Card className="w-full max-w-[420px] border-border/60 shadow-lg shadow-black/5">
        <CardHeader className="space-y-1 pb-2">
          <CardTitle className="text-2xl font-semibold tracking-tight">
            {title}
          </CardTitle>
          <CardDescription className="text-base">{subtitle}</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="space-y-5"
              noValidate
            >
              {serverError ? (
                <Alert variant="destructive" className="border-destructive/40">
                  <AlertCircle className="h-4 w-4" />
                  <AlertTitle>Sign-in failed</AlertTitle>
                  <AlertDescription>{serverError}</AlertDescription>
                </Alert>
              ) : null}

              <FormField
                control={form.control}
                name="identifier"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email or phone</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        autoComplete="username"
                        placeholder="name@school.org or +1 234 567 8900"
                        disabled={form.formState.isSubmitting}
                      />
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
                    <div className="flex items-center justify-between gap-2">
                      <FormLabel>Password</FormLabel>
                    </div>
                    <FormControl>
                      <div className="relative">
                        <Input
                          {...field}
                          type={showPassword ? "text" : "password"}
                          autoComplete="current-password"
                          className="pr-11"
                          disabled={form.formState.isSubmitting}
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="icon"
                          className="absolute right-0 top-0 h-9 w-9 text-muted-foreground hover:text-foreground"
                          onClick={() => setShowPassword((v) => !v)}
                          aria-label={
                            showPassword ? "Hide password" : "Show password"
                          }
                          disabled={form.formState.isSubmitting}
                        >
                          {showPassword ? (
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
                name="rememberMe"
                render={({ field }) => (
                  <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border border-border/80 bg-muted/30 px-3 py-3">
                    <FormControl>
                      <Checkbox
                        checked={field.value}
                        onCheckedChange={(v) => field.onChange(v === true)}
                        disabled={form.formState.isSubmitting}
                      />
                    </FormControl>
                    <div className="space-y-1 leading-none">
                      <FormLabel className="font-normal">
                        Remember me on this device
                      </FormLabel>
                      <p className="text-xs text-muted-foreground">
                        When off, closing the browser ends the session for this
                        tab.
                      </p>
                      <FormMessage />
                    </div>
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
                    Signing in…
                  </>
                ) : (
                  "Sign in"
                )}
              </Button>
            </form>
          </Form>

          <p className="mt-6 text-center text-xs text-muted-foreground">
            Wrong portal?{" "}
            <Link
              to="/login"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Student
            </Link>
            {" · "}
            <Link
              to="/admin/login"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Admin
            </Link>
            {" · "}
            <Link
              to="/teacher/login"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Teacher
            </Link>
            {" · "}
            <Link
              to="/super-admin/login"
              className="font-medium text-primary underline-offset-4 hover:underline"
            >
              Super Admin
            </Link>
            {loginPath !== "/login" ? (
              <>
                <span className="mx-1">·</span>
                <span className="text-muted-foreground">
                  You are on: {formatRoleLabel(expectedRole)}
                </span>
              </>
            ) : null}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
