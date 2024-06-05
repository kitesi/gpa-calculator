"use client";
import { signIn } from "next-auth/react";
import { Button } from "@headlessui/react";

export default function LoginButton() {
    return (
        <Button
            onClick={() => signIn()}
            className="m-8 rounded-md bg-green-500 px-7 py-2.5 text-sm font-semibold text-black shadow-sm hover:bg-green-600 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600"
        >
            Login
        </Button>
    );
}
