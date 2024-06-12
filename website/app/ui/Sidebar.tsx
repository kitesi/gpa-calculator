import { auth } from "@/auth";
import LoginButton from "./LoginButton";
import LoggedButtons from "./LoggedButtons";
import Classes from "./Classes";

export default async function Sidebar() {
    const session = await auth();

    return (
        <section className="h-full w-80 overflow-auto bg-gray-800">
            {!session && <LoginButton />}
            {session && (
                <LoggedButtons
                    name={session?.user?.name || ""}
                    avatar={session?.user?.image || ""}
                />
            )}
            <Classes></Classes>
        </section>
    );
}
