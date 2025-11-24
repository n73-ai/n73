import { authVerify } from "@/api/users";
import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useAuthStore } from "@/store/auth";
import Spinner from "@/components/spinner";
import { useEffect } from "react";

export default function Verify() {
  const { jwt } = useParams();
  const navigate = useNavigate();

  const { setAuthState } = useAuthStore();

  const { data, isLoading, isError } = useQuery<any, Error>({
    queryKey: ["auth"],
    queryFn: () => authVerify(jwt!),
    enabled: !!jwt,
    retry: false,
  });

  useEffect(() => {
    if (data) {
      setAuthState(data.token, data.exp, data.email, true);
      navigate("/")
    }
  },[data])

  return (
    <div className="flex items-center justify-center min-h-screen">
      {isLoading && (
        <div className="flex items-center">
          Authenticating
          <Spinner />
        </div>
      )}
      {isError && (
        <div className="flex items-center gap-[3px]">
          <p>An error occurred. Please try the authentication flow again</p>
          <Link to="/login" className="underline text-white">here</Link>
        </div>
      )}
    </div>
  );
}
