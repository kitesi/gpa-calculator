import { auth } from "@/auth";
import LoginButton from "./LoginButton";
import LoggedButtons from "./LoggedButtons";
import Classes from "./Classes";
import clsx from "clsx";
import ToggleSidebarButton from "./ToggleSidebarButton";

export default async function Sidebar() {
    const session = await auth();

    return (
        <>
            <ToggleSidebarButton></ToggleSidebarButton>
            <section
                className={clsx(
                    "absolute h-full w-full -translate-x-full overflow-auto border-r-2 border-r-gray-900 bg-gray-900 transition-transform md:static md:w-80",
                    "peer-aria-pressed:pointer-events-auto peer-aria-pressed:translate-x-0 md:translate-x-0",
                )}
            >
                {!session && <LoginButton />}
                {session && (
                    <LoggedButtons
                        name={session?.user?.name || ""}
                        avatar={session?.user?.image || ""}
                    />
                )}
                <Classes></Classes>
            </section>
        </>
    );
}
