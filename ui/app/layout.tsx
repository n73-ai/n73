import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import "./globals.css";
import { Providers } from "@/components/providers";

export const metadata: Metadata = {
  title: "Zustack AI",
  description: "Generate web applications with AI.",
  icons: {
    icon: '/logo.png', // o '/favicon.svg'
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={GeistSans.className}>
      <body>
          <Providers>{children}</Providers>
      </body>
    </html>
  );
}
