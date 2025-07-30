"use client";
import Image from "next/image";
import Navbar from "@/components/navigation/navbar";
import CreateProject from "@/components/create_project";

export default function Page() {
  return (
    <>
      <div className="relative min-h-screen overflow-hidden">
        <Navbar />
        <div className="bb">
          <Image
            src="/gradient.svg"
            alt="Cool gradient"
            width={1200}
            height={1200}
            className="absolute top-[800px] left-1/2 transform bg-black
        -translate-x-1/2 scale-[2] pointer-events-none select-none"
          />
        </div>

        <div
          className="absolute inset-0 pointer-events-none z-10"
          style={{
            backgroundImage: "url(/grain.png)",
            backgroundSize: "100px 100px",
            backgroundRepeat: "repeat",
            mixBlendMode: "overlay",
            opacity: 0.9,
          }}
        />
        <div className="flex justify-center pt-[200px] relative z-20">
          <div className="flex flex-col">
            <h1 className="scroll-m-20 text-center text-4xl font-extrabold tracking-tight text-balance">
              Build, Preview, and Ship with AI.
            </h1>
            <p className="text-center leading-7 [&:not(:first-child)]:mt-6">
              Create apps and websites by chatting with AI.
            </p>
            <CreateProject />
          </div>
        </div>
      </div>
    </>
  );
}
