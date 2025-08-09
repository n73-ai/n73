import { Link, useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import { User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuthStore } from "@/store/auth";
import ZustackLogo from "./zustack-logo";
import toast from "react-hot-toast";

export default function Navbar() {
  const { isAuth, logout, email } = useAuthStore();
  const navigate = useNavigate();

  return (
    <nav className="container mx-auto px-[10px] mt-[10px] xl:px-[200px]">
      <div className="flex justify-between items-center">
        <div className="flex gap-[40px]">
          <Link
            to="/"
            className="flex items-center gap-[10px] text-[25px] font-bold"
          >
            <ZustackLogo size={35} />
            <span>n73</span>
          </Link>
        </div>
        {isAuth ? (
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="outline" className="flex gap-[5px]">
                <User />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>{email}</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => toast.error("This page do not exist, yet!")}
              >
                Profile
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => toast.error("This page do not exist, yet!")}
              >
                Billing
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => logout()}>
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        ) : (
          <div className="flex gap-[10px]">
            <Button variant="secondary" onClick={() => navigate("/login")}>
              Sign In
            </Button>
            <Button variant="default" onClick={() => navigate("/signup")}>
              Get Started
            </Button>
          </div>
        )}
      </div>
    </nav>
  );
}
