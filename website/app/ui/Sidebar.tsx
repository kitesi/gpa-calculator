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
                    "absolute h-full w-full max-w-96 flex-shrink-0 -translate-x-full overflow-auto border-r-[1px] border-r-midnight-700 bg-midnight-900 transition-transform md:static md:w-80",
                    "peer-aria-pressed:pointer-events-auto peer-aria-pressed:translate-x-0 md:translate-x-0",
                )}
            >
                {/* <h1 className="mb-2 text-lg font-bold">
                    Kite's GPA Calculator
                </h1> */}
                {!session && <p className="p-4">Log in!</p>}
                {/* {!session && (
                    <div className="flex w-full justify-center p-8">
                        <LoginButton></LoginButton>
                    </div>
                )} */}
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
