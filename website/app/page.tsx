import { auth } from "@/auth";
import LoginButton from "./ui/LoginButton";

export default async function Home() {
    const session = await auth();

    if (session) {
        return <></>;
    }

    return (
        <div className="p-20">
            <h1 className="text-2xl font-semibold md:mb-2 md:text-4xl lg:text-5xl">
                Kite's GPA Calculator
            </h1>
            <h2 className="mb-10 text-xl font-normal text-gray-400 md:text-2xl lg:text-3xl">
                A GPA calculator that stores your grades.
            </h2>

            <p className="mb-10 max-w-[60ch] leading-9">
                This website is a GPA calculator. You can use it to calculate
                your grades for your class, your overall GPA, what you need to
                get on your final to get a certain grade, and more. How it
                differs from other gpa calculators is that it stores your grades
                so that you don't have to re-enter everything every time. To get
                started, log in with your Google account.
            </p>
            <LoginButton />
        </div>
    );
}
