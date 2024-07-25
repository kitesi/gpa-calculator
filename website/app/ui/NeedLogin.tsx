import Link from "next/link";
import Error from "@/app/ui/Error";

export default function NeedLogin() {
    return (
        <Error message="">
            You must be logged in to add a class. To log in visit the{" "}
            <Link href="/" className="text-blue-300 underline">
                home page.
            </Link>
        </Error>
    );
}
