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
        <div className="flex items-center justify-around border-b-[1px] border-b-midnight-700 p-5">
            <Image
                src={avatar}
                alt={name}
                width={40}
                height={40}
                className="block rounded-md"
            />

            <div className="flex flex-col justify-start space-y-3 pl-6">
                <p className="font-bold">{name}</p>

                <Button
                    onClick={() => signOut()}
                    className="bg-my-red block rounded-md border-none px-3 py-1 text-sm font-semibold text-white focus:ring-2 focus:ring-red-600 focus:ring-offset-2"
                >
                    Logout
                </Button>
                <Link
                    href="/grades/new"
                    className="bg-my-green block whitespace-nowrap rounded-md border-none px-8 py-1 text-sm font-semibold text-white focus:ring-2 focus:ring-green-600 focus:ring-offset-2"
                >
                    Add Class
                </Link>
                <Link
                    href="/grades"
                    className="block rounded-md border-none bg-green-500 px-3 py-1 text-center text-sm font-semibold text-black focus:ring-2 focus:ring-green-600 focus:ring-offset-2"
                >
                    GPA
                </Link>
            </div>
        </div>
    );
}
