import { Inter } from "next/font/google";

import QueryWrapper from "./QueryWrapper";
import Sidebar from "./ui/Sidebar";
import { Toaster } from "react-hot-toast";

import type { Metadata } from "next";
import "./globals.css";
import { SessionProvider } from "next-auth/react";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
    title: "Kite's GPA Calculator",
    description: "A GPA calculator that stores your grades.",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en" className="h-full">
            <body
                className={
                    inter.className + " h-full bg-midnight-800 text-white"
                }
            >
                <QueryWrapper>
                    <SessionProvider>
                        <div className="flex h-full">
                            <Sidebar />
                            <Toaster></Toaster>
                            {children}
                        </div>
                    </SessionProvider>
                </QueryWrapper>
            </body>
        </html>
    );
}
