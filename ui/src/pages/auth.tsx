import Spinner from "@/components/spinner";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import ZustackLogo from "@/components/zustack-logo";
import { useLocation, useNavigate } from "react-router-dom";
import { useMutation } from "@tanstack/react-query";
import { authLink, authVerify } from "@/api/users";
import toast from "react-hot-toast";
import type { ErrorResponse } from "@/lib/types";
import { useAuthStore } from "@/store/auth";

export default function Auth() {
  const [email, setEmail] = useState("");
  const [otp, setOtp] = useState("");
  const [step, setStep] = useState(0);
  const location = useLocation();
  const { setAuthState } = useAuthStore();
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

  const authVerifyMutation = useMutation({
    mutationFn: () => authVerify(otp),
    onSuccess: (response) => {
      setAuthState(response.token, response.exp, response.email, true);
      navigate("/")
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

  const handleSubmitAuthVerify = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (otp === "") {
      toast.error("The One time password is required.");
      return;
    }
    authVerifyMutation.mutate();
  };

  return (
    <div className="w-full lg:grid lg:min-h-screen lg:grid-cols-2 xl:min-h-screen">
      <div className="flex items-center justify-center py-12">
        <div className="mx-auto grid w-full max-w-sm gap-4 px-[5px]">
          <div className="grid">
            {location.pathname === "/login" && (
              <h3 className="text-3xl font-semibold tracking-tight">Log in</h3>
            )}

            {location.pathname === "/signup" && (
              <h3 className="text-3xl font-semibold tracking-tight">
                Create your account
              </h3>
            )}
          </div>

          {step === 0 && (
            <form onSubmit={handleSubmitAuthLink} className="grid">
              <div className="grid gap-[10px]">
                <Label htmlFor="email">Email</Label>
                <Input
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  id="email"
                  type="email"
                  placeholder="Email"
                  required
                />
                <Button type="submit" variant="default">
                  {location.pathname === "/login" && <span>Log in</span>}
                  {location.pathname === "/signup" && (
                    <span>Create account</span>
                  )}
                  {authLinkMutation.isPending && <Spinner />}
                </Button>
              </div>
            </form>
          )}

          {step === 1 && (
            <form onSubmit={handleSubmitAuthVerify} className="grid">
              <div className="grid gap-2">
                <Label htmlFor="otp">One time password</Label>
                <Input
                  value={otp}
                  onChange={(e) => setOtp(e.target.value)}
                  id="otp"
                  type="text"
                  placeholder="One time password"
                  required
                />
                <Button type="submit" variant="default">
                  {authVerifyMutation.isPending && <Spinner />}
                  <span>Submit</span>
                </Button>
              </div>
            </form>
          )}

          {location.pathname === "/login" && (
            <div className="flex justify-center gap-1 text-center items-center">
              <p>Don't have an account?</p>
              <Button onClick={() => navigate("/signup")} variant={"link"}>
                Sign up
              </Button>
            </div>
          )}

          {location.pathname === "/signup" && (
            <div className="flex justify-center gap-1 text-center items-center">
              <p>Have an account?</p>
              <Button onClick={() => navigate("/login")} variant={"link"}>
                Log in
              </Button>
            </div>
          )}

          {step === 1 && (
            <div className="flex justify-center gap-1 text-center items-center">
              <p>Didn't receive the code?</p>
              <Button variant={"link"}>Send Again</Button>
            </div>
          )}

          <div className="flex justify-center gap-1 text-center items-center">
            <p>Go back</p>
            <Button onClick={() => navigate("/")} variant={"link"}>
              Home
            </Button>
          </div>
        </div>
      </div>
      <div
        className="hidden lg:flex flex-1 min-h-screen items-center 
        justify-center border-l border-dashed border-zinc-700"
      >
        <ZustackLogo size={300} />
      </div>
    </div>
  );
}
