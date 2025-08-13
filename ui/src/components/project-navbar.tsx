import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronLeft, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import Settings from "./settings";
import { useAuthStore } from "@/store/auth";
import toast from "react-hot-toast";

export default function ProjectNavbar() {
  const { logout, email } = useAuthStore();

  return (
    <nav className="px-[10px] h-[60px] flex items-center bg-secondary/50">
      <div className="flex justify-between w-full">
        <div className="flex gap-[10px]">
          <Link to="/">
            <Button variant="outline" className="flex gap-[5px]">
              <ChevronLeft />
            </Button>
          </Link>
          <Settings />
        </div>
        <div className="flex gap-[10px]">
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="outline" className="flex gap-[5px]">
                <User />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuLabel>{email}</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => logout()}>
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </nav>
  );
}
