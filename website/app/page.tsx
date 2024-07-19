import { auth } from "@/auth";
import LoginButton from "./ui/LoginButton";

export default async function Home() {
    const session = await auth();

    if (session) {
        return <></>;
    }

    return (
        <div className="p-5">
            <h1 className="mb-5 border-b-4 border-b-gray-500 text-2xl font-semibold md:text-4xl">
                Kite's GPA Calculator
            </h1>
            <p className="mb-5 max-w-[60ch]">
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
