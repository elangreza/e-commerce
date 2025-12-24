import { auth, User } from "@/auth";
import Header from "@/components/Header";
import LogoutButton from "@/components/LogoutButton";
import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google"; // Import User from "@/auth"
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Toko saya",
  description: "a minimalistic microservice demo",
};



export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const session = await auth();
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <div className="flex justify-center items-center  bg-zinc-50 dark:bg-gray-600">
          <div className="flex justify-center items-center max-w-5xl mx-auto w-full">
            <Header />
            <LogoutButton user={session?.user as User} />
          </div>
        </div>

        {/* <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-gray-600"> */}
        <div >
          {children}
        </div>
      </body>
    </html>
  );
}
