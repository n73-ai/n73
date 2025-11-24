import Spinner from "@/components/spinner";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import ZustackLogo from "@/components/zustack-logo";
import { useLocation, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { authLink } from "@/api/users";
import toast from "react-hot-toast";
import type { ErrorResponse } from "@/lib/types";

export default function Auth() {
  const [email, setEmail] = useState("");
  const [step, setStep] = useState(0);
  const location = useLocation();
  const navigate = useNavigate();

  const authLinkMutation = useMutation({
    mutationFn: () => authLink(email),
    onSuccess: () => {
      setStep(1);
    },
    onError: (error: ErrorResponse) => {
      toast.error(error.response.data.error || "An unexpected error occurred.");
    },
  });

  const handleSubmitAuthLink = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (email === "") {
      toast.error("The Email is required.");
      return;
    }
    authLinkMutation.mutate();
  };

  return (
    <div className="w-full lg:grid lg:min-h-screen lg:grid-cols-2 xl:min-h-screen">
      <div className="relative flex items-center justify-center px-6 py-12 min-h-screen lg:min-h-auto">
        <div className="w-full max-w-sm">
          <div className="mb-8 text-center lg:text-left">
            {location.pathname === "/login" && (
              <h1 className="text-3xl font-semibold tracking-tight">Log in</h1>
            )}
            {location.pathname === "/signup" && (
              <h1 className="text-3xl font-semibold tracking-tight">
                Create your account
              </h1>
            )}
            <p className="text-muted-foreground mt-2">
              {location.pathname === "/login"
                ? "Welcome back! Please sign in to your account"
                : "Get started by creating your account"}
            </p>
          </div>

          <div className="space-y-6">
            {step === 0 && (
              <form onSubmit={handleSubmitAuthLink} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-sm font-medium">
                    Email address
                  </Label>
                  <Input
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    id="email"
                    type="email"
                    placeholder="Enter your email"
                    required
                    className=""
                  />
                </div>
                <Button
                  type="submit"
                  variant="secondary"
                  className="w-full"
                  disabled={authLinkMutation.isPending}
                >
                  {authLinkMutation.isPending ? (
                    <Spinner />
                  ) : (
                    <>
                      {location.pathname === "/login"
                        ? "Log in"
                        : "Create account"}
                    </>
                  )}
                </Button>
              </form>
            )}

            {step === 1 && (
              <p>
                We’ve sent you an email containing a link to complete your
                authentication.
              </p>
            )}
          </div>

          <div className="mt-8 space-y-4">
            {location.pathname === "/login" && step === 0 && (
              <div className="text-center">
                <span className="text-sm text-muted-foreground">
                  Don't have an account?{" "}
                </span>
                <Button
                  onClick={() => navigate("/signup")}
                  variant="link"
                  className="p-0 h-auto font-medium"
                >
                  Sign up
                </Button>
              </div>
            )}

            {location.pathname === "/signup" && step === 0 && (
              <div className="text-center">
                <span className="text-sm text-muted-foreground">
                  Already have an account?{" "}
                </span>
                <Button
                  onClick={() => navigate("/login")}
                  variant="link"
                  className="p-0 h-auto font-medium"
                >
                  Log in
                </Button>
              </div>
            )}

            {step === 1 && (
              <div className="space-y-3">
                <div className="text-center">
                  <span className="text-sm text-muted-foreground">
                    Didn't receive the email?{" "}
                  </span>
                  <Button
                    onClick={() => authLinkMutation.mutate()}
                    variant="link"
                    className="p-0 h-auto font-medium"
                    disabled={authLinkMutation.isPending}
                  >
                    {authLinkMutation.isPending ? "Sending..." : "Send again"}
                  </Button>
                </div>

                <div className="text-center">
                  <span className="text-sm text-muted-foreground">
                    Wrong email?{" "}
                  </span>
                  <Button
                    onClick={() => setStep(0)}
                    variant="link"
                    className="p-0 h-auto font-medium"
                  >
                    Go back
                  </Button>
                </div>
              </div>
            )}

            <div className="pt-4 border-t border-border">
              <div className="text-center">
                <Button
                  onClick={() => navigate("/")}
                  variant="ghost"
                  className="text-sm text-muted-foreground"
                >
                  ← Back to Home
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="hidden lg:flex flex-1 min-h-screen items-center justify-center border-l border-dashed border-zinc-700">
        <ZustackLogo size={400} />
      </div>
    </div>
  );
}
