"use client";
import { signIn } from "next-auth/react";
import { Button } from "@headlessui/react";

export default function LoginButton() {
    return (
        <Button
            onClick={() => signIn()}
            className="text-md transform rounded-md bg-blue-700 px-20 py-2.5 font-semibold text-white shadow-md transition duration-300 ease-in-out hover:scale-105 hover:bg-blue-600 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600 focus-visible:ring-offset-2"
        >
            Login
        </Button>
    );
}
