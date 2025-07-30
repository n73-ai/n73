"use client";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { LogOut, User } from "lucide-react";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { useAuthStore } from "@/store/auth";
import { useRouter } from "next/navigation";
import ZustackLogo from "../zustack_logo";

export default function Navbar() {
  const { isAuth, logout } = useAuthStore();
  const router = useRouter();

  return (
    <nav className="container mx-auto px-[10px] mt-[10px] xl:px-[200px]">
      <div className="flex justify-between items-center">
        <div className="flex gap-[40px]">
          <Link href="/" className="flex gap-2 text-[25px] items-center">
            <ZustackLogo />
            <span className="text-2xl font-bold 
            tracking-tight text-balance">
              Zustack
            </span>
          </Link>
          {isAuth && (
            <div className="hidden md:flex gap-2 items-center">
              <Button
                onClick={() => router.push("/some")}
                variant={
                  location.pathname.includes("/some") ? "secondary" : "outline"
                }
              >
                Some
              </Button>
            </div>
          )}
        </div>
        {isAuth ? (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="secondary" size="icon" className="rounded-full">
                <span className="sr-only">Toggle user menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <Link href="/profile">
                <DropdownMenuItem className="flex gap-2">
                  <User className="h-4 w-4 text-zinc-300" />
                  Profile
                </DropdownMenuItem>
              </Link>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onClick={() => {
                  logout();
                  router.push("/");
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
            <Button
              variant="outline"
              size="sm"
              onClick={() => router.push("/login")}
            >
              Sign In
            </Button>
            <Button
              variant="default"
              size="sm"
              onClick={() => router.push("/login")}
            >
              Get Started
            </Button>
          </div>
        )}
      </div>
    </nav>
  );
}
