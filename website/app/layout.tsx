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
            <head>
                <link rel="icon" href="/favicon.ico" sizes="any" />
                <link
                    rel="apple-touch-icon"
                    sizes="76x76"
                    href={`/static/favicons/apple-touch-icon.png`}
                />
                <link
                    rel="icon"
                    type="image/png"
                    sizes="32x32"
                    href={`/static/favicons/favicon-32x32.png`}
                />
                <link
                    rel="icon"
                    type="image/png"
                    sizes="16x16"
                    href={`/static/favicons/favicon-16x16.png`}
                />
                {/* <link rel="manifest" href={`${basePath}/static/favicons/site.webmanifest`} /> */}
                <link
                    rel="mask-icon"
                    href={`/static/favicons/safari-pinned-tab.svg`}
                    color="#5bbad5"
                />
                <meta name="msapplication-TileColor" content="#000000" />
            </head>
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
