import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  ChevronLeft,
  Settings,
  User,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";

export default function ProjectNavbar() {
  return (
    <nav className="px-[10px] h-[60px] flex items-center bg-secondary/50">
      <div className="flex justify-between w-full">
        <div className="flex gap-[10px]">
          <Link to="/">
          <Button variant="outline" className="flex gap-[5px]">
            <ChevronLeft />
          </Button>
          </Link>
          <Button variant="outline" className="flex gap-[5px]">
            <Settings />
          </Button>
        </div>
        <div className="flex gap-[10px]">
          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="outline" className="flex gap-[5px]">
                <User />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-[300px]">
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Billing</DropdownMenuItem>
              <DropdownMenuItem>Subscription</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Logout</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </nav>
  );
}
