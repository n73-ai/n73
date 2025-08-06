import { Link, useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import { CircleUser, LogOut, User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuthStore } from "@/store/auth";
import ZustackLogo from "./zustack-logo";

export default function Navbar() {
  const { isAuth, logout } = useAuthStore();
  const navigate = useNavigate();

  return (
    <nav className="container mx-auto px-[10px] mt-[10px] xl:px-[200px]">
      <div className="flex justify-between items-center">
        <div className="flex gap-[40px]">
          <Link
            to="/"
            className="flex items-center gap-2 text-[25px] font-bold"
          >
            <ZustackLogo size={40} />
            <span>AI Zustack</span>
          </Link>
        </div>
        {isAuth ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="secondary" size="icon" className="rounded-full">
                <CircleUser className="h-5 w-5" />
                <span className="sr-only">Toggle user menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <Link to="/profile">
                <DropdownMenuItem className="flex gap-2">
                  <User className="h-4 w-4 text-zinc-300" />
                  Profile
                </DropdownMenuItem>
              </Link>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => {
                  logout();
                  navigate("/");
                }}
                className="flex gap-2"
              >
                <LogOut className="w-4 h-4 text-zinc-300" />
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
