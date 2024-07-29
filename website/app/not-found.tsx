import Link from "next/link";
import Error from "@/app/ui/Error";

export default function Page404() {
    return (
        <Error message="">
            404 - Page not found. Go back to the{" "}
            <Link href="/" className="text-blue-300 underline">
                home page.
            </Link>
        </Error>
    );
}
