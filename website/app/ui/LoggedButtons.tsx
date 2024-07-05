"use client";
import { signOut } from "next-auth/react";
import { Button } from "@headlessui/react";
import Image from "next/image";
import Link from "next/link";

interface Props {
    name: string;
    avatar: string;
}

export default function LoggedButtons({ name, avatar }: Props) {
    return (
        <div className="flex items-center justify-around bg-gray-900 px-5 py-6">
            <Image
                src={avatar}
                alt={name}
                width={40}
                height={40}
                className="block rounded-full"
            />

            <div className="flex flex-col justify-start space-y-3 pl-6">
                <p className="font-bold">{name}</p>

                <Button
                    onClick={() => signOut()}
                    className="block rounded-lg border border-red-800 bg-red-500 px-3 py-1 text-sm font-semibold text-white hover:bg-red-600 focus:ring-2 focus:ring-red-600 focus:ring-offset-2"
                >
                    Logout
                </Button>
                <Link
                    href="/grades/new"
                    className="block whitespace-nowrap rounded-lg border border-green-800 bg-green-500 px-8 py-1 text-sm font-semibold text-black hover:bg-green-600 focus:ring-2 focus:ring-green-600 focus:ring-offset-2"
                >
                    Add Class
                </Link>
                <Link
                    href="/grades"
                    className="block rounded-lg border border-green-800 bg-green-500 px-3 py-1 text-center text-sm font-semibold text-black hover:bg-green-600 focus:ring-2 focus:ring-green-600 focus:ring-offset-2"
                >
                    GPA
                </Link>
            </div>
        </div>
    );
}
